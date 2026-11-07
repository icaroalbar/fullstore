// Package stock defines the Inventory aggregate for managing product stock.
package stock

import (
	"errors"
	"sync"
)

var (
	// ErrInsufficientStock returned when trying to reserve without enough available quantity.
	ErrInsufficientStock = errors.New("insufficient stock")
	// ErrInvalidQuantity returned when quantity is <= 0.
	ErrInvalidQuantity = errors.New("invalid quantity: must be greater than zero")
	// ErrProductNotFound returned when product is not found in inventory.
	ErrProductNotFound = errors.New("product not found in inventory")
	// ErrInvalidOperation returned when trying to return or deduct more than reserved.
	ErrInvalidOperation = errors.New("invalid operation: quantity exceeds reserved stock")
)

// Item represents a product with its quantity in an operation (reserve, return, deduct).
type Item struct {
	ProductID string
	Quantity  int32
}

// Inventory is the aggregate root for stock management.
// It maintains a map of products, each with available and reserved quantities.
// Fields are private; access only through methods. Mutex protects concurrent access.
type Inventory struct {
	stocks map[string]*productStock
	mu     sync.Mutex
}

// productStock holds available and reserved quantities for a single product.
type productStock struct {
	available int32
	reserved  int32
}

// NewInventory creates an empty inventory.
func NewInventory() *Inventory {
	return &Inventory{
		stocks: make(map[string]*productStock),
	}
}

// Replenish adds quantity to a product's available stock.
// Creates the product entry if it doesn't exist.
func (inv *Inventory) Replenish(productID string, quantity int32) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	inv.mu.Lock()
	defer inv.mu.Unlock()

	if _, exists := inv.stocks[productID]; !exists {
		inv.stocks[productID] = &productStock{}
	}

	inv.stocks[productID].available += quantity
	return nil
}

// Available returns the available quantity for a product.
// Acquires lock to ensure atomic read in concurrent environment.
func (inv *Inventory) Available(productID string) (int32, error) {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	stock, exists := inv.stocks[productID]
	if !exists {
		return 0, ErrProductNotFound
	}

	return stock.available, nil
}

// Reserved returns the reserved quantity for a product.
// Acquires lock to ensure atomic read in concurrent environment.
func (inv *Inventory) Reserved(productID string) (int32, error) {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	stock, exists := inv.stocks[productID]
	if !exists {
		return 0, ErrProductNotFound
	}

	return stock.reserved, nil
}

// Reserve attempts to reserve quantities for a list of items.
// Validates all items before making any changes (all-or-nothing).
// Phase 1: validate all items and aggregate quantities per product (in local map).
// Phase 2: validate availability for all products.
// Phase 3: apply changes (decrement available, increment reserved) only if phases 1 & 2 pass.
func (inv *Inventory) Reserve(items []Item) error {
	if len(items) == 0 {
		return nil
	}

	inv.mu.Lock()
	defer inv.mu.Unlock()

	// Phase 1: Validate and aggregate quantities per product
	reserved := make(map[string]int32) // productID -> total quantity to reserve

	for _, item := range items {
		// Validate quantity
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}

		// Check product exists
		if _, exists := inv.stocks[item.ProductID]; !exists {
			return ErrProductNotFound
		}

		// Sum quantities if product appears multiple times
		reserved[item.ProductID] += item.Quantity
	}

	// Phase 2: Validate availability and apply changes
	for productID, totalQty := range reserved {
		if inv.stocks[productID].available < totalQty {
			return ErrInsufficientStock
		}
	}

	// All validations passed; apply changes
	for productID, totalQty := range reserved {
		stock := inv.stocks[productID]
		stock.available -= totalQty
		stock.reserved += totalQty
	}

	return nil
}

// validateAndAggregate validates items and returns a map of aggregated quantities per product.
// Phase 1: validate all items (quantity > 0, product exists).
// Returns aggregated quantities if all valid, or an error if any validation fails.
// This helper is used by both Return and Deduct (all-or-nothing pattern).
func (inv *Inventory) validateAndAggregate(items []Item) (map[string]int32, error) {
	if len(items) == 0 {
		return make(map[string]int32), nil
	}

	aggregated := make(map[string]int32)

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}

		if _, exists := inv.stocks[item.ProductID]; !exists {
			return nil, ErrProductNotFound
		}

		aggregated[item.ProductID] += item.Quantity
	}

	return aggregated, nil
}

// Return attempts to return reserved quantities for a list of items.
// Validates all items before making any changes (all-or-nothing).
// Phase 1 & 2: validate and aggregate (via helper).
// Phase 3: validate reserved is enough, then apply (increment available, decrement reserved).
func (inv *Inventory) Return(items []Item) error {
	if len(items) == 0 {
		return nil
	}

	inv.mu.Lock()
	defer inv.mu.Unlock()

	// Phase 1 & 2: validate and aggregate
	aggregated, err := inv.validateAndAggregate(items)
	if err != nil {
		return err
	}

	// Phase 3: validate reserved is enough for all products
	for productID, totalQty := range aggregated {
		if inv.stocks[productID].reserved < totalQty {
			return ErrInvalidOperation
		}
	}

	// All validations passed; apply changes
	for productID, totalQty := range aggregated {
		stock := inv.stocks[productID]
		stock.available += totalQty
		stock.reserved -= totalQty
	}

	return nil
}

// Deduct attempts to deduct reserved quantities permanently for a list of items.
// Validates all items before making any changes (all-or-nothing).
// Phase 1 & 2: validate and aggregate (via helper).
// Phase 3: validate reserved is enough, then apply (decrement reserved only).
func (inv *Inventory) Deduct(items []Item) error {
	if len(items) == 0 {
		return nil
	}

	inv.mu.Lock()
	defer inv.mu.Unlock()

	// Phase 1 & 2: validate and aggregate
	aggregated, err := inv.validateAndAggregate(items)
	if err != nil {
		return err
	}

	// Phase 3: validate reserved is enough for all products
	for productID, totalQty := range aggregated {
		if inv.stocks[productID].reserved < totalQty {
			return ErrInvalidOperation
		}
	}

	// All validations passed; apply changes (decrement reserved only)
	for productID, totalQty := range aggregated {
		stock := inv.stocks[productID]
		stock.reserved -= totalQty
	}

	return nil
}

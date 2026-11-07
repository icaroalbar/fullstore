// Package stock defines the Inventory aggregate for managing product stock.
package stock

import "sync"

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

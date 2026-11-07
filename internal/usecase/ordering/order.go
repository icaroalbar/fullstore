// Package ordering defines use cases for order management (CriarPedido, CancelarPedido, ExpirarPedidos).
package ordering

import (
	"errors"
	"fmt"
	"time"

	"github.com/icaroalbar/fullstore/internal/domain/order"
	"github.com/icaroalbar/fullstore/internal/domain/product"
	"github.com/icaroalbar/fullstore/internal/domain/stock"
)

// Errors declared by the use case
var (
	ErrProductNotFound    = errors.New("product not found")
	ErrProductUnavailable = errors.New("product unavailable")
	ErrOrderNotFound      = errors.New("order not found")
)

// ItemInput represents a line item input for creating an order.
type ItemInput struct {
	ProductID string
	Quantity  int32
}

// ProductCatalog is a port for fetching products from the catalog.
type ProductCatalog interface {
	Get(productID string) (*product.Product, error) // ErrProductNotFound
}

// OrderRepository is a port for persisting and retrieving orders.
type OrderRepository interface {
	Save(o *order.Order) error
	Update(orderID string, fn func(*order.Order) error) error // ErrOrderNotFound
	ListPending() ([]*order.Order, error)
}

// InventoryPort is a port for managing stock reservations.
type InventoryPort interface {
	Reserve(items []stock.Item) error // tudo-ou-nada; stock.ErrInsufficientStock
	Return(items []stock.Item) error  // não idempotente
}

// Clock is a port for injecting the current time.
type Clock interface {
	Now() time.Time
}

// OrderUseCases implements the three use cases for order management.
type OrderUseCases struct {
	catalog ProductCatalog
	orders  OrderRepository
	stock   InventoryPort
	clock   Clock
}

// NewOrderUseCases creates a new OrderUseCases with injected ports.
func NewOrderUseCases(
	catalog ProductCatalog,
	orders OrderRepository,
	stock InventoryPort,
	clock Clock,
) *OrderUseCases {
	return &OrderUseCases{
		catalog: catalog,
		orders:  orders,
		stock:   stock,
		clock:   clock,
	}
}

// CriarPedido creates a new order with the given items.
// Order: 1) fetch products, 2) NewOrder in memory, 3) Reserve stock, 4) Save.
// If Save fails, Return is called to compensate.
func (uc *OrderUseCases) CriarPedido(items []ItemInput) (*order.Order, error) {
	// Phase 1: Fetch products and validate
	orderItems := make([]order.Item, len(items))
	for i, item := range items {
		p, err := uc.catalog.Get(item.ProductID)
		if err != nil {
			return nil, err
		}
		// PR-04: Product must be available (active)
		if !p.IsAvailable() {
			return nil, ErrProductUnavailable
		}
		orderItems[i] = order.Item{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     p.Price(),
		}
	}

	// Phase 2: Create order in memory (validates PE-01, PE-02, PE-03)
	o, err := order.NewOrder(orderItems, uc.clock.Now())
	if err != nil {
		return nil, err
	}

	// Phase 3: Reserve stock (tudo-ou-nada; ES-01, ES-02)
	stockItems := make([]stock.Item, len(orderItems))
	for i, item := range orderItems {
		stockItems[i] = stock.Item{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	err = uc.stock.Reserve(stockItems)
	if err != nil {
		return nil, err
	}

	// Phase 4: Save order. If it fails, compensate by returning stock.
	err = uc.orders.Save(o)
	if err != nil {
		// Compensation: return the reserved stock
		returnErr := uc.stock.Return(stockItems)
		if returnErr != nil {
			// Both Save and Return failed; join errors for visibility
			return nil, errors.Join(err, returnErr)
		}
		return nil, err
	}

	return o, nil
}

// CancelarPedido cancels an order and returns its stock.
// Atomic Update first (changes state), then Return (effect outside order) only if Update succeeds.
// If Return fails, error propagates with order ID so operations can investigate.
func (uc *OrderUseCases) CancelarPedido(orderID string) error {
	var items []order.Item

	// Update: read order, call o.Cancel(), all under lock (atomic).
	// Copy items INSIDE the function passed to Update, before releasing the lock.
	err := uc.orders.Update(orderID, func(o *order.Order) error {
		items = o.Items()
		return o.Cancel()
	})

	// If Update fails (ErrOrderNotFound, ErrInvalidTransition), return error immediately.
	// Do NOT call Return (no side effect if order state didn't change).
	if err != nil {
		return err
	}

	// Update succeeded: order was canceled. Now return reserved stock.
	stockItems := make([]stock.Item, len(items))
	for i, item := range items {
		stockItems[i] = stock.Item{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	returnErr := uc.stock.Return(stockItems)
	if returnErr != nil {
		// Return failed after successful Cancel: stock is trapped. Propagate error with order ID.
		return fmt.Errorf("order %s canceled but stock return failed: %w", orderID, returnErr)
	}

	return nil
}

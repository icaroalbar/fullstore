// Package ordering defines use cases for order management (CriarPedido, CancelarPedido, ExpirarPedidos).
package ordering

import (
	"errors"
	"time"

	"github.com/icaroalbar/fullstore/internal/domain/order"
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
	Get(productID string) (*Product, error) // ErrProductNotFound
}

// Product represents a product in the catalog with pricing and availability.
type Product struct {
	ID     string
	Active bool
	Price  int64 // centavos
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

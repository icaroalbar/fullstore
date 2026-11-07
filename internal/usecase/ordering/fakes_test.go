package ordering

import (
	"sync"
	"time"

	"github.com/icaroalbar/fullstore/internal/domain/order"
	"github.com/icaroalbar/fullstore/internal/domain/product"
	"github.com/icaroalbar/fullstore/internal/domain/stock"
)

// ProductCatalogFake implements ProductCatalog for testing
type ProductCatalogFake struct {
	products map[string]*product.Product
	err      error
	mu       sync.Mutex
}

func NewProductCatalogFake() *ProductCatalogFake {
	return &ProductCatalogFake{
		products: make(map[string]*product.Product),
	}
}

func (f *ProductCatalogFake) SetProduct(p *product.Product) {
	f.SetProductByID(p.ID(), p)
}

func (f *ProductCatalogFake) SetProductByID(productID string, p *product.Product) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.products[productID] = p
}

func (f *ProductCatalogFake) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

func (f *ProductCatalogFake) Get(productID string) (*product.Product, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return nil, f.err
	}
	if p, ok := f.products[productID]; ok {
		return p, nil
	}
	return nil, ErrProductNotFound
}

// OrderRepositoryFake implements OrderRepository for testing
type OrderRepositoryFake struct {
	orders map[string]*order.Order
	err    error
	mu     sync.Mutex
}

func NewOrderRepositoryFake() *OrderRepositoryFake {
	return &OrderRepositoryFake{
		orders: make(map[string]*order.Order),
	}
}

func (f *OrderRepositoryFake) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

func (f *OrderRepositoryFake) Save(o *order.Order) error {
	if o == nil {
		panic("Save called with nil order")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return f.err
	}
	// Store the order under lock
	f.orders[o.ID()] = o
	return nil
}

func (f *OrderRepositoryFake) Update(orderID string, fn func(*order.Order) error) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return f.err
	}

	o, ok := f.orders[orderID]
	if !ok {
		return ErrOrderNotFound
	}

	// Run function under lock on the stored order
	err := fn(o)
	return err
}

func (f *OrderRepositoryFake) ListPending() ([]*order.Order, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return nil, f.err
	}

	result := make([]*order.Order, 0)
	for _, o := range f.orders {
		if o.Status() == order.StatusPending {
			result = append(result, o)
		}
	}
	return result, nil
}

// GetForTest retrieves an order for testing purposes
func (f *OrderRepositoryFake) GetForTest(orderID string) (*order.Order, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if o, ok := f.orders[orderID]; ok {
		return o, nil
	}
	return nil, ErrOrderNotFound
}

// InventoryPortFake implements InventoryPort for testing
type InventoryPortFake struct {
	stock map[string]int64
	err   error
	mu    sync.Mutex
}

func NewInventoryPortFake() *InventoryPortFake {
	return &InventoryPortFake{
		stock: make(map[string]int64),
	}
}

func (f *InventoryPortFake) SetStock(productID string, qty int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stock[productID] = qty
}

func (f *InventoryPortFake) GetStock(productID string) int64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.stock[productID]
}

func (f *InventoryPortFake) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

func (f *InventoryPortFake) Reserve(items []stock.Item) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return f.err
	}

	// Tudo ou nada: validate all before reserving any
	for _, item := range items {
		available, ok := f.stock[item.ProductID]
		if !ok || available < int64(item.Quantity) {
			return stock.ErrInsufficientStock
		}
	}

	// All items available; reserve them
	for _, item := range items {
		f.stock[item.ProductID] -= int64(item.Quantity)
	}

	return nil
}

func (f *InventoryPortFake) Return(items []stock.Item) error {
	if items == nil {
		panic("Return called with nil items")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return f.err
	}

	for _, item := range items {
		f.stock[item.ProductID] += int64(item.Quantity)
	}

	return nil
}

// ClockFake implements Clock for testing
type ClockFake struct {
	now time.Time
	mu  sync.Mutex
}

func NewClockFake(now time.Time) *ClockFake {
	return &ClockFake{now: now}
}

func (f *ClockFake) SetNow(now time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = now
}

func (f *ClockFake) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.now
}

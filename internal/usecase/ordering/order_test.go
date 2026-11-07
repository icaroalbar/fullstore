package ordering

import (
	"testing"
	"time"

	"github.com/icaroalbar/fullstore/internal/domain/order"
	"github.com/icaroalbar/fullstore/internal/domain/stock"
)

// T002: Portas e erros do caso de uso
func TestOrderUseCases_ConstructorAndFakes(t *testing.T) {
	// Fakes
	catalogFake := &ProductCatalogFake{}
	orderRepoFake := &OrderRepositoryFake{}
	inventoryFake := &InventoryPortFake{}
	clockFake := &ClockFake{now: time.Now()}

	// Constructor should succeed
	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)
	if uc == nil {
		t.Errorf("NewOrderUseCases should return non-nil OrderUseCases")
	}
}

// Fake implementations
type ProductCatalogFake struct {
	products map[string]*Product
	err      error
}

func (f *ProductCatalogFake) Get(productID string) (*Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	if p, ok := f.products[productID]; ok {
		return p, nil
	}
	return nil, ErrProductNotFound
}

type OrderRepositoryFake struct {
	orders map[string]*order.Order
	err    error
}

func (f *OrderRepositoryFake) Save(o *order.Order) error {
	if f.err != nil {
		return f.err
	}
	f.orders[o.ID()] = o
	return nil
}

func (f *OrderRepositoryFake) Update(orderID string, fn func(*order.Order) error) error {
	if f.err != nil {
		return f.err
	}
	if o, ok := f.orders[orderID]; ok {
		return fn(o)
	}
	return ErrOrderNotFound
}

func (f *OrderRepositoryFake) ListPending() ([]*order.Order, error) {
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

type InventoryPortFake struct {
	stock map[string]int64
	err   error
}

func (f *InventoryPortFake) Reserve(items []stock.Item) error {
	if f.err != nil {
		return f.err
	}
	for _, item := range items {
		if available, ok := f.stock[item.ProductID]; !ok || available < int64(item.Quantity) {
			return stock.ErrInsufficientStock
		}
		f.stock[item.ProductID] -= int64(item.Quantity)
	}
	return nil
}

func (f *InventoryPortFake) Return(items []stock.Item) error {
	if f.err != nil {
		return f.err
	}
	for _, item := range items {
		f.stock[item.ProductID] += int64(item.Quantity)
	}
	return nil
}

type ClockFake struct {
	now time.Time
}

func (f *ClockFake) Now() time.Time {
	return f.now
}

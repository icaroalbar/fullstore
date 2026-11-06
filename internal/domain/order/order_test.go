package order

import (
	"fmt"
	"testing"
	"time"
)

// T001: PE-01 — validação de quantidade de itens
func TestOrder_PE01_EmptyOrder(t *testing.T) {
	_, err := NewOrder([]Item{}, time.Now())
	if err != ErrEmptyOrder {
		t.Errorf("empty order should return ErrEmptyOrder, got %v", err)
	}
}

func TestOrder_PE01_TooManyItems(t *testing.T) {
	items := make([]Item, 21)
	for i := 0; i < 21; i++ {
		items[i] = Item{ProductID: "p1", Quantity: 1, Price: 100}
	}
	_, err := NewOrder(items, time.Now())
	if err != ErrTooManyItems {
		t.Errorf("21 items should return ErrTooManyItems, got %v", err)
	}
}

func TestOrder_PE01_ValidItemCount(t *testing.T) {
	tests := []int{1, 10, 20}
	for _, count := range tests {
		items := make([]Item, count)
		for i := 0; i < count; i++ {
			items[i] = Item{ProductID: fmt.Sprintf("p%d", i), Quantity: 1, Price: 100}
		}
		_, err := NewOrder(items, time.Now())
		if err != nil {
			t.Errorf("order with %d items should not error, got %v", count, err)
		}
	}
}

// T001: PE-02 — validação de quantidade por item
func TestOrder_PE02_InvalidQuantityZero(t *testing.T) {
	items := []Item{
		{ProductID: "p1", Quantity: 0, Price: 100},
	}
	_, err := NewOrder(items, time.Now())
	if err != ErrInvalidQuantity {
		t.Errorf("item with qty 0 should return ErrInvalidQuantity, got %v", err)
	}
}

func TestOrder_PE02_InvalidQuantityTooHigh(t *testing.T) {
	items := []Item{
		{ProductID: "p1", Quantity: 11, Price: 100},
	}
	_, err := NewOrder(items, time.Now())
	if err != ErrInvalidQuantity {
		t.Errorf("item with qty 11 should return ErrInvalidQuantity, got %v", err)
	}
}

func TestOrder_PE02_ValidQuantity(t *testing.T) {
	tests := []int{1, 5, 10}
	for _, qty := range tests {
		items := []Item{
			{ProductID: "p1", Quantity: int32(qty), Price: 100},
		}
		_, err := NewOrder(items, time.Now())
		if err != nil {
			t.Errorf("item with qty %d should not error, got %v", qty, err)
		}
	}
}

// T001: Verificar Status, ID, CreatedAt
func TestOrder_PE01_StatusPending(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o, _ := NewOrder(items, time.Now())
	if o.Status() != "pending" {
		t.Errorf("new order should have status pending, got %s", o.Status())
	}
}

func TestOrder_PE01_HasUniqueID(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o1, _ := NewOrder(items, time.Now())
	o2, _ := NewOrder(items, time.Now())
	if o1.ID() == o2.ID() {
		t.Errorf("orders should have unique IDs, got %s and %s", o1.ID(), o2.ID())
	}
	if o1.ID() == "" {
		t.Errorf("order ID should not be empty")
	}
}

func TestOrder_PE01_HasCreatedAt(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	before := time.Now()
	o, _ := NewOrder(items, time.Now())
	after := time.Now()
	if o.CreatedAt().Before(before) || o.CreatedAt().After(after) {
		t.Errorf("order CreatedAt should be around now, got %v", o.CreatedAt())
	}
}

// T001+: Getters protection — Items() retorna cópia, não referência
func TestOrder_ItemsReturnsCopy(t *testing.T) {
	originalItems := []Item{
		{ProductID: "p1", Quantity: 2, Price: 100},
		{ProductID: "p2", Quantity: 1, Price: 200},
	}
	o, _ := NewOrder(originalItems, time.Now())
	originalTotal := o.Total()

	// Pega os itens via getter
	retrievedItems := o.Items()

	// Altera o slice retornado
	retrievedItems[0].Price = 9999

	// Confirma que o pedido não mudou
	if o.Total() != originalTotal {
		t.Errorf("altering returned Items slice should not change Order.Total(), got %d, want %d", o.Total(), originalTotal)
	}

	// Pega items novamente e confirma que tem 2 (não 3)
	itemsAgain := o.Items()
	if len(itemsAgain) != 2 {
		t.Errorf("order should still have 2 items, got %d", len(itemsAgain))
	}
	if itemsAgain[0].Price != 100 {
		t.Errorf("order's first item price should still be 100, got %d", itemsAgain[0].Price)
	}
}

// T001+: Getters — Total, Status, ID
func TestOrder_GetterID(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o, _ := NewOrder(items, time.Now())
	if o.ID() == "" {
		t.Errorf("ID() should not be empty")
	}
	// ID deve ser UUID-like (36 chars com hífens)
	if len(o.ID()) != 36 {
		t.Errorf("ID() should be UUID format (36 chars), got %d", len(o.ID()))
	}
}

func TestOrder_GetterTotal(t *testing.T) {
	items := []Item{
		{ProductID: "p1", Quantity: 2, Price: 100}, // 200
		{ProductID: "p2", Quantity: 3, Price: 150}, // 450
	}
	o, _ := NewOrder(items, time.Now())
	expectedTotal := int64(200 + 450)
	if o.Total() != expectedTotal {
		t.Errorf("Total() should be %d, got %d", expectedTotal, o.Total())
	}
}

func TestOrder_GetterStatus(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o, _ := NewOrder(items, time.Now())
	if o.Status() != "pending" {
		t.Errorf("Status() should be 'pending', got %q", o.Status())
	}
}

func TestOrder_GetterPaidAtNil(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o, _ := NewOrder(items, time.Now())
	if o.PaidAt() != nil {
		t.Errorf("PaidAt() should be nil for new order, got %v", o.PaidAt())
	}
}

func TestOrder_GetterExpiredAtNil(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o, _ := NewOrder(items, time.Now())
	if o.ExpiredAt() != nil {
		t.Errorf("ExpiredAt() should be nil for new order, got %v", o.ExpiredAt())
	}
}

// Validações adicionais: preço <= 0, duplicata, cópia de slice, injeção relógio
func TestOrder_InvalidPriceZero(t *testing.T) {
	items := []Item{
		{ProductID: "p1", Quantity: 1, Price: 0},
	}
	_, err := NewOrder(items, time.Now())
	if err != ErrInvalidPrice {
		t.Errorf("item with price 0 should return ErrInvalidPrice, got %v", err)
	}
}

func TestOrder_InvalidPriceNegative(t *testing.T) {
	items := []Item{
		{ProductID: "p1", Quantity: 1, Price: -100},
	}
	_, err := NewOrder(items, time.Now())
	if err != ErrInvalidPrice {
		t.Errorf("item with negative price should return ErrInvalidPrice, got %v", err)
	}
}

func TestOrder_DuplicateProductID(t *testing.T) {
	items := []Item{
		{ProductID: "p1", Quantity: 1, Price: 100},
		{ProductID: "p1", Quantity: 2, Price: 200},
	}
	_, err := NewOrder(items, time.Now())
	if err != ErrDuplicateItem {
		t.Errorf("duplicate ProductID should return ErrDuplicateItem, got %v", err)
	}
}

func TestOrder_SliceCopiedOnCreation(t *testing.T) {
	originalItems := []Item{
		{ProductID: "p1", Quantity: 1, Price: 100},
		{ProductID: "p2", Quantity: 2, Price: 200},
	}
	o, _ := NewOrder(originalItems, time.Now())

	// Altera o slice original
	originalItems[0].Price = 9999

	// Confirma que o pedido tem o total original
	expectedTotal := int64(100*1 + 200*2)
	if o.Total() != expectedTotal {
		t.Errorf("order should have original total, got %d, want %d", o.Total(), expectedTotal)
	}

	// Confirma que o pedido tem 2 items (não 3)
	if len(o.Items()) != 2 {
		t.Errorf("order should have 2 items, got %d", len(o.Items()))
	}
}

func TestOrder_CreatedAtUseInjectedTime(t *testing.T) {
	customTime := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o, _ := NewOrder(items, customTime)

	if o.CreatedAt() != customTime {
		t.Errorf("order CreatedAt should use injected time, got %v, want %v", o.CreatedAt(), customTime)
	}
}

func TestOrder_StatusIsConstant(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o, _ := NewOrder(items, time.Now())
	if o.Status() != StatusPending {
		t.Errorf("new order status should be StatusPending, got %v", o.Status())
	}
}

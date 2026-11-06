package order

import (
	"testing"
	"time"
)

// T001: PE-01 — validação de quantidade de itens
func TestOrder_PE01_EmptyOrder(t *testing.T) {
	_, err := NewOrder([]Item{})
	if err != ErrEmptyOrder {
		t.Errorf("empty order should return ErrEmptyOrder, got %v", err)
	}
}

func TestOrder_PE01_TooManyItems(t *testing.T) {
	items := make([]Item, 21)
	for i := 0; i < 21; i++ {
		items[i] = Item{ProductID: "p1", Quantity: 1, Price: 100}
	}
	_, err := NewOrder(items)
	if err != ErrTooManyItems {
		t.Errorf("21 items should return ErrTooManyItems, got %v", err)
	}
}

func TestOrder_PE01_ValidItemCount(t *testing.T) {
	tests := []int{1, 10, 20}
	for _, count := range tests {
		items := make([]Item, count)
		for i := 0; i < count; i++ {
			items[i] = Item{ProductID: "p1", Quantity: 1, Price: 100}
		}
		_, err := NewOrder(items)
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
	_, err := NewOrder(items)
	if err != ErrInvalidQuantity {
		t.Errorf("item with qty 0 should return ErrInvalidQuantity, got %v", err)
	}
}

func TestOrder_PE02_InvalidQuantityTooHigh(t *testing.T) {
	items := []Item{
		{ProductID: "p1", Quantity: 11, Price: 100},
	}
	_, err := NewOrder(items)
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
		_, err := NewOrder(items)
		if err != nil {
			t.Errorf("item with qty %d should not error, got %v", qty, err)
		}
	}
}

// T001: Verificar Status, ID, CreatedAt
func TestOrder_PE01_StatusPending(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o, _ := NewOrder(items)
	if o.Status() != "pending" {
		t.Errorf("new order should have status pending, got %s", o.Status())
	}
}

func TestOrder_PE01_HasUniqueID(t *testing.T) {
	items := []Item{{ProductID: "p1", Quantity: 1, Price: 100}}
	o1, _ := NewOrder(items)
	o2, _ := NewOrder(items)
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
	o, _ := NewOrder(items)
	after := time.Now()
	if o.CreatedAt().Before(before) || o.CreatedAt().After(after) {
		t.Errorf("order CreatedAt should be around now, got %v", o.CreatedAt())
	}
}

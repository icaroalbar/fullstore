package stock

import (
	"testing"
)

// T001: Tipos base Inventory, Replenish e getters

func TestInventory_ReplenishPositive(t *testing.T) {
	inv := NewInventory()
	err := inv.Replenish("p1", 5)
	if err != nil {
		t.Errorf("Replenish(p1, 5) should not error, got %v", err)
	}

	available, err := inv.Available("p1")
	if err != nil {
		t.Errorf("Available(p1) should not error, got %v", err)
	}
	if available != 5 {
		t.Errorf("after Replenish(5), Available should be 5, got %d", available)
	}
}

func TestInventory_ReplenishInvalidQuantity(t *testing.T) {
	inv := NewInventory()
	err := inv.Replenish("p1", 0)
	if err != ErrInvalidQuantity {
		t.Errorf("Replenish(p1, 0) should return ErrInvalidQuantity, got %v", err)
	}
}

func TestInventory_MultipleReplenish(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 5); err != nil {
		t.Errorf("first Replenish should not error, got %v", err)
	}
	err := inv.Replenish("p1", 3)
	if err != nil {
		t.Errorf("second Replenish should not error, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 8 {
		t.Errorf("after two Replenish (5+3), Available should be 8, got %d", available)
	}
}

func TestInventory_GetterProductNotFound(t *testing.T) {
	inv := NewInventory()
	_, err := inv.Available("unknown")
	if err != ErrProductNotFound {
		t.Errorf("Available(unknown) should return ErrProductNotFound, got %v", err)
	}
}

func TestInventory_GettersUseLock(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	// Teste simples que getters não pânico; -race detectará data races
	available, err := inv.Available("p1")
	if err != nil {
		t.Errorf("Available should not error, got %v", err)
	}
	if available != 10 {
		t.Errorf("Available should be 10, got %d", available)
	}

	reserved, err := inv.Reserved("p1")
	if err != nil {
		t.Errorf("Reserved should not error, got %v", err)
	}
	if reserved != 0 {
		t.Errorf("Reserved should be 0, got %d", reserved)
	}
}

// T002: Reserve com validação — ES-01, ES-02

func TestInventory_ReserveEmptyList(t *testing.T) {
	inv := NewInventory()
	err := inv.Reserve([]Item{})
	if err != nil {
		t.Errorf("Reserve with empty list should not error, got %v", err)
	}
}

func TestInventory_ES01_ReserveDecrementsAvailable(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}})
	if err != nil {
		t.Errorf("Reserve(3) on 10 available should not error, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 7 {
		t.Errorf("after Reserve(3), available should be 7, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 3 {
		t.Errorf("after Reserve(3), reserved should be 3, got %d", reserved)
	}
}

func TestInventory_ES02_InsufficientStock(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 15}})
	if err != ErrInsufficientStock {
		t.Errorf("Reserve(15) on 10 should return ErrInsufficientStock, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 10 {
		t.Errorf("after failed Reserve, available should still be 10, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 0 {
		t.Errorf("after failed Reserve, reserved should still be 0, got %d", reserved)
	}
}

func TestInventory_ReserveInvalidQuantity(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 0}})
	if err != ErrInvalidQuantity {
		t.Errorf("Reserve(0) should return ErrInvalidQuantity, got %v", err)
	}
}

func TestInventory_ReserveProductNotFound(t *testing.T) {
	inv := NewInventory()
	err := inv.Reserve([]Item{{ProductID: "unknown", Quantity: 1}})
	if err != ErrProductNotFound {
		t.Errorf("Reserve with unknown product should return ErrProductNotFound, got %v", err)
	}
}

func TestInventory_ReserveMultipleItems(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish p1 should not error, got %v", err)
	}
	if err := inv.Replenish("p2", 5); err != nil {
		t.Errorf("Replenish p2 should not error, got %v", err)
	}

	err := inv.Reserve([]Item{
		{ProductID: "p1", Quantity: 3},
		{ProductID: "p2", Quantity: 2},
	})
	if err != nil {
		t.Errorf("Reserve multiple items should not error, got %v", err)
	}

	available1, _ := inv.Available("p1")
	if available1 != 7 {
		t.Errorf("p1 available should be 7, got %d", available1)
	}

	available2, _ := inv.Available("p2")
	if available2 != 3 {
		t.Errorf("p2 available should be 3, got %d", available2)
	}
}

func TestInventory_ReserveProductRepetition(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	err := inv.Reserve([]Item{
		{ProductID: "p1", Quantity: 2},
		{ProductID: "p1", Quantity: 3},
	})
	if err != nil {
		t.Errorf("Reserve with repeated product should sum and not error, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 5 {
		t.Errorf("after Reserve(2+3), available should be 5, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 5 {
		t.Errorf("after Reserve(2+3), reserved should be 5, got %d", reserved)
	}
}

func TestInventory_ES02b_AtomicFail(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish p1 should not error, got %v", err)
	}
	if err := inv.Replenish("p2", 10); err != nil {
		t.Errorf("Replenish p2 should not error, got %v", err)
	}

	err := inv.Reserve([]Item{
		{ProductID: "p1", Quantity: 3},
		{ProductID: "p2", Quantity: 15},
	})
	if err != ErrInsufficientStock {
		t.Errorf("Reserve should fail atomically (p2 has only 10), got %v", err)
	}

	available1, _ := inv.Available("p1")
	if available1 != 10 {
		t.Errorf("after failed Reserve, p1 available should still be 10, got %d", available1)
	}

	available2, _ := inv.Available("p2")
	if available2 != 10 {
		t.Errorf("after failed Reserve, p2 available should still be 10, got %d", available2)
	}
}

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

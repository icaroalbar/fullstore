package stock

import (
	"sync"
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

// T003: Return e Deduct — ES-03, ES-04

func TestInventory_ES03a_ReturnAfterCancel(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Return([]Item{{ProductID: "p1", Quantity: 3}})
	if err != nil {
		t.Errorf("Return(3) should not error, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 10 {
		t.Errorf("after Return, available should be 10, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 0 {
		t.Errorf("after Return, reserved should be 0, got %d", reserved)
	}
}

func TestInventory_ES03b_ReturnAfterExpire(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Return([]Item{{ProductID: "p1", Quantity: 3}})
	if err != nil {
		t.Errorf("Return(3) on expire should not error, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 10 {
		t.Errorf("after Return (expire), available should be 10, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 0 {
		t.Errorf("after Return (expire), reserved should be 0, got %d", reserved)
	}
}

func TestInventory_ReturnInvalidQuantity(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}
	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Return([]Item{{ProductID: "p1", Quantity: 0}})
	if err != ErrInvalidQuantity {
		t.Errorf("Return(0) should return ErrInvalidQuantity, got %v", err)
	}
}

func TestInventory_ReturnProductNotFound(t *testing.T) {
	inv := NewInventory()
	err := inv.Return([]Item{{ProductID: "unknown", Quantity: 1}})
	if err != ErrProductNotFound {
		t.Errorf("Return with unknown product should return ErrProductNotFound, got %v", err)
	}
}

func TestInventory_ReturnTooMuch(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}
	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Return([]Item{{ProductID: "p1", Quantity: 5}})
	if err != ErrInvalidOperation {
		t.Errorf("Return(5) on reserved=3 should return ErrInvalidOperation, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 7 {
		t.Errorf("after failed Return, available should still be 7, got %d", available)
	}
}

func TestInventory_ReturnProductRepetition(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 5}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Return([]Item{
		{ProductID: "p1", Quantity: 2},
		{ProductID: "p1", Quantity: 3},
	})
	if err != nil {
		t.Errorf("Return with repeated product should sum and not error, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 10 {
		t.Errorf("after Return(2+3), available should be 10, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 0 {
		t.Errorf("after Return(2+3), reserved should be 0, got %d", reserved)
	}
}

func TestInventory_ReturnAtomicFail(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish p1 should not error, got %v", err)
	}
	if err := inv.Replenish("p2", 10); err != nil {
		t.Errorf("Replenish p2 should not error, got %v", err)
	}
	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}, {ProductID: "p2", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Return([]Item{
		{ProductID: "p1", Quantity: 2},
		{ProductID: "p2", Quantity: 5},
	})
	if err != ErrInvalidOperation {
		t.Errorf("Return should fail atomically (p2 reserved=3, try return 5), got %v", err)
	}

	available1, _ := inv.Available("p1")
	if available1 != 7 {
		t.Errorf("after failed Return, p1 available should still be 7, got %d", available1)
	}

	reserved1, _ := inv.Reserved("p1")
	if reserved1 != 3 {
		t.Errorf("after failed Return, p1 reserved should still be 3, got %d", reserved1)
	}
}

func TestInventory_ES04_DeductAfterPay(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Deduct([]Item{{ProductID: "p1", Quantity: 3}})
	if err != nil {
		t.Errorf("Deduct(3) should not error, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 7 {
		t.Errorf("after Deduct, available should still be 7, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 0 {
		t.Errorf("after Deduct, reserved should be 0, got %d", reserved)
	}
}

func TestInventory_DeductProductRepetition(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 5}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Deduct([]Item{
		{ProductID: "p1", Quantity: 2},
		{ProductID: "p1", Quantity: 3},
	})
	if err != nil {
		t.Errorf("Deduct with repeated product should sum and not error, got %v", err)
	}

	available, _ := inv.Available("p1")
	if available != 5 {
		t.Errorf("after Deduct(2+3), available should still be 5, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 0 {
		t.Errorf("after Deduct(2+3), reserved should be 0, got %d", reserved)
	}
}

func TestInventory_DeductInvalidQuantity(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}
	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Deduct([]Item{{ProductID: "p1", Quantity: 0}})
	if err != ErrInvalidQuantity {
		t.Errorf("Deduct(0) should return ErrInvalidQuantity, got %v", err)
	}
}

func TestInventory_DeductProductNotFound(t *testing.T) {
	inv := NewInventory()
	err := inv.Deduct([]Item{{ProductID: "unknown", Quantity: 1}})
	if err != ErrProductNotFound {
		t.Errorf("Deduct with unknown product should return ErrProductNotFound, got %v", err)
	}
}

func TestInventory_DeductTooMuch(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}
	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Deduct([]Item{{ProductID: "p1", Quantity: 5}})
	if err != ErrInvalidOperation {
		t.Errorf("Deduct(5) on reserved=3 should return ErrInvalidOperation, got %v", err)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 3 {
		t.Errorf("after failed Deduct, reserved should still be 3, got %d", reserved)
	}
}

func TestInventory_DeductAtomicFail(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish p1 should not error, got %v", err)
	}
	if err := inv.Replenish("p2", 10); err != nil {
		t.Errorf("Replenish p2 should not error, got %v", err)
	}
	if err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 3}, {ProductID: "p2", Quantity: 3}}); err != nil {
		t.Errorf("Reserve should not error, got %v", err)
	}

	err := inv.Deduct([]Item{
		{ProductID: "p1", Quantity: 2},
		{ProductID: "p2", Quantity: 5},
	})
	if err != ErrInvalidOperation {
		t.Errorf("Deduct should fail atomically (p2 reserved=3, try deduct 5), got %v", err)
	}

	reserved1, _ := inv.Reserved("p1")
	if reserved1 != 3 {
		t.Errorf("after failed Deduct, p1 reserved should still be 3, got %d", reserved1)
	}
}

// T004: Concorrência — ES-05

func TestInventory_ES05_RaceCondition(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	successCount := 0
	failureCount := 0
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := inv.Reserve([]Item{{ProductID: "p1", Quantity: 1}})
			mu.Lock()
			switch err {
			case nil:
				successCount++
			case ErrInsufficientStock:
				failureCount++
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	if successCount != 10 {
		t.Errorf("ES05: expected exactly 10 successes, got %d", successCount)
	}
	if failureCount != 40 {
		t.Errorf("ES05: expected exactly 40 failures, got %d", failureCount)
	}
}

func TestInventory_ES05_FinalState(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = inv.Reserve([]Item{{ProductID: "p1", Quantity: 1}})
		}()
	}
	wg.Wait()

	available, _ := inv.Available("p1")
	if available != 0 {
		t.Errorf("after race, available should be 0, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 10 {
		t.Errorf("after race, reserved should be 10, got %d", reserved)
	}
}

func TestInventory_ES05_ConcurrentReads(t *testing.T) {
	inv := NewInventory()
	if err := inv.Replenish("p1", 10); err != nil {
		t.Errorf("Replenish should not error, got %v", err)
	}

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = inv.Reserve([]Item{{ProductID: "p1", Quantity: 1}})
		}()
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = inv.Available("p1")
			_, _ = inv.Reserved("p1")
		}()
	}

	wg.Wait()

	available, _ := inv.Available("p1")
	if available != 0 {
		t.Errorf("after concurrent reads/writes, available should be 0, got %d", available)
	}

	reserved, _ := inv.Reserved("p1")
	if reserved != 10 {
		t.Errorf("after concurrent reads/writes, reserved should be 10, got %d", reserved)
	}
}

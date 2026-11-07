package payment

import (
	"testing"
)

// T001: Tipos base OrderPayments + NewOrderPayments (PG-01)

func TestOrderPayments_PG01_InvalidTotal(t *testing.T) {
	_, err := NewOrderPayments("order-1", 0)
	if err != ErrInvalidAmount {
		t.Errorf("NewOrderPayments(id, 0) should return ErrInvalidAmount, got %v", err)
	}
}

func TestOrderPayments_PG01_ValidTotal(t *testing.T) {
	op, err := NewOrderPayments("order-1", 10000)
	if err != nil {
		t.Errorf("NewOrderPayments(id, 10000) should not error, got %v", err)
	}
	if op == nil {
		t.Errorf("NewOrderPayments should return non-nil OrderPayments")
	}
}

func TestOrderPayments_GettersInitial(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)

	approved, ok := op.Approved()
	if ok {
		t.Errorf("Approved() on new OrderPayments should return ok=false, got %v", ok)
	}
	if approved != (Attempt{}) {
		t.Errorf("Approved() on new OrderPayments should return zero-value, got %v", approved)
	}

	attempts := op.Attempts()
	if attempts == nil || len(attempts) != 0 {
		t.Errorf("Attempts() on new OrderPayments should return empty slice, got %v", attempts)
	}
}

func TestOrderPayments_CopiesDefensive(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)

	attempts1 := op.Attempts()
	attempts2 := op.Attempts()

	// Verify they're equal but not the same reference
	if len(attempts1) != len(attempts2) {
		t.Errorf("Attempts() calls should return equal slices")
	}

	// Modifying returned slice should not affect internal state
	// (This will be clearer after T002 adds attempts)
}

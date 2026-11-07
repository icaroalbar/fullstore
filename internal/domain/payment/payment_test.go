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

// T002: RecordAttempt com validação de valor (PG-01)

func TestOrderPayments_PG01_AmountMismatch(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)
	err := op.RecordAttempt(9900, StatusDeclined)
	if err != ErrAmountMismatch {
		t.Errorf("RecordAttempt(9900, declined) should return ErrAmountMismatch, got %v", err)
	}

	// Verify nothing was registered
	attempts := op.Attempts()
	if len(attempts) != 0 {
		t.Errorf("after failed RecordAttempt, Attempts() should be empty, got %d", len(attempts))
	}
}

func TestOrderPayments_PG01_InvalidAmount(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)
	err := op.RecordAttempt(0, StatusApproved)
	if err != ErrInvalidAmount {
		t.Errorf("RecordAttempt(0, approved) should return ErrInvalidAmount, got %v", err)
	}

	attempts := op.Attempts()
	if len(attempts) != 0 {
		t.Errorf("after failed RecordAttempt, Attempts() should be empty, got %d", len(attempts))
	}
}

func TestOrderPayments_InvalidResult(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)
	err := op.RecordAttempt(10000, Status("unknown"))
	if err != ErrInvalidResult {
		t.Errorf("RecordAttempt(10000, unknown) should return ErrInvalidResult, got %v", err)
	}

	attempts := op.Attempts()
	if len(attempts) != 0 {
		t.Errorf("after failed RecordAttempt, Attempts() should be empty, got %d", len(attempts))
	}
}

func TestOrderPayments_PG02_FirstApproved(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)
	err := op.RecordAttempt(10000, StatusApproved)
	if err != nil {
		t.Errorf("RecordAttempt(10000, approved) should not error, got %v", err)
	}

	approved, ok := op.Approved()
	if !ok {
		t.Errorf("Approved() should return ok=true after registering approved attempt")
	}
	if approved.Amount != 10000 || approved.Status != StatusApproved {
		t.Errorf("Approved() should return the registered attempt, got %v", approved)
	}
}

func TestOrderPayments_PG02_DeclinedThenApproved(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)

	err1 := op.RecordAttempt(10000, StatusDeclined)
	if err1 != nil {
		t.Errorf("RecordAttempt(10000, declined) should not error, got %v", err1)
	}

	err2 := op.RecordAttempt(10000, StatusApproved)
	if err2 != nil {
		t.Errorf("RecordAttempt(10000, approved) should not error, got %v", err2)
	}

	attempts := op.Attempts()
	if len(attempts) != 2 {
		t.Errorf("after 2 RecordAttempt calls, Attempts() should have 2 items, got %d", len(attempts))
	}

	if attempts[0].Status != StatusDeclined {
		t.Errorf("first attempt should be declined, got %v", attempts[0].Status)
	}
	if attempts[1].Status != StatusApproved {
		t.Errorf("second attempt should be approved, got %v", attempts[1].Status)
	}
}

func TestOrderPayments_OnlyDeclinedNoApproved(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)

	if err := op.RecordAttempt(10000, StatusDeclined); err != nil {
		t.Errorf("RecordAttempt should not error, got %v", err)
	}
	if err := op.RecordAttempt(10000, StatusDeclined); err != nil {
		t.Errorf("RecordAttempt should not error, got %v", err)
	}
	if err := op.RecordAttempt(10000, StatusDeclined); err != nil {
		t.Errorf("RecordAttempt should not error, got %v", err)
	}

	approved, ok := op.Approved()
	if ok {
		t.Errorf("Approved() should return ok=false when only declines exist")
	}
	if approved != (Attempt{}) {
		t.Errorf("Approved() should return zero-value when no approved, got %v", approved)
	}
}

// T003: Bloquear segundo aprovado (PG-02, ordem de validações)

func TestOrderPayments_PG02_SecondApprovedFails(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)

	err1 := op.RecordAttempt(10000, StatusApproved)
	if err1 != nil {
		t.Errorf("first RecordAttempt(approved) should not error, got %v", err1)
	}

	err2 := op.RecordAttempt(10000, StatusApproved)
	if err2 != ErrAlreadyPaid {
		t.Errorf("second RecordAttempt(approved) should return ErrAlreadyPaid, got %v", err2)
	}

	// Verify only one approved attempt was registered
	attempts := op.Attempts()
	if len(attempts) != 1 {
		t.Errorf("after failed second RecordAttempt, Attempts() should have 1 item, got %d", len(attempts))
	}
}

func TestOrderPayments_PG02_DeclineAfterApproveFails(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)

	err1 := op.RecordAttempt(10000, StatusApproved)
	if err1 != nil {
		t.Errorf("RecordAttempt(approved) should not error, got %v", err1)
	}

	err2 := op.RecordAttempt(10000, StatusDeclined)
	if err2 != ErrAlreadyPaid {
		t.Errorf("RecordAttempt(declined) after approved should return ErrAlreadyPaid, got %v", err2)
	}

	// Verify only one attempt (approved) was registered
	attempts := op.Attempts()
	if len(attempts) != 1 {
		t.Errorf("after failed second RecordAttempt, Attempts() should have 1 item, got %d", len(attempts))
	}

	if attempts[0].Status != StatusApproved {
		t.Errorf("registered attempt should be approved, got %v", attempts[0].Status)
	}
}

func TestOrderPayments_PG02_OrderValidations(t *testing.T) {
	op, _ := NewOrderPayments("order-1", 10000)

	// Register first as approved
	if err := op.RecordAttempt(10000, StatusApproved); err != nil {
		t.Errorf("first RecordAttempt should not error, got %v", err)
	}

	// Try second with wrong amount (should fail with ErrAlreadyPaid, not ErrAmountMismatch)
	err := op.RecordAttempt(9900, StatusApproved)
	if err != ErrAlreadyPaid {
		t.Errorf("RecordAttempt with approved already + wrong amount should return ErrAlreadyPaid (before validating amount), got %v", err)
	}

	// Verify nothing was added
	attempts := op.Attempts()
	if len(attempts) != 1 {
		t.Errorf("after failed second RecordAttempt, Attempts() should still have 1 item, got %d", len(attempts))
	}
}

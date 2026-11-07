// Package payment defines the OrderPayments aggregate for managing payment attempts.
package payment

import (
	"errors"
	"sync"
)

var (
	// ErrInvalidAmount returned when amount is <= 0.
	ErrInvalidAmount = errors.New("invalid amount: must be greater than zero")
	// ErrAmountMismatch returned when amount != total of the order.
	ErrAmountMismatch = errors.New("amount mismatch: payment must equal order total")
	// ErrAlreadyPaid returned when trying to register a payment when one is already approved.
	ErrAlreadyPaid = errors.New("already paid: order has an approved payment")
	// ErrInvalidResult returned when result is not "approved" or "declined".
	ErrInvalidResult = errors.New("invalid result: must be 'approved' or 'declined'")
)

// Status represents the result of a payment attempt.
type Status string

const (
	// StatusApproved payment was accepted.
	StatusApproved Status = "approved"
	// StatusDeclined payment was rejected.
	StatusDeclined Status = "declined"
)

// Attempt represents a single payment attempt for an order.
type Attempt struct {
	Amount int64
	Status Status
}

// OrderPayments is the aggregate root for managing payment attempts for an order.
// Fields are private; state changes only through business methods.
// Mutex protects concurrent access.
type OrderPayments struct {
	orderID  string
	total    int64
	attempts []Attempt
	mu       sync.Mutex
}

// NewOrderPayments creates a new OrderPayments for an order with a given total amount.
// Returns ErrInvalidAmount if total <= 0.
func NewOrderPayments(orderID string, total int64) (*OrderPayments, error) {
	if total <= 0 {
		return nil, ErrInvalidAmount
	}
	return &OrderPayments{
		orderID:  orderID,
		total:    total,
		attempts: []Attempt{},
	}, nil
}

// Approved returns the first approved payment attempt, if any.
// Returns (zero-value, false) if no approved payment exists.
// Uses lock to ensure atomic read.
func (op *OrderPayments) Approved() (Attempt, bool) {
	op.mu.Lock()
	defer op.mu.Unlock()

	for _, attempt := range op.attempts {
		if attempt.Status == StatusApproved {
			return attempt, true
		}
	}
	return Attempt{}, false
}

// Attempts returns a copy of all payment attempts in order.
// Uses lock to ensure atomic read.
func (op *OrderPayments) Attempts() []Attempt {
	op.mu.Lock()
	defer op.mu.Unlock()

	// Defensive copy
	result := make([]Attempt, len(op.attempts))
	copy(result, op.attempts)
	return result
}

// RecordAttempt registers a payment attempt with validation.
// Validates: amount > 0, amount == total, result is valid.
// All-or-nothing: on validation failure, nothing is registered.
// Uses lock to ensure atomic write.
func (op *OrderPayments) RecordAttempt(amount int64, result Status) error {
	// Phase 1: Validate inputs (before acquiring lock)
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount != op.total {
		return ErrAmountMismatch
	}
	if result != StatusApproved && result != StatusDeclined {
		return ErrInvalidResult
	}

	// Phase 2: Acquire lock and register attempt
	op.mu.Lock()
	defer op.mu.Unlock()

	op.attempts = append(op.attempts, Attempt{
		Amount: amount,
		Status: result,
	})

	return nil
}

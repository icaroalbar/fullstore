package stock

import "errors"

var (
	// ErrInsufficientStock returned when trying to reserve without enough available quantity.
	ErrInsufficientStock = errors.New("insufficient stock")
	// ErrInvalidQuantity returned when quantity is <= 0.
	ErrInvalidQuantity = errors.New("invalid quantity: must be greater than zero")
	// ErrProductNotFound returned when product is not found in inventory.
	ErrProductNotFound = errors.New("product not found in inventory")
	// ErrInvalidOperation returned when trying to return or deduct more than reserved.
	ErrInvalidOperation = errors.New("invalid operation: quantity exceeds reserved stock")
)

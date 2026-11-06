// Package order define o agregado Order (Pedido) com validações de itens,
// transições de status e expiração de pedidos pending.
package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Status represents an order's state.
type Status string

const (
	// StatusPending order awaiting payment.
	StatusPending Status = "pending"
	// StatusPaid order payment confirmed.
	StatusPaid Status = "paid"
	// StatusCanceled order canceled by customer or expired.
	StatusCanceled Status = "canceled"
)

var (
	// ErrEmptyOrder returned when creating order with zero items.
	ErrEmptyOrder = errors.New("order must have at least 1 item")
	// ErrTooManyItems returned when order exceeds 20 items.
	ErrTooManyItems = errors.New("order cannot have more than 20 items")
	// ErrInvalidQuantity returned when item quantity is not 1-10.
	ErrInvalidQuantity = errors.New("item quantity must be between 1 and 10")
	// ErrInvalidPrice returned when item price is <= 0.
	ErrInvalidPrice = errors.New("item price must be greater than zero")
	// ErrDuplicateItem returned when two items have the same ProductID.
	ErrDuplicateItem = errors.New("duplicate product in order")
	// ErrInvalidTransition returned on invalid order status change.
	ErrInvalidTransition = errors.New("invalid order status transition")
	// ErrProductUnavailable returned when product is disabled or missing.
	ErrProductUnavailable = errors.New("product is not available")
)

// Item is a value object representing a line item in an order.
type Item struct {
	ProductID string
	Quantity  int32 // 1-10
	Price     int64 // centavos, snapshot do produto
}

// Order is the aggregate root for an order (Pedido) in the domain.
// It manages a collection of items, status transitions, and expiration.
// Fields are private; state changes only through business methods.
type Order struct {
	id        string
	items     []Item
	status    Status
	total     int64 // centavos
	createdAt time.Time
	paidAt    *time.Time
	expiredAt *time.Time
}

// ID returns the order's unique identifier.
func (o *Order) ID() string {
	return o.id
}

// Items returns a copy of the order's items (prevents external mutation).
func (o *Order) Items() []Item {
	itemsCopy := make([]Item, len(o.items))
	copy(itemsCopy, o.items)
	return itemsCopy
}

// Status returns the current status (pending, paid, or canceled).
func (o *Order) Status() Status {
	return o.status
}

// Total returns the order's total in centavos.
func (o *Order) Total() int64 {
	return o.total
}

// CreatedAt returns when the order was created.
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// PaidAt returns when the order was paid (nil if not paid).
func (o *Order) PaidAt() *time.Time {
	return o.paidAt
}

// ExpiredAt returns when the order was expired (nil if not expired).
func (o *Order) ExpiredAt() *time.Time {
	return o.expiredAt
}

// NewOrder cria um novo pedido com validações PE-01, PE-02, e verificações.
// O pedido nasce com status pending. now é injetado para testes.
func NewOrder(items []Item, now time.Time) (*Order, error) {
	// PE-01: validar quantidade de itens (1-20)
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}
	if len(items) > 20 {
		return nil, ErrTooManyItems
	}

	// Copiar slice de entrada (defensive copy)
	itemsCopy := make([]Item, len(items))
	copy(itemsCopy, items)

	// Rastrear ProductIDs vistos para detectar duplicatas
	seen := make(map[string]bool)

	// PE-02: validar quantidade por item (1-10), preço > 0, e duplicatas
	for _, item := range itemsCopy {
		if item.Quantity < 1 || item.Quantity > 10 {
			return nil, ErrInvalidQuantity
		}
		if item.Price <= 0 {
			return nil, ErrInvalidPrice
		}
		if seen[item.ProductID] {
			return nil, ErrDuplicateItem
		}
		seen[item.ProductID] = true
	}

	// PE-03: calcular total (soma de preço × quantidade)
	total := int64(0)
	for _, item := range itemsCopy {
		total += item.Price * int64(item.Quantity)
	}

	return &Order{
		id:        uuid.New().String(),
		items:     itemsCopy,
		status:    StatusPending,
		total:     total,
		createdAt: now,
	}, nil
}

// Pay transiciona status de "pending" para "paid" (PE-04).
func (o *Order) Pay() error {
	// TODO: Implementar
	return errors.New("not implemented")
}

// Cancel transiciona status de "pending" para "canceled" (PE-05).
func (o *Order) Cancel() error {
	// TODO: Implementar
	return errors.New("not implemented")
}

// IsExpired retorna true se o pedido é pending e passou de 30 minutos (PE-06).
func (o *Order) IsExpired() bool {
	// TODO: Implementar
	return false
}

// Expire muda o status para "canceled" se expirado (PE-06).
func (o *Order) Expire() error {
	// TODO: Implementar
	return errors.New("not implemented")
}

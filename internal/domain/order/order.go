// Package order define o agregado Order (Pedido) com validações de itens,
// transições de status e expiração de pedidos pending.
package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrEmptyOrder returned when creating order with zero items.
	ErrEmptyOrder = errors.New("order must have at least 1 item")
	// ErrTooManyItems returned when order exceeds 20 items.
	ErrTooManyItems = errors.New("order cannot have more than 20 items")
	// ErrInvalidQuantity returned when item quantity is not 1-10.
	ErrInvalidQuantity = errors.New("item quantity must be between 1 and 10")
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
	status    string // "pending", "paid", "canceled"
	total     int64  // centavos
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

// Status returns the current status ("pending", "paid", or "canceled").
func (o *Order) Status() string {
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

// NewOrder cria um novo pedido com validações PE-01 e PE-02.
// O pedido nasce com status "pending".
func NewOrder(items []Item) (*Order, error) {
	// PE-01: validar quantidade de itens (1-20)
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}
	if len(items) > 20 {
		return nil, ErrTooManyItems
	}

	// PE-02: validar quantidade por item (1-10)
	for _, item := range items {
		if item.Quantity < 1 || item.Quantity > 10 {
			return nil, ErrInvalidQuantity
		}
	}

	// PE-03: calcular total (soma de preço × quantidade)
	total := int64(0)
	for _, item := range items {
		total += item.Price * int64(item.Quantity)
	}

	return &Order{
		id:        uuid.New().String(),
		items:     items,
		status:    "pending",
		total:     total,
		createdAt: time.Now(),
	}, nil
}

// CalculateTotal retorna a soma de (preço × quantidade) para todos os itens.
func (o *Order) CalculateTotal() int64 {
	// TODO: Implementar
	return 0
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

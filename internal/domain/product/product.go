package product

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var (
	ErrInvalidName  = errors.New("invalid product name: must be 3-100 characters (UTF-8)")
	ErrInvalidPrice = errors.New("invalid product price: must be greater than zero")
)

// Product is an aggregate root for products in the catalog.
// Fields are private; state changes only through business methods.
type Product struct {
	id        string
	name      string
	price     int64 // em centavos, BRL
	disabled  bool
	createdAt time.Time
}

// ID returns the product's unique identifier.
func (p *Product) ID() string {
	return p.id
}

// Name returns the product's name.
func (p *Product) Name() string {
	return p.name
}

// Price returns the product's price in centavos.
func (p *Product) Price() int64 {
	return p.price
}

// Disabled returns whether the product is disabled.
func (p *Product) Disabled() bool {
	return p.disabled
}

// CreatedAt returns when the product was created.
func (p *Product) CreatedAt() time.Time {
	return p.createdAt
}

// NewProduct cria um novo produto com validações PR-01 e PR-02.
// O produto nasce desativado (disabled = true). now é injetado para testes.
func NewProduct(name string, price int64, now time.Time) (*Product, error) {
	// PR-01: validar nome com UTF-8 (3-100 caracteres)
	nameLen := utf8.RuneCountInString(name)
	if nameLen < 3 || nameLen > 100 {
		return nil, ErrInvalidName
	}

	// PR-02: validar preço > 0
	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	return &Product{
		id:        uuid.New().String(),
		name:      name,
		price:     price,
		disabled:  true,
		createdAt: now,
	}, nil
}

// Activate ativa o produto se preço > 0 (PR-03).
func (p *Product) Activate() error {
	if p.price <= 0 {
		return ErrInvalidPrice
	}
	p.disabled = false
	return nil
}

// Deactivate desativa o produto (idempotente, PR-03).
func (p *Product) Deactivate() {
	p.disabled = true
}

// IsAvailable retorna true se o produto não está desativado (PR-04).
func (p *Product) IsAvailable() bool {
	return !p.disabled
}

// ChangePrice altera o preço com validação PR-02/PR-06.
func (p *Product) ChangePrice(newPrice int64) error {
	if newPrice <= 0 {
		return ErrInvalidPrice
	}
	p.price = newPrice
	return nil
}

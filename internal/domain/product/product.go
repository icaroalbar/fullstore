package product

import (
	"errors"
	"time"
	"github.com/google/uuid"
)

var (
	ErrInvalidName  = errors.New("invalid product name: must be 3-100 characters")
	ErrInvalidPrice = errors.New("invalid product price: must be greater than zero")
)

type Product struct {
	ID       string
	Name     string
	Price    int64 // em centavos, BRL
	Disabled bool
	CreatedAt time.Time
}

// NewProduct cria um novo produto com validações PR-01 e PR-02.
// O produto nasce desativado (disabled = true).
func NewProduct(name string, price int64) (*Product, error) {
	if len(name) < 3 || len(name) > 100 {
		return nil, ErrInvalidName
	}
	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	return &Product{
		ID:        uuid.New().String(),
		Name:      name,
		Price:     price,
		Disabled:  true,
		CreatedAt: time.Now(),
	}, nil
}

// Activate ativa o produto se preço > 0 (PR-03).
func (p *Product) Activate() error {
	if p.Price <= 0 {
		return ErrInvalidPrice
	}
	p.Disabled = false
	return nil
}

// Deactivate desativa o produto (idempotente, PR-03).
func (p *Product) Deactivate() {
	p.Disabled = true
}

// IsAvailable retorna true se o produto não está desativado (PR-04).
func (p *Product) IsAvailable() bool {
	return !p.Disabled
}

// ChangePrice altera o preço com validação PR-02/PR-06.
func (p *Product) ChangePrice(newPrice int64) error {
	if newPrice <= 0 {
		return ErrInvalidPrice
	}
	p.Price = newPrice
	return nil
}

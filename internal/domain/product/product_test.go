package product

import (
	"strings"
	"testing"
)

// T001: PR-01 — validação de nome
func TestProduct_PR01_InvalidNameEmpty(t *testing.T) {
	_, err := NewProduct("", 100)
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}
}

func TestProduct_PR01_InvalidNameTooShort(t *testing.T) {
	_, err := NewProduct("TV", 100)
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName for 2-char name, got %v", err)
	}
}

func TestProduct_PR01_InvalidNameTooLong(t *testing.T) {
	longName := ""
	for i := 0; i < 101; i++ {
		longName += "a"
	}
	_, err := NewProduct(longName, 100)
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName for 101-char name, got %v", err)
	}
}

func TestProduct_PR01_ValidNameMinMax(t *testing.T) {
	tests := []string{
		"abc",           // min
		"Notebook Dell", // normal
		strings.Repeat("a", 100), // max
	}
	for _, name := range tests {
		_, err := NewProduct(name, 100)
		if err != nil {
			t.Errorf("valid name %q should not error, got %v", name, err)
		}
	}
}

// T001: PR-02 — validação de preço
func TestProduct_PR02_InvalidPriceZero(t *testing.T) {
	_, err := NewProduct("Produto", 0)
	if err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice for price 0, got %v", err)
	}
}

func TestProduct_PR02_InvalidPriceNegative(t *testing.T) {
	_, err := NewProduct("Produto", -100)
	if err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice for negative price, got %v", err)
	}
}

func TestProduct_PR02_ValidPrice(t *testing.T) {
	_, err := NewProduct("Produto", 100)
	if err != nil {
		t.Errorf("valid price 100 should not error, got %v", err)
	}
}

// T002: PR-03 — ativação e desativação
func TestProduct_PR03_StartsDisabled(t *testing.T) {
	p, _ := NewProduct("Produto", 100)
	if !p.Disabled {
		t.Errorf("new product should be disabled, got Disabled=%v", p.Disabled)
	}
}

func TestProduct_PR03_ActivateWithValidPrice(t *testing.T) {
	p, _ := NewProduct("Produto", 100)
	err := p.Activate()
	if err != nil {
		t.Errorf("activate with valid price should not error, got %v", err)
	}
	if p.Disabled {
		t.Errorf("after activate, Disabled should be false, got %v", p.Disabled)
	}
}

func TestProduct_PR03_ActivateRequiresValidPrice(t *testing.T) {
	// Cria produto com preço válido, depois tenta ativar com preço inválido
	p, _ := NewProduct("Produto", 100)
	p.Price = 0 // torna preço inválido
	err := p.Activate()
	if err != ErrInvalidPrice {
		t.Errorf("activate with invalid price should return ErrInvalidPrice, got %v", err)
	}
}

func TestProduct_PR03_DeactivateActive(t *testing.T) {
	p, _ := NewProduct("Produto", 100)
	p.Activate()
	p.Deactivate()
	if !p.Disabled {
		t.Errorf("after deactivate, Disabled should be true, got %v", p.Disabled)
	}
}

func TestProduct_PR03_DeactivateDisabledIdempotent(t *testing.T) {
	p, _ := NewProduct("Produto", 100)
	// já está disabled
	p.Deactivate()
	if !p.Disabled {
		t.Errorf("deactivate on already disabled should be idempotent, got Disabled=%v", p.Disabled)
	}
}

// T003: PR-04 — disponibilidade no catálogo
func TestProduct_PR04_DisabledNotAvailable(t *testing.T) {
	p, _ := NewProduct("Produto", 100)
	if p.IsAvailable() {
		t.Errorf("disabled product should not be available, got IsAvailable()=true")
	}
}

func TestProduct_PR04_ActiveAvailable(t *testing.T) {
	p, _ := NewProduct("Produto", 100)
	p.Activate()
	if !p.IsAvailable() {
		t.Errorf("active product should be available, got IsAvailable()=false")
	}
}

// T003: PR-05 — imutabilidade de preço no pedido (snapshot)
// Nota: cenários 16-17 envolvem Order aggregate, testados em order_test.go
// Aqui validamos que Product.Price é só leitura em contexto
func TestProduct_PR05_PriceSnapshot1(t *testing.T) {
	p, _ := NewProduct("Produto", 5000) // R$ 50,00 em centavos
	originalPrice := p.Price
	if originalPrice != 5000 {
		t.Errorf("product price should be %d, got %d", 5000, originalPrice)
	}
}

func TestProduct_PR05_PriceSnapshot2(t *testing.T) {
	p, _ := NewProduct("Produto", 5000)
	// Simula mudança de preço no produto
	p.Price = 10000
	// O preço foi mudado no objeto. Em produção, Order teria guardado 5000.
	// Este teste só valida que Price é um campo mutável no domínio,
	// e que quem cria Order é responsável por snapshot.
	if p.Price != 10000 {
		t.Errorf("product price should be updatable, got %d", p.Price)
	}
}

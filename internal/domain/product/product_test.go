package product

import (
	"strings"
	"testing"
	"time"
)

// T001: PR-01 — validação de nome
func TestProduct_PR01_InvalidNameEmpty(t *testing.T) {
	_, err := NewProduct("", 100, time.Now())
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}
}

func TestProduct_PR01_InvalidNameTooShort(t *testing.T) {
	_, err := NewProduct("TV", 100, time.Now())
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName for 2-char name, got %v", err)
	}
}

func TestProduct_PR01_InvalidNameTooLong(t *testing.T) {
	longName := ""
	for i := 0; i < 101; i++ {
		longName += "a"
	}
	_, err := NewProduct(longName, 100, time.Now())
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName for 101-char name, got %v", err)
	}
}

func TestProduct_PR01_ValidNameMinMax(t *testing.T) {
	tests := []string{
		"abc",                    // min
		"Notebook Dell",          // normal
		strings.Repeat("a", 100), // max
	}
	for _, name := range tests {
		_, err := NewProduct(name, 100, time.Now())
		if err != nil {
			t.Errorf("valid name %q should not error, got %v", name, err)
		}
	}
}

// T001: PR-02 — validação de preço
func TestProduct_PR02_InvalidPriceZero(t *testing.T) {
	_, err := NewProduct("Produto", 0, time.Now())
	if err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice for price 0, got %v", err)
	}
}

func TestProduct_PR02_InvalidPriceNegative(t *testing.T) {
	_, err := NewProduct("Produto", -100, time.Now())
	if err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice for negative price, got %v", err)
	}
}

func TestProduct_PR02_ValidPrice(t *testing.T) {
	_, err := NewProduct("Produto", 100, time.Now())
	if err != nil {
		t.Errorf("valid price 100 should not error, got %v", err)
	}
}

// T002: PR-03 — ativação e desativação
func TestProduct_PR03_StartsDisabled(t *testing.T) {
	p, _ := NewProduct("Produto", 100, time.Now())
	if !p.Disabled() {
		t.Errorf("new product should be disabled, got Disabled()=%v", p.Disabled())
	}
}

func TestProduct_PR03_ActivateWithValidPrice(t *testing.T) {
	p, _ := NewProduct("Produto", 100, time.Now())
	err := p.Activate()
	if err != nil {
		t.Errorf("activate with valid price should not error, got %v", err)
	}
	if p.Disabled() {
		t.Errorf("after activate, Disabled() should be false, got %v", p.Disabled())
	}
}

// Nota: TestProduct_PR03_ActivateRequiresValidPrice não é testável agora
// pois campos privados e ChangePrice com validação impedem preço inválido.
// A validação PE-02 é garantida em NewProduct e ChangePrice.

func TestProduct_PR03_DeactivateActive(t *testing.T) {
	p, _ := NewProduct("Produto", 100, time.Now())
	if err := p.Activate(); err != nil {
		t.Errorf("Activate should not error, got %v", err)
	}
	p.Deactivate()
	if !p.Disabled() {
		t.Errorf("after deactivate, Disabled() should be true, got %v", p.Disabled())
	}
}

func TestProduct_PR03_DeactivateDisabledIdempotent(t *testing.T) {
	p, _ := NewProduct("Produto", 100, time.Now())
	// já está disabled
	p.Deactivate()
	if !p.Disabled() {
		t.Errorf("deactivate on already disabled should be idempotent, got Disabled()=%v", p.Disabled())
	}
}

// T003: PR-04 — disponibilidade no catálogo
func TestProduct_PR04_DisabledNotAvailable(t *testing.T) {
	p, _ := NewProduct("Produto", 100, time.Now())
	if p.IsAvailable() {
		t.Errorf("disabled product should not be available, got IsAvailable()=true")
	}
}

func TestProduct_PR04_ActiveAvailable(t *testing.T) {
	p, _ := NewProduct("Produto", 100, time.Now())
	if err := p.Activate(); err != nil {
		t.Errorf("Activate should not error, got %v", err)
	}
	if !p.IsAvailable() {
		t.Errorf("active product should be available, got IsAvailable()=false")
	}
}

// T003: PR-05 — imutabilidade de preço no pedido (snapshot)
// Nota: cenários 16-17 envolvem Order aggregate, testados em order_test.go
// Aqui validamos que Product.Price() retorna valor correto
func TestProduct_PR05_PriceSnapshot1(t *testing.T) {
	p, _ := NewProduct("Produto", 5000, time.Now()) // R$ 50,00 em centavos
	originalPrice := p.Price()
	if originalPrice != 5000 {
		t.Errorf("product price should be %d, got %d", 5000, originalPrice)
	}
}

func TestProduct_PR05_PriceSnapshot2(t *testing.T) {
	p, _ := NewProduct("Produto", 5000, time.Now())
	// Simula mudança de preço no produto via ChangePrice
	if err := p.ChangePrice(10000); err != nil {
		t.Errorf("ChangePrice should not error, got %v", err)
	}
	// O preço foi mudado. Em produção, Order teria guardado 5000.
	// Este teste valida que Price é mutável apenas via ChangePrice.
	if p.Price() != 10000 {
		t.Errorf("product price should be updated, got %d", p.Price())
	}
}

// T003+: PR-06 — alterar preço com validação
func TestProduct_PR06_ChangePriceValid(t *testing.T) {
	p, _ := NewProduct("Produto", 5000, time.Now()) // R$ 50,00
	err := p.ChangePrice(10000)                     // R$ 100,00
	if err != nil {
		t.Errorf("ChangePrice(10000) should not error, got %v", err)
	}
	if p.Price() != 10000 {
		t.Errorf("after ChangePrice(10000), Price() should be 10000, got %d", p.Price())
	}
}

func TestProduct_PR06_ChangePriceZero(t *testing.T) {
	p, _ := NewProduct("Produto", 5000, time.Now())
	err := p.ChangePrice(0)
	if err != ErrInvalidPrice {
		t.Errorf("ChangePrice(0) should return ErrInvalidPrice, got %v", err)
	}
	if p.Price() != 5000 {
		t.Errorf("after ChangePrice(0) error, Price() should remain 5000, got %d", p.Price())
	}
}

func TestProduct_PR06_ChangePriceNegative(t *testing.T) {
	p, _ := NewProduct("Produto", 5000, time.Now())
	err := p.ChangePrice(-100)
	if err != ErrInvalidPrice {
		t.Errorf("ChangePrice(-100) should return ErrInvalidPrice, got %v", err)
	}
	if p.Price() != 5000 {
		t.Errorf("after ChangePrice(-100) error, Price() should remain 5000, got %d", p.Price())
	}
}

// Testes adicionais: UTF-8, injeção relógio, getters, name com acentos
func TestProduct_UTF8Name(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"café", true},                    // 4 runes (valid)
		{"açúcar", true},                  // 6 runes (valid)
		{"naçõescomplexas", true},         // 17 runes (valid)
		{"café com açúcar", true},         // 17 runes, espaço (valid)
		{"ab", false},                     // 2 runes (too short)
		{"ã", false},                      // 1 rune (too short)
		{strings.Repeat("é", 101), false}, // 101 runes (too long)
	}
	for _, tt := range tests {
		_, err := NewProduct(tt.name, 100, time.Now())
		if tt.valid && err != nil {
			t.Errorf("name %q should be valid, got error %v", tt.name, err)
		}
		if !tt.valid && err != ErrInvalidName {
			t.Errorf("name %q should be invalid, got %v", tt.name, err)
		}
	}
}

func TestProduct_InjectedTime(t *testing.T) {
	customTime := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	p, _ := NewProduct("Produto", 100, customTime)
	if p.CreatedAt() != customTime {
		t.Errorf("CreatedAt() should use injected time, got %v, want %v", p.CreatedAt(), customTime)
	}
}

func TestProduct_Getters(t *testing.T) {
	customTime := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	p, _ := NewProduct("TestProduct", 5000, customTime)

	if p.ID() == "" {
		t.Errorf("ID() should not be empty")
	}
	if p.Name() != "TestProduct" {
		t.Errorf("Name() should be TestProduct, got %s", p.Name())
	}
	if p.Price() != 5000 {
		t.Errorf("Price() should be 5000, got %d", p.Price())
	}
	if !p.Disabled() {
		t.Errorf("Disabled() should be true for new product")
	}
	if p.CreatedAt() != customTime {
		t.Errorf("CreatedAt() should be %v, got %v", customTime, p.CreatedAt())
	}
}

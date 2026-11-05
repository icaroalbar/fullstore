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

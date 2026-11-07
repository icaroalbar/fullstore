package ordering

import (
	"errors"
	"testing"
	"time"

	"github.com/icaroalbar/fullstore/internal/domain/order"
	"github.com/icaroalbar/fullstore/internal/domain/stock"
)

// T002: Portas e erros do caso de uso
func TestOrderUseCases_ConstructorAndFakes(t *testing.T) {
	// Fakes
	catalogFake := NewProductCatalogFake()
	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	clockFake := NewClockFake(time.Now())

	// Constructor should succeed
	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)
	if uc == nil {
		t.Errorf("NewOrderUseCases should return non-nil OrderUseCases")
	}
}

// Testes dos fakes (validam comportamentos críticos)
func TestInventoryPortFake_ReserveTudoOuNada(t *testing.T) {
	fake := NewInventoryPortFake()
	fake.SetStock("p1", 10)
	fake.SetStock("p2", 5)

	items := []stock.Item{
		{ProductID: "p1", Quantity: 10},
		{ProductID: "p2", Quantity: 10}, // Não há 10 de p2
	}

	err := fake.Reserve(items)
	if err != stock.ErrInsufficientStock {
		t.Errorf("Reserve should fail, got %v", err)
	}

	// Verificar que nada foi reservado (tudo ou nada)
	s1 := fake.GetStock("p1")
	s2 := fake.GetStock("p2")
	if s1 != 10 || s2 != 5 {
		t.Errorf("Reserve falhou e não deve ter alterado nada; p1=%d (esperado 10), p2=%d (esperado 5)", s1, s2)
	}
}

func TestOrderRepositoryFake_UpdateRetornaErro(t *testing.T) {
	fake := NewOrderRepositoryFake()
	o, _ := order.NewOrder(
		[]order.Item{{ProductID: "p1", Quantity: 1, Price: 100}},
		time.Now(),
	)
	_ = fake.Save(o)

	// Update que retorna erro sem alterar nada
	err := fake.Update(o.ID(), func(_ *order.Order) error {
		return errors.New("test error")
	})

	if err == nil {
		t.Errorf("Update com erro deveria retornar erro")
	}

	// Recuperar a ordem e verificar que status não foi alterado pela função
	recovered, _ := fake.GetForTest(o.ID())
	if recovered.Status() != order.StatusPending {
		t.Errorf("Order status deveria ser pending; status=%v", recovered.Status())
	}
}

func TestOrderRepositoryFake_SaveNaoPanic(t *testing.T) {
	fake := NewOrderRepositoryFake()
	o, _ := order.NewOrder(
		[]order.Item{{ProductID: "p1", Quantity: 1, Price: 100}},
		time.Now(),
	)

	// Save não deve dar panic
	err := fake.Save(o)
	if err != nil {
		t.Errorf("Save should not error, got %v", err)
	}
}

func TestInventoryPortFake_ReturnNaoPanic(t *testing.T) {
	fake := NewInventoryPortFake()
	fake.SetStock("p1", 10)

	// Return não deve dar panic
	items := []stock.Item{{ProductID: "p1", Quantity: 1}}
	err := fake.Return(items)
	if err != nil {
		t.Errorf("Return should not error, got %v", err)
	}
}

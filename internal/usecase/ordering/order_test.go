package ordering

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/icaroalbar/fullstore/internal/domain/order"
	"github.com/icaroalbar/fullstore/internal/domain/product"
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

// T004: CancelarPedido validações e transição (cenários 06-08b)

func TestCancelarPedido_06_ES03_CancelaPending(t *testing.T) {
	// Setup: pedido pending com estoque reservado
	catalogFake := NewProductCatalogFake()
	p1, _ := product.NewProduct("prod1", 10000, time.Now())
	if err := p1.Activate(); err != nil {
		t.Fatalf("Activate should not fail: %v", err)
	}
	catalogFake.SetProductByID("prod1", p1)

	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 100)
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	// Criar pedido (CriarPedido)
	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	o, err := uc.CriarPedido(items)
	if err != nil {
		t.Fatalf("CriarPedido failed: %v", err)
	}
	orderID := o.ID()

	// Estoque deve estar reservado: 100 - 5 = 95
	if inventoryFake.GetStock("prod1") != 95 {
		t.Fatalf("Stock deve estar reservado em 95, got %d", inventoryFake.GetStock("prod1"))
	}

	// Cancelar pedido
	err = uc.CancelarPedido(orderID)

	if err != nil {
		t.Errorf("CancelarPedido com pending deve suceder, got %v", err)
	}

	// Verificar que estoque foi devolvido: 95 + 5 = 100
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque deve ter sido devolvido; esperado 100, got %d", inventoryFake.GetStock("prod1"))
	}

	// Verificar status do pedido
	recovered, _ := orderRepoFake.GetForTest(orderID)
	if recovered.Status() != order.StatusCanceled {
		t.Errorf("Pedido deve estar canceled, got %v", recovered.Status())
	}
}

func TestCancelarPedido_07_PE05_CancelaPago(t *testing.T) {
	// Setup: pedido paid
	catalogFake := NewProductCatalogFake()
	p1, _ := product.NewProduct("prod1", 10000, time.Now())
	if err := p1.Activate(); err != nil {
		t.Fatalf("Activate should not fail: %v", err)
	}
	catalogFake.SetProductByID("prod1", p1)

	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 100)
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	// Criar pedido
	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	o, err := uc.CriarPedido(items)
	if err != nil {
		t.Fatalf("CriarPedido failed: %v", err)
	}
	orderID := o.ID()

	// Simular transição para paid via Update (apenas para teste)
	_ = orderRepoFake.Update(orderID, func(ord *order.Order) error {
		return ord.Pay(clockFake.Now())
	})

	// Estoque deve estar reservado: 100 - 5 = 95
	stockBefore := inventoryFake.GetStock("prod1")

	// Tentar cancelar pedido pago
	err = uc.CancelarPedido(orderID)

	if err != order.ErrInvalidTransition {
		t.Errorf("CancelarPedido com paid deve retornar ErrInvalidTransition, got %v", err)
	}

	// Verificar que estoque NÃO foi devolvido
	if inventoryFake.GetStock("prod1") != stockBefore {
		t.Errorf("Estoque não deve ter sido alterado, esperado %d, got %d", stockBefore, inventoryFake.GetStock("prod1"))
	}
}

func TestCancelarPedido_08_CancelaCancelado(t *testing.T) {
	// Setup: pedido canceled
	catalogFake := NewProductCatalogFake()
	p1, _ := product.NewProduct("prod1", 10000, time.Now())
	if err := p1.Activate(); err != nil {
		t.Fatalf("Activate should not fail: %v", err)
	}
	catalogFake.SetProductByID("prod1", p1)

	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 100)
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	// Criar pedido
	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	o, err := uc.CriarPedido(items)
	if err != nil {
		t.Fatalf("CriarPedido failed: %v", err)
	}
	orderID := o.ID()

	// Cancelar uma vez
	err = uc.CancelarPedido(orderID)
	if err != nil {
		t.Fatalf("Primeiro CancelarPedido falhou: %v", err)
	}

	// Estoque deve estar em 100 (devolvido)
	if inventoryFake.GetStock("prod1") != 100 {
		t.Fatalf("Estoque deve estar em 100 após primeiro cancelamento, got %d", inventoryFake.GetStock("prod1"))
	}

	// Tentar cancelar de novo
	err = uc.CancelarPedido(orderID)

	if err != order.ErrInvalidTransition {
		t.Errorf("CancelarPedido em canceled deve retornar ErrInvalidTransition, got %v", err)
	}

	// Verificar que estoque não foi devolvido novamente (teria voltado a 105)
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque não deve ser devolvido de novo; esperado 100, got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestCancelarPedido_08b_PedidoInexistente(t *testing.T) {
	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 100)
	clockFake := NewClockFake(time.Now())

	catalogFake := NewProductCatalogFake()
	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	// Tentar cancelar pedido inexistente
	err := uc.CancelarPedido("pedido-inexistente")

	if err != ErrOrderNotFound {
		t.Errorf("CancelarPedido com pedido inexistente deve retornar ErrOrderNotFound, got %v", err)
	}

	// Verificar que estoque não foi alterado
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque não deve ter sido alterado, got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestCancelarPedido_UpdateFaila(t *testing.T) {
	orderRepoFake := NewOrderRepositoryFake()
	orderRepoFake.SetError(errors.New("update error"))
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 100)
	clockFake := NewClockFake(time.Now())

	catalogFake := NewProductCatalogFake()
	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	// Tentar cancelar quando Update falha
	err := uc.CancelarPedido("any-order-id")

	if err == nil {
		t.Errorf("CancelarPedido com Update falhando deve retornar erro")
	}

	// Verificar que Return não foi chamado
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque não deve ter sido alterado (Return não chamado), got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestCriarPedido_04c_SaveEReturnFalham(t *testing.T) {
	catalogFake := NewProductCatalogFake()
	p1, _ := product.NewProduct("prod1", 10000, time.Now())
	if err := p1.Activate(); err != nil {
		t.Fatalf("Activate should not fail: %v", err)
	}
	catalogFake.SetProductByID("prod1", p1)

	orderRepoFake := NewOrderRepositoryFake()
	orderRepoFake.SetError(errors.New("save error"))
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 100)
	inventoryFake.SetReturnError(errors.New("return error"))
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	o, err := uc.CriarPedido(items)

	if err == nil {
		t.Errorf("CriarPedido com Save e Return falhando deve retornar erro")
	}
	if o != nil {
		t.Errorf("CriarPedido com erro não deve retornar order")
	}

	// Verificar que errors.Join retorna ambos os erros
	if !errors.Is(err, errors.New("save error")) {
		// Note: errors.Is não funciona com strings diretas de errors.New,
		// mas devemos verificar que o erro do Save está no joined error
		errStr := err.Error()
		if !strings.Contains(errStr, "save error") {
			t.Errorf("Erro joinado deve conter 'save error', got %v", errStr)
		}
	}
	if !strings.Contains(err.Error(), "return error") {
		t.Errorf("Erro joinado deve conter 'return error', got %v", err)
	}
}

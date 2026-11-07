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

// T003: CriarPedido validações e reserva (cenários 01-05)

func TestCriarPedido_01_ES01_ReservaEstoque(t *testing.T) {
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

	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	o, err := uc.CriarPedido(items)

	if err != nil {
		t.Errorf("CriarPedido com produtos ativos deve suceder, got %v", err)
	}
	if o == nil {
		t.Errorf("CriarPedido deve retornar order não-nil")
	}

	// Verificar reserva: 100 - 5 = 95 disponíveis
	if inventoryFake.GetStock("prod1") != 95 {
		t.Errorf("Estoque deve ter sido reservado; esperado 95, got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestCriarPedido_02_PR04_ProdutoDesativado(t *testing.T) {
	catalogFake := NewProductCatalogFake()
	p1, _ := product.NewProduct("prod1", 10000, time.Now())
	// Não ativar (status desativado)
	catalogFake.SetProductByID("prod1", p1)

	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 100)
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	o, err := uc.CriarPedido(items)

	if err != ErrProductUnavailable {
		t.Errorf("CriarPedido com produto desativado deve retornar ErrProductUnavailable, got %v", err)
	}
	if o != nil {
		t.Errorf("CriarPedido com erro não deve retornar order, got %v", o)
	}

	// Verificar que nada foi reservado
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque não deve ter sido alterado, got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestCriarPedido_03_ES02_SemEstoque(t *testing.T) {
	catalogFake := NewProductCatalogFake()
	p1, _ := product.NewProduct("prod1", 10000, time.Now())
	if err := p1.Activate(); err != nil {
		t.Fatalf("Activate should not fail: %v", err)
	}
	p2, _ := product.NewProduct("prod2", 5000, time.Now())
	if err := p2.Activate(); err != nil {
		t.Fatalf("Activate should not fail: %v", err)
	}
	catalogFake.SetProductByID("prod1", p1)
	catalogFake.SetProductByID("prod2", p2)

	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 10)
	inventoryFake.SetStock("prod2", 3) // Não há 5 de prod2
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	items := []ItemInput{{ProductID: "prod1", Quantity: 5}, {ProductID: "prod2", Quantity: 5}}
	o, err := uc.CriarPedido(items)

	if err != stock.ErrInsufficientStock {
		t.Errorf("CriarPedido sem estoque deve retornar ErrInsufficientStock, got %v", err)
	}
	if o != nil {
		t.Errorf("CriarPedido com erro não deve retornar order")
	}

	// Verificar tudo-ou-nada: nada foi reservado
	if inventoryFake.GetStock("prod1") != 10 {
		t.Errorf("prod1 não deve ter sido alterado (tudo-ou-nada), got %d", inventoryFake.GetStock("prod1"))
	}
	if inventoryFake.GetStock("prod2") != 3 {
		t.Errorf("prod2 não deve ter sido alterado (tudo-ou-nada), got %d", inventoryFake.GetStock("prod2"))
	}
}

func TestCriarPedido_04_SaveFaila(t *testing.T) {
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
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	o, err := uc.CriarPedido(items)

	if err == nil {
		t.Errorf("CriarPedido com Save falho deve retornar erro")
	}
	if o != nil {
		t.Errorf("CriarPedido com erro não deve retornar order")
	}

	// Verificar que Return foi chamado (estoque restaurado)
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque deve ter sido devolvido (compensação), got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestCriarPedido_04b_PedidoInvalido(t *testing.T) {
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

	// Pedido com quantidade inválida (0)
	items := []ItemInput{{ProductID: "prod1", Quantity: 0}}
	o, err := uc.CriarPedido(items)

	if err == nil {
		t.Errorf("CriarPedido com pedido inválido deve retornar erro")
	}
	if o != nil {
		t.Errorf("CriarPedido com erro não deve retornar order")
	}

	// Verificar que nada foi reservado
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque não deve ter sido alterado, got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestCriarPedido_05_ProdutoInexistente(t *testing.T) {
	catalogFake := NewProductCatalogFake()
	// Não adicionar produto prod1
	catalogFake.SetError(ErrProductNotFound) // Simula erro de busca

	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	o, err := uc.CriarPedido(items)

	if err != ErrProductNotFound {
		t.Errorf("CriarPedido com produto inexistente deve retornar ErrProductNotFound, got %v", err)
	}
	if o != nil {
		t.Errorf("CriarPedido com erro não deve retornar order")
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

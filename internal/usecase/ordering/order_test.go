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

	errSave := errors.New("save error")
	errReturn := errors.New("return error")

	orderRepoFake := NewOrderRepositoryFake()
	orderRepoFake.SetError(errSave)
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 100)
	inventoryFake.SetReturnError(errReturn)
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
	if !errors.Is(err, errSave) {
		t.Errorf("Erro joinado deve conter errSave, got %v", err)
	}
	if !errors.Is(err, errReturn) {
		t.Errorf("Erro joinado deve conter errReturn, got %v", err)
	}
}

// T005: ExpirarPedidos com Update + clock (cenários 09-10)

func TestExpirarPedidos_09_PE06_ExpiraAntigosDevolvePending(t *testing.T) {
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

	// Avanço relógio em 31 minutos (>30 min)
	clockFake.SetNow(time.Now().Add(31 * time.Minute))

	// Rodar expiração
	err = uc.ExpirarPedidos()

	if err != nil {
		t.Errorf("ExpirarPedidos deve suceder, got %v", err)
	}

	// Verificar que pedido foi cancelado com expiredAt
	recovered, _ := orderRepoFake.GetForTest(orderID)
	if recovered.Status() != order.StatusCanceled {
		t.Errorf("Pedido deve estar canceled, got %v", recovered.Status())
	}
	if recovered.ExpiredAt() == nil {
		t.Errorf("ExpiredAt deve estar setado")
	}

	// Verificar que estoque foi devolvido
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque deve ter sido devolvido; esperado 100, got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestExpirarPedidos_09b_NaoExpiraJaCancelado(t *testing.T) {
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

	// Simular que o pedido foi cancelado via Update (fora do escopo de ExpirarPedidos)
	_ = orderRepoFake.Update(orderID, func(ord *order.Order) error {
		// Simular devolução manual para estado paid
		_ = inventoryFake.Return([]stock.Item{{ProductID: "prod1", Quantity: 5}})
		return ord.Pay(clockFake.Now())
	})

	// Estoque deve estar em 100 (devolvido)
	if inventoryFake.GetStock("prod1") != 100 {
		t.Fatalf("Estoque deve estar em 100, got %d", inventoryFake.GetStock("prod1"))
	}

	// Avanço relógio em 31 minutos
	clockFake.SetNow(time.Now().Add(31 * time.Minute))

	// Rodar expiração (ListPending não inclui paid, então este pedido não será processado)
	err = uc.ExpirarPedidos()

	// Não deve retornar erro
	if err != nil {
		t.Errorf("ExpirarPedidos deve suceder, got %v", err)
	}

	// Verificar que estoque não foi alterado
	if inventoryFake.GetStock("prod1") != 100 {
		t.Errorf("Estoque não deve ser alterado; esperado 100, got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestExpirarPedidos_10_ReturnFailaAcumulaErro(t *testing.T) {
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

	// Configurar Return para falhar
	errReturn := errors.New("return error")
	inventoryFake.SetReturnError(errReturn)

	// Avanço relógio em 31 minutos
	clockFake.SetNow(time.Now().Add(31 * time.Minute))

	// Rodar expiração
	err = uc.ExpirarPedidos()

	if err == nil {
		t.Errorf("ExpirarPedidos deve retornar erro quando Return falha")
	}

	// Verificar que erro contém o ID do pedido
	if !strings.Contains(err.Error(), orderID) {
		t.Errorf("Erro deve conter ID do pedido %s, got %v", orderID, err)
	}

	// Verificar que errors.Is acha o erro do Return
	if !errors.Is(err, errReturn) {
		t.Errorf("Erro deve conter errReturn, got %v", err)
	}

	// Estoque não deve ter sido devolvido (Return falhou)
	if inventoryFake.GetStock("prod1") != 95 {
		t.Errorf("Estoque deve estar em 95 (Return falhou); got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestExpirarPedidos_MistoNovoEAntigo(t *testing.T) {
	catalogFake := NewProductCatalogFake()
	p1, _ := product.NewProduct("prod1", 10000, time.Now())
	if err := p1.Activate(); err != nil {
		t.Fatalf("Activate should not fail: %v", err)
	}
	catalogFake.SetProductByID("prod1", p1)

	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 200)
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	// Criar pedido antigo (será expirado)
	clockFake.SetNow(time.Now().Add(-31 * time.Minute))
	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	oAntigo, _ := uc.CriarPedido(items)
	idAntigo := oAntigo.ID()

	// Criar pedido novo (não será expirado, pois clock foi resetado)
	clockFake.SetNow(time.Now())
	oNovo, _ := uc.CriarPedido(items)
	idNovo := oNovo.ID()

	// Rodar expiração com o relógio atual
	err := uc.ExpirarPedidos()

	if err != nil {
		t.Errorf("ExpirarPedidos com misto deve suceder, got %v", err)
	}

	// Verificar que pedido antigo foi expirado
	recoveredAntigo, _ := orderRepoFake.GetForTest(idAntigo)
	if recoveredAntigo.Status() != order.StatusCanceled {
		t.Errorf("Pedido antigo deve estar canceled, got %v", recoveredAntigo.Status())
	}

	// Verificar que pedido novo NÃO foi expirado
	recoveredNovo, _ := orderRepoFake.GetForTest(idNovo)
	if recoveredNovo.Status() != order.StatusPending {
		t.Errorf("Pedido novo deve estar pending (not expired), got %v", recoveredNovo.Status())
	}

	// Verificar estoque: 200 - 5 (antigo) = 195 (antigo) → expirar devolve 5 → 195; - 5 (novo) = 190 após novo
	// Mas antigo expira ANTES de novo ser criado no relógio, então:
	// Inicial 200 → antigo -5 = 195 → novo -5 = 190
	// Expiração: antigo retorna +5 = 195, novo ErrNotExpired (não toca)
	if inventoryFake.GetStock("prod1") != 195 {
		t.Errorf("Estoque deve estar em 195 (antigo devolvido); got %d", inventoryFake.GetStock("prod1"))
	}
}

// T006: Teste do problema conhecido (Return falha após Cancel bem-sucedido)

func TestCancelarPedido_ReturnFails_PropagaErroComID(t *testing.T) {
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

	// Estoque reservado: 100 - 5 = 95
	if inventoryFake.GetStock("prod1") != 95 {
		t.Fatalf("Stock deve estar em 95, got %d", inventoryFake.GetStock("prod1"))
	}

	// Configurar Return para falhar
	errReturn := errors.New("return error")
	inventoryFake.SetReturnError(errReturn)

	// Cancelar pedido (Update funciona, Return falha)
	err = uc.CancelarPedido(orderID)

	if err == nil {
		t.Errorf("CancelarPedido com Return falhando deve retornar erro")
	}

	// Verificar que erro contém o ID do pedido
	if !strings.Contains(err.Error(), orderID) {
		t.Errorf("Erro deve conter ID do pedido %s, got %v", orderID, err)
	}

	// Verificar que errors.Is encontra o erro do Return
	if !errors.Is(err, errReturn) {
		t.Errorf("Erro deve conter errReturn, got %v", err)
	}

	// Verificar que o pedido foi cancelado mesmo assim
	recovered, _ := orderRepoFake.GetForTest(orderID)
	if recovered.Status() != order.StatusCanceled {
		t.Errorf("Pedido deve estar canceled mesmo com Return falhando, got %v", recovered.Status())
	}

	// Verificar que estoque NÃO foi devolvido (Return falhou)
	if inventoryFake.GetStock("prod1") != 95 {
		t.Errorf("Estoque deve estar preso em 95 (Return falhou); got %d", inventoryFake.GetStock("prod1"))
	}
}

// T005 polimento: ListPending falha e Update falha com erro real

func TestExpirarPedidos_ListPendingFalha(t *testing.T) {
	catalogFake := NewProductCatalogFake()
	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	// Configurar ListPending para falhar
	orderRepoFake.SetError(errors.New("listpending error"))

	// Rodar expiração
	err := uc.ExpirarPedidos()

	if err == nil {
		t.Errorf("ExpirarPedidos com ListPending falhando deve retornar erro")
	}

	// Nenhum Return deve ter sido chamado
	if inventoryFake.GetStock("prod1") != 0 {
		t.Errorf("Estoque não deve ter mudado; got %d", inventoryFake.GetStock("prod1"))
	}
}

func TestExpirarPedidos_UpdateFalhaComErroReal(t *testing.T) {
	catalogFake := NewProductCatalogFake()
	p1, _ := product.NewProduct("prod1", 10000, time.Now())
	if err := p1.Activate(); err != nil {
		t.Fatalf("Activate should not fail: %v", err)
	}
	catalogFake.SetProductByID("prod1", p1)

	orderRepoFake := NewOrderRepositoryFake()
	inventoryFake := NewInventoryPortFake()
	inventoryFake.SetStock("prod1", 200)
	clockFake := NewClockFake(time.Now())

	uc := NewOrderUseCases(catalogFake, orderRepoFake, inventoryFake, clockFake)

	// Criar um pedido antigo
	clockFake.SetNow(time.Now().Add(-31 * time.Minute))
	items := []ItemInput{{ProductID: "prod1", Quantity: 5}}
	_, _ = uc.CriarPedido(items)

	// Voltar ao relógio atual e configurar Update para falhar
	clockFake.SetNow(time.Now())
	orderRepoFake.SetError(errors.New("update error"))

	// Rodar expiração
	err := uc.ExpirarPedidos()

	if err == nil {
		t.Errorf("ExpirarPedidos com Update falhando deve retornar erro")
	}

	// Verificar que nenhum Return foi chamado (Update falhou)
	if inventoryFake.GetStock("prod1") != 195 {
		t.Errorf("Estoque não deve ter devolvido (Update falhou); got %d", inventoryFake.GetStock("prod1"))
	}
}

# Tarefas 005: Casos de Uso — Pedido

## T001: Mudar Order.Expire para retornar ErrInvalidTransition (Spec 002)

**O quê:** Pré-requisito de Spec 005. Order.Expire (em domain/order/order.go) era idempotente. Mudar para retornar `ErrInvalidTransition` se o pedido não estiver pending. Impede Return duplicado em ExpirarPedidos.

**Critério de aceite:**
- [ ] `TestOrder_PE06_ExpireNotPendingReturnsError` — Expire em pedido não pending retorna ErrInvalidTransition (não nil)
- [ ] Spec 002 (docs/specs/002-pedido/spec.md) atualizada: PE-06 documenta novo comportamento
- [ ] `make check` passa

**Arquivo:** `internal/domain/order/order.go` + `order_test.go` + `docs/specs/002-pedido/spec.md`

---

## T002: Portas e erros do caso de uso (ES-01, ES-02, ES-03)

**O quê:** Struct `OrderUseCases`, construtor, injeção de 4 portas (ProductCatalog, OrderRepository, InventoryPort, Clock). Erros em var block. Adaptadores falsos (fakes) nos testes.

**Critério de aceite:**
- [ ] `internal/usecase/ordering/order.go` criado com construtor
- [ ] Portas: ProductCatalog, OrderRepository (*order.Order), InventoryPort, Clock
- [ ] Erros em var: ErrProductNotFound, ErrProductUnavailable, ErrOrderNotFound
- [ ] ItemInput struct: ProductID, Quantity
- [ ] Fakes completos para todas as portas (ProductCatalogFake, OrderRepositoryFake, InventoryPortFake, ClockFake)
- [ ] `make check` passa

**Arquivo:** `internal/usecase/ordering/order.go` + `order_test.go`

---

## T003: CriarPedido validações e reserva (Cenários 01-05, ES-01, ES-02, PR-04, PE-01-03)

**O quê:** Implementar `CriarPedido([]ItemInput) (*order.Order, error)` com ordem: 1) buscar produtos (catálogo), 2) NewOrder memória (preços buscados), 3) Reserve, 4) Save. Se (4) falha, chamar Return (compensação). Propagar stock.ErrInsufficientStock (não redeclarar).

**Critério de aceite:**
- [ ] `TestCriarPedido_01_ES01_ReservaEstoque` — produtos ativos, estoque reservado
- [ ] `TestCriarPedido_02_PR04_ProdutoDesativado` — ErrProductUnavailable, sem Reserve
- [ ] `TestCriarPedido_03_ES02_SemEstoque` — stock.ErrInsufficientStock, tudo-ou-nada
- [ ] `TestCriarPedido_04_SaveFaila` — Save falha, Return chamado, erro propagado
- [ ] `TestCriarPedido_04b_PedidoInvalido` — NewOrder inválido (PE-01-03), sem Reserve
- [ ] `TestCriarPedido_05_ProdutoInexistente` — ErrProductNotFound, sem efeito
- [ ] `make check` passa

**Arquivo:** `internal/usecase/ordering/order.go` + `order_test.go`

---

## T004: CancelarPedido validações e devolução (Cenários 06-08b, ES-03, PE-05)

**O quê:** Implementar `CancelarPedido(orderID string) error` com ordem: 1) Get pedido, 2) Cancel (memória), 3) Save, 4) Return. Se Save falha, não chamar Return (sem efeito colateral órfão).

**Critério de aceite:**
- [ ] `TestCancelarPedido_06_ES03_CancelaPending` — pending → canceled, Return chamado, estoque devolvido
- [ ] `TestCancelarPedido_07_PE05_CancelaPago` — pago → ErrInvalidTransition (de Order.Cancel), sem Return
- [ ] `TestCancelarPedido_08_CancelaCancelado` — já canceled → ErrInvalidTransition, sem Return duplicado
- [ ] `TestCancelarPedido_08b_PedidoInexistente` — inexistente → ErrOrderNotFound (antes de Cancel)
- [ ] `TestCancelarPedido_SaveFaila` — Save falha após Cancel, Return não chamado
- [ ] `make check` passa

**Arquivo:** `internal/usecase/ordering/order.go` + `order_test.go`

---

## T005: ExpirarPedidos com clock (Cenários 09-10, PE-06, ES-03)

**O quê:** Implementar `ExpirarPedidos() error` que: 1) lista pending, 2) para cada: Get, Expire(clock.Now()), Save, Return. Se Return falha, acumula erro com ID do pedido. Propaga errors.Join de todos os erros. Sem parar no primeiro erro.

**Critério de aceite:**
- [ ] `TestExpirarPedidos_09_PE06_ExpiraAntigosDevolvePending` — >30 min → expired (canceled com expiredAt), Return chamado
- [ ] `TestExpirarPedidos_09b_NaoExpiraJaCancelado` — já cancelado → ErrInvalidTransition (de Expire), sem Return
- [ ] `TestExpirarPedidos_10_ReturnFailaAcumulaErro` — Return falha, erro acumulado com ID, demais processados, errors.Join propagado
- [ ] `make check` passa

**Arquivo:** `internal/usecase/ordering/order.go` + `order_test.go`

---

## T006: Teste do problema conhecido (Return falha após Cancel/Expire)

**O quê:** Cenários 04 (Save falha em CriarPedido) e 10 (Return falha em ExpirarPedidos) já testam compensações bem-sucedidas. Aqui: demonstrar o problema: se Order.Cancel/Expire der OK, Save OK, mas Return falha, a reserva fica presa. Erro deve ser propagado com ID do pedido.

**Critério de aceite:**
- [ ] `TestCancelarPedido_ReturnFails_PropagaErroComID` — Cancel OK, Save OK, Return falha; erro propagado contém ID do pedido
- [ ] `TestExpirarPedidos_ReturnFailaAcumulado` — Expire OK, Save OK, Return falha; erro acumulado em errors.Join com ID
- [ ] Comentário no código: "Estoque fica preso; operacional rastreia pelo ID e manualmente compensa"
- [ ] `make check` passa

**Arquivo:** `internal/usecase/ordering/order_test.go`

---

## T007: Concorrência — dois cancelamentos simultâneos (ES-05 pattern)

**O quê:** Tarefa extra: verificar que dois CancelarPedido simultâneos no mesmo pedido não causam Return duplicado. Um ganha a transição (pending → canceled), outro falha com ErrInvalidTransition. Apenas um Return chamado. Teste com -race.

**Critério de aceite:**
- [ ] `TestCancelarPedido_RaceCondition_DoisCancelamentosSimultaneos` — 50 goroutines (canal start), todos tentam cancelar o mesmo pending; 1 sucesso (Return 1×), 49 falham com ErrInvalidTransition (sem Return)
- [ ] `go test -race ./internal/usecase/ordering` passa sem data races
- [ ] `make check` passa

**Arquivo:** `internal/usecase/ordering/order_test.go`

---

## T008: Cobertura >= 80% + status implementada

**O quê:** Verificar cobertura, atualizar spec.md Status: implementada.

**Critério de aceite:**
- [ ] `go test -cover ./internal/usecase/ordering` mostra >= 80%
- [ ] Spec 005: Status mudado para `implementada`
- [ ] `make check` passa

**Arquivo:** `docs/specs/005-casos-de-uso-pedido/spec.md`

---

## Ordem de execução

1. **T001** — Pré-requisito: Order.Expire mudança (Spec 002)
2. **T002** — Setup: portas, erros, construtor
3. **T003** — CriarPedido: ordem, validações, compensação
4. **T004** — CancelarPedido: transição, compensação segura
5. **T005** — ExpirarPedidos: clock, errors.Join, sem parar no erro
6. **T006** — Teste do problema conhecido (Return falha com ID)
7. **T007** — Concorrência: -race com dois cancelamentos
8. **T008** — Cobertura + spec implementada

Cada tarefa finaliza com `make check` verde e diff explicado.

---

## Regras cobertas

✅ ES-01 (CriarPedido reserva estoque) — T003
✅ ES-02 (Sem estoque suficiente: tudo-ou-nada) — T003
✅ ES-03 (Cancelar/expirar devolve) — T004, T005
✅ PR-04 (Produto desativado) — T003
✅ PE-06 (Expira em 30 min) — T005

## Pacote

**internal/usecase/ordering:** 3 casos de uso (CriarPedido, CancelarPedido, ExpirarPedidos)
- Alias: `order "internal/domain/order"` (evita conflito com pacote usecase/order que não existe)

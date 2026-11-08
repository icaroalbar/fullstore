# Tarefas 005: Casos de Uso — Pedido

## T001: Mudar Order.Expire para retornar ErrInvalidTransition (Spec 002)

**O quê:** Pré-requisito de Spec 005. Order.Expire (em domain/order/order.go) era idempotente. Mudar para retornar `ErrInvalidTransition` se o pedido não estiver pending. Impede Return duplicado em ExpirarPedidos.

**Critério de aceite:**
- [x] `TestOrder_PE06_ExpireNotPendingReturnsError` — Expire em pedido não pending retorna ErrInvalidTransition (não nil)
- [x] `TestOrder_PE06_ExpireOnPaidOrderReturnsError` — Expire em pedido paid retorna ErrInvalidTransition
- [x] Spec 002 (docs/specs/002-pedido/spec.md) atualizada: PE-06 documenta novo comportamento
- [x] `make check` passa

**Arquivo:** `internal/domain/order/order.go` + `order_test.go` + `docs/specs/002-pedido/spec.md`

---

## T002: Portas e erros do caso de uso (ES-01, ES-02, ES-03)

**O quê:** Struct `OrderUseCases`, construtor, injeção de 4 portas (ProductCatalog, OrderRepository, InventoryPort, Clock). Erros em var block. Adaptadores falsos (fakes) nos testes com thread-safety.

**Critério de aceite:**
- [x] `internal/usecase/ordering/order.go` criado com construtor
- [x] Portas: ProductCatalog (*product.Product), OrderRepository (*order.Order), InventoryPort, Clock
- [x] Erros em var: ErrProductNotFound, ErrProductUnavailable, ErrOrderNotFound
- [x] ItemInput struct: ProductID, Quantity (int32)
- [x] Fakes em `fakes_test.go`: ProductCatalogFake, OrderRepositoryFake, InventoryPortFake, ClockFake com mutex
- [x] Testes dos fakes: Reserve tudo-ou-nada, Update transacional, Save/Return sem panic
- [x] `make check` passa (100% cobertura)

**Arquivo:** `internal/usecase/ordering/order.go` + `order_test.go` + `fakes_test.go`

---

## T003: CriarPedido validações e reserva (Cenários 01-05, ES-01, ES-02, PR-04, PE-01-03)

**O quê:** Implementar `CriarPedido([]ItemInput) (*order.Order, error)` com ordem: 1) buscar produtos (catálogo), 2) NewOrder memória (preços buscados), 3) Reserve, 4) Save. Se (4) falha, chamar Return (compensação). Propagar stock.ErrInsufficientStock (não redeclarar).

**Critério de aceite:**
- [x] `TestCriarPedido_01_ES01_ReservaEstoque` — produtos ativos, estoque reservado
- [x] `TestCriarPedido_02_PR04_ProdutoDesativado` — ErrProductUnavailable, sem Reserve
- [x] `TestCriarPedido_03_ES02_SemEstoque` — stock.ErrInsufficientStock, tudo-ou-nada
- [x] `TestCriarPedido_04_SaveFaila` — Save falha, Return chamado, erro propagado
- [x] `TestCriarPedido_04b_PedidoInvalido` — NewOrder inválido (PE-01-03), sem Reserve
- [x] `TestCriarPedido_04c_SaveEReturnFalham` — Save e Return falham, errors.Join(saveErr, returnErr)
- [x] `TestCriarPedido_05_ProdutoInexistente` — ErrProductNotFound, sem efeito
- [x] `make check` passa (com -race)

**Arquivo:** `internal/usecase/ordering/order.go` + `order_test.go`

---

## T004: CancelarPedido com Update atômico (Cenários 06-08b, ES-03, PE-05)

**O quê:** Implementar `CancelarPedido(orderID string) error` com Update atômico: `repo.Update(orderID, func(o *order.Order) error { return o.Cancel() })`. Depois: Return. Se Update falha (ErrOrderNotFound, ErrInvalidTransition), não chamar Return.

**Critério de aceite:**
- [x] `TestCancelarPedido_06_ES03_CancelaPending` — Update ok, Return chamado, estoque devolvido
- [x] `TestCancelarPedido_07_PE05_CancelaPago` — Update retorna ErrInvalidTransition (de o.Cancel), sem Return
- [x] `TestCancelarPedido_08_CancelaCancelado` — Update retorna ErrInvalidTransition, sem Return duplicado
- [x] `TestCancelarPedido_08b_PedidoInexistente` — Update retorna ErrOrderNotFound, sem Return
- [x] `TestCancelarPedido_UpdateFaila` — Update falha, Return não chamado
- [x] `make check` passa (100% coverage, 6 testes CancelarPedido)

**Arquivo:** `internal/usecase/ordering/order.go` + `order_test.go`

---

## T005: ExpirarPedidos com Update + clock (Cenários 09-10, PE-06, ES-03)

**O quê:** Implementar `ExpirarPedidos() error` que: 1) lista pending (cópias), 2) para cada: `Update(id, func(o *order.Order) error { return o.Expire(clock.Now()) })`, 3) Return. Se o.Expire retorna ErrNotExpired (pedido novo), ignorar (não é erro). Se Return falha, acumula erro com ID. Propaga errors.Join de erros reais. Sem parar no erro.

**Critério de aceite:**
- [x] `TestExpirarPedidos_09_PE06_ExpiraAntigosDevolvePending` — >30 min → expired (canceled com expiredAt), Return chamado
- [x] `TestExpirarPedidos_09b_NaoExpiraJaCancelado` — já cancelado → ErrInvalidTransition (de o.Expire), sem Return
- [x] `TestExpirarPedidos_10_ReturnFailaAcumulaErro` — Return falha, erro acumulado com ID, demais processados, errors.Join propagado
- [x] `TestExpirarPedidos_MistoNovoEAntigo` — pendidos novos + antigos; antigos expiram+devolvem, novos ignorados, erro vazio (sem erros reais)
- [x] `TestExpirarPedidos_ListPendingFalha` — ListPending falha, nenhum Return chamado
- [x] `TestExpirarPedidos_UpdateFalhaComErroReal` — Update falha com erro real, acumula com ID, errors.Is acha erro original, demais processados
- [x] `make check` passa (100% coverage, variável `pendentes` em vez de `pendidos`)

**Arquivo:** `internal/usecase/ordering/order.go` + `order_test.go`

---

## T006: Teste do problema conhecido (Return falha após Cancel/Expire)

**O quê:** Demonstrar: se Order.Cancel/Expire der OK, Update OK, mas Return falha, a reserva fica presa. Erro deve ser propagado com ID do pedido. T005 já testa erro acumulado em ExpirarPedidos; aqui: teste explícito de CancelarPedido com Return falha.

**Critério de aceite:**
- [x] `TestCancelarPedido_ReturnFails_PropagaErroComID` — Update ok (o.Cancel() ok), Return falha; erro propagado contém ID do pedido
- [x] Comentário no código: "Estoque fica preso; operacional rastreia pelo ID e manualmente compensa"
- [x] `make check` passa (93.3% coverage)

**Arquivo:** `internal/usecase/ordering/order_test.go`

---

## T007: Concorrência — dois cancelamentos simultâneos (ES-05 pattern)

**O quê:** Tarefa extra: verificar que Update serializa: dois CancelarPedido simultâneos no mesmo pedido, um ganha a transição (pending → canceled), outro falha com ErrInvalidTransition. Apenas um Return chamado. Teste com -race.

**Critério de aceite:**
- [ ] `TestCancelarPedido_RaceCondition_DoisCancelamentosSimultaneos` — 50 goroutines (canal start), todos tentam `repo.Update` do mesmo pending; 1 sucesso (o.Cancel ok, Return 1×), 49 falham com ErrInvalidTransition (o.Cancel retorna erro, sem Return)
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

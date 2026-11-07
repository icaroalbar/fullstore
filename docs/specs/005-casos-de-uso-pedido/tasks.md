# Tarefas 005: Casos de Uso — Pedido

## T001: Mudar Order.Expire para retornar ErrInvalidTransition (Spec 002)

**O quê:** Pré-requisito de Spec 005. Order.Expire (em domain/order/order.go) atualmente é idempotente (retorna nil para não pending). Mudar para retornar `ErrInvalidTransition` se o pedido não estiver pending.

**Por quê:** Impede que ExpirarPedidos tente Return duplicado; garante que o uso é seguro em Spec 005.

**Critério de aceite:**
- [ ] `TestOrder_PE06_ExpireBlocksNotPending` — Expire em pedido já canceled retorna ErrInvalidTransition
- [ ] Spec 002 atualizada para refletir novo comportamento
- [ ] `make check` passa

**Arquivo:** `internal/domain/order/order.go` + `order_test.go`

---

## T002: Portas e erros do caso de uso (ES-01, ES-02, ES-03)

**O quê:** Struct `OrderUseCases`, construtor, injeção de 4 portas (ProductCatalog, OrderRepository, InventoryPort, Clock). Erros do caso de uso em var block (ErrProductNotFound, ErrProductUnavailable, ErrOrderNotFound).

**Critério de aceite:**
- [ ] `internal/usecase/order/order.go` criado com construtor
- [ ] Portas declaradas como interfaces
- [ ] Erros em var block
- [ ] Testes falsos (fakes) para cada porta
- [ ] `make check` passa

**Arquivo:** `internal/usecase/order/order.go` + `order_test.go`

---

## T003: CriarPedido validações e reserva (Cenários 01-05, ES-01, ES-02, PR-04, PE-01-03)

**O quê:** Implementar CriarPedido(items) com ordem: 1) catálogo, 2) NewOrder memória, 3) Reserve, 4) Save. Cobrir todos os cenários de erro.

**Critério de aceite:**
- [ ] `TestCriarPedido_ES01_CriarReserva` — 01: produtos ativos, estoque reservado
- [ ] `TestCriarPedido_PR04_ProdutoDesativado` — 02: ErrProductUnavailable, sem efeito
- [ ] `TestCriarPedido_ES02_SemEstoque` — 03: ErrInsufficientStock, tudo-ou-nada
- [ ] `TestCriarPedido_04_RollbackGravacao` — 04: gravar falha, Return chamado
- [ ] `TestCriarPedido_04b_PedidoInvalido` — 04b: ordem inválida (PE-01-03), sem efeito
- [ ] `TestCriarPedido_05_ProdutoInexistente` — 05: ErrProductNotFound, sem efeito
- [ ] `make check` passa

**Arquivo:** `internal/usecase/order/order.go` + `order_test.go`

---

## T004: CancelarPedido validações e devolução (Cenários 06-08b, ES-03, PE-05)

**O quê:** Implementar CancelarPedido(orderID) com validação de transição (Order.Cancel) antes de Return.

**Critério de aceite:**
- [ ] `TestCancelarPedido_ES03_CancelaPending` — 06: pending → canceled, estoque devolvido
- [ ] `TestCancelarPedido_PE05_CancelaPago` — 07: pago → erro ErrInvalidTransition, sem Return
- [ ] `TestCancelarPedido_08_CancelaCancelado` — 08: já canceled → erro, sem Return duplicado
- [ ] `TestCancelarPedido_08b_PedidoInexistente` — 08b: inexistente → ErrOrderNotFound
- [ ] `make check` passa

**Arquivo:** `internal/usecase/order/order.go` + `order_test.go`

---

## T005: ExpirarPedidos com clock e retry (Cenários 09-10, PE-06, ES-03)

**O quê:** Implementar ExpirarPedidos(agora time.Time) que lista pending, chama order.Expire(agora), depois Return. Coleta erros sem parar.

**Critério de aceite:**
- [ ] `TestExpirarPedidos_PE06_ExpiraAntigosDevolvePending` — 09: >30 min → canceled com expiredAt, estoque devolvido
- [ ] `TestExpirarPedidos_09b_NaoExpiraJaCancelado` — 09b: já canceled → ErrInvalidTransition, sem Return
- [ ] `TestExpirarPedidos_10_ColetaErrosRetry` — 10: Return falha em um, demais continuam, erro reportado
- [ ] `make check` passa

**Arquivo:** `internal/usecase/order/order.go` + `order_test.go`

---

## T006: Teste de compensação (Return falha)

**O quê:** Cenário 04 e 10 tocam em "gravar falha" e "Return falha". Ambas testadas em T003 e T005. Aqui: adicionar teste explícito que demonstra o problema conhecido e a decisão (log + reportar).

**Critério de aceite:**
- [ ] `TestReturnFail_DeixaReservaEmAberto` — demonstra que se Return falha após Cancel/Expire bem-sucedido, estoque fica preso
- [ ] Teste valida que erro é registrado (log ou retorno)
- [ ] Comentário explica: "Operacional investigará e compensará"
- [ ] `make check` passa

**Arquivo:** `internal/usecase/order/order_test.go`

---

## T007: Cobertura >= 80% + status implementada

**O quê:** Verificar cobertura, atualizar spec.md Status: implementada.

**Critério de aceite:**
- [ ] `go test -cover ./internal/usecase/order` mostra >= 80%
- [ ] Spec 005: Status mudado para `implementada`
- [ ] `make check` passa

**Arquivo:** `docs/specs/005-casos-de-uso-pedido/spec.md`

---

## Ordem de execução

1. **T001** — Pré-requisito: Order.Expire mudança (Spec 002)
2. **T002** — Setup: portas, erros, construtor
3. **T003** — CriarPedido: ordem, validações, compensação
4. **T004** — CancelarPedido: transição, idempotência segura
5. **T005** — ExpirarPedidos: clock, retry, coleta erros
6. **T006** — Teste do problema conhecido (Return falha)
7. **T007** — Cobertura + spec implementada

Cada tarefa finaliza com `make check` verde e diff explicado.

---

## Regras cobertas

✅ ES-01 (CriarPedido reserva estoque) — T003
✅ ES-02 (Sem estoque suficiente: tudo-ou-nada) — T003
✅ ES-03 (Cancelar/expirar devolve) — T004, T005
✅ PR-04 (Produto desativado) — T003
✅ PE-06 (Expira em 30 min) — T005

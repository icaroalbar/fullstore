# Tarefas 003: Estoque

## T001: Tipos base Stock, Replenish e getters — ES-01 base

**O quê:** Criar struct Stock com productID, available, reserved, mutex. Implementar Replenish (adiciona ao disponível) e getters Available(), Reserved().

**Critério de aceite:**
- [ ] `TestStock_ReplenishPositive` — Replenish(5) em estoque vazio → available = 5
- [ ] `TestStock_GetterAvailable` — Available() retorna disponível correto
- [ ] `TestStock_GetterReserved` — Reserved() retorna reservado correto
- [ ] `TestStock_ReplenishInvalidQuantity` — Replenish(0) ou Replenish(-1) retorna ErrInvalidQuantity
- [ ] `TestStock_MultipleReplenish` — múltiplas reposições se somam

**Arquivo:** `internal/domain/stock/stock.go` + `stock_test.go` + `errors.go`

---

## T002: Reserve com validação — ES-01, ES-02

**O quê:** Implementar Reserve(items []Item) error que valida tudo antes de mudar quantidades. Retorna ErrInsufficientStock se faltar, ErrInvalidQuantity se qty<=0, soma produtos repetidos.

**Critério de aceite:**
- [ ] `TestStock_ES01_ReserveDecrementsAvailable` — Reserve(3) em 10 disponível → available = 7, reserved = 3
- [ ] `TestStock_ES02_InsufficientStock` — Reserve(15) em 10 disponível → ErrInsufficientStock, nada muda
- [ ] `TestStock_ReserveInvalidQuantity` — Reserve com qty <= 0 → ErrInvalidQuantity
- [ ] `TestStock_ReserveMultipleItems` — Reserve([Notebook×3, Mouse×2]) funciona com múltiplos itens
- [ ] `TestStock_ReserveProductRepetition` — Reserve([Notebook×2, Notebook×3]) soma para 5 e valida tudo
- [ ] `TestStock_ES02b_TwoProductsOneFails` — Reserve([A×3, B×15]) com B tendo 10 → nada é reservado

**Arquivo:** `internal/domain/stock/stock.go` + `stock_test.go`

---

## T003: Return e Deduct — ES-03, ES-04

**O quê:** Implementar Return(items []Item) para devolver reserva (cancel/expire) e Deduct(items []Item) para baixar de vez (pagamento).

**Critério de aceite:**
- [ ] `TestStock_ES03a_ReturnAfterCancel` — Reserve 3, depois Return 3 → available = 10, reserved = 0
- [ ] `TestStock_ES03b_ReturnAfterExpire` — Reserve 3, depois Return 3 → mesmo resultado
- [ ] `TestStock_ReturnInvalidOperation` — Return(5) quando reserved = 3 → ErrInvalidOperation
- [ ] `TestStock_ES04_DeductAfterPay` — Reserve 3 → available = 7, reserved = 3; Deduct 3 → available = 7, reserved = 0 (vendidas de vez)
- [ ] `TestStock_DeductInvalidOperation` — Deduct(5) quando reserved = 3 → ErrInvalidOperation

**Arquivo:** `internal/domain/stock/stock.go` + `stock_test.go`

---

## T004: Teste de concorrência — ES-05

**O quê:** Teste com 50 goroutines tentando comprar 1 unidade cada de estoque com 10 unidades. Deve haver exatamente 10 sucesso, 40 falhas com ErrInsufficientStock. Final: available = 0, reserved = 10.

**Critério de aceite:**
- [ ] `TestStock_ES05_RaceCondition` — 50 goroutines disputam 10 unidades; exatamente 10 Reserve sucesso, 40 ErrInsufficientStock
- [ ] `TestStock_ES05_FinalState` — Após race, available = 0, reserved = 10
- [ ] `-race` flag passa sem data races

**Arquivo:** `internal/domain/stock/stock_test.go`

---

## T005: Cobertura >= 80% e status implementada

**O quê:** Verificar cobertura do pacote domain/stock, confirmar >= 80%, atualizar spec.md para Status: implementada, fazer commit.

**Critério de aceite:**
- [ ] `go test -cover ./internal/domain/stock` mostra >= 80%
- [ ] Spec 003: Status mudado para `implementada`
- [ ] `make check` passa (vet + lint + test -race -cover)

**Arquivo:** `docs/specs/003-estoque/spec.md`

---

## Ordem de execução

1. T001 — base com Replenish + getters
2. T002 — Reserve (validação tudo-ou-nada)
3. T003 — Return + Deduct (transições)
4. T004 — concorrência (ES-05)
5. T005 — cobertura + spec implementada

Cada tarefa finaliza com `make check` passando e diff explicado (comparando com OO).

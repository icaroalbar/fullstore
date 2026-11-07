# Tarefas 003: Estoque

## T001: Tipos base Inventory, Replenish e getters

**O quê:** Criar struct Inventory com mapa productID → {disponível, reservado} + mutex. Implementar Replenish(productID, qty), Available(productID), Reserved(productID) com lock nos getters.

**Critério de aceite:**
- [x] `TestInventory_ReplenishPositive` — Replenish("p1", 5) → Available("p1") = 5
- [x] `TestInventory_ReplenishInvalidQuantity` — Replenish("p1", 0) retorna ErrInvalidQuantity
- [x] `TestInventory_MultipleReplenish` — múltiplas reposições se somam
- [x] `TestInventory_GetterProductNotFound` — Available("unknown") retorna ErrProductNotFound
- [x] `TestInventory_GettersUseLock` — getters adquirem lock (verificar -race)

**Arquivo:** `internal/domain/stock/inventory.go` + `inventory_test.go` + `errors.go`

---

## T002: Reserve com validação — ES-01, ES-02

**O quê:** Implementar Reserve(items []Item) que valida tudo antes de mudar quantidades. Retorna ErrInsufficientStock se faltar, ErrInvalidQuantity se qty≤0, soma produtos repetidos.

**Critério de aceite:**
- [x] `TestInventory_ES01_ReserveDecrementsAvailable` — Reserve([{p1, 3}]) em p1=10 → available=7, reserved=3
- [x] `TestInventory_ES02_InsufficientStock` — Reserve([{p1, 15}]) em p1=10 → ErrInsufficientStock, nada muda
- [x] `TestInventory_ReserveInvalidQuantity` — Reserve com qty ≤ 0 → ErrInvalidQuantity
- [x] `TestInventory_ReserveMultipleItems` — Reserve([{p1,3}, {p2,2}]) funciona
- [x] `TestInventory_ReserveProductRepetition` — Reserve([{p1,2}, {p1,3}]) soma para 5 antes de validar
- [x] `TestInventory_ES02b_AtomicFail` — Reserve([{p1,3}, {p2,15}]) com p2=10 → nada é reservado (tudo-ou-nada)

**Arquivo:** `internal/domain/stock/inventory.go` + `inventory_test.go`

---

## T003: Return e Deduct — ES-03, ES-04

**O quê:** Implementar Return(items []Item) e Deduct(items []Item). Ambos validam tudo antes de mudar, rejeitam qty≤0, retornam ErrInvalidOperation se quantidade > reservado.

**Critério de aceite:**
- [x] `TestInventory_ES03a_ReturnAfterCancel` — Reserve([{p1,3}]), depois Return([{p1,3}]) → available=10, reserved=0
- [x] `TestInventory_ES03b_ReturnAfterExpire` — mesmo resultado
- [x] `TestInventory_ReturnInvalidQuantity` — Return com qty ≤ 0 → ErrInvalidQuantity
- [x] `TestInventory_ReturnProductNotFound` — Return([{unknown, 1}]) → ErrProductNotFound
- [x] `TestInventory_ReturnTooMuch` — Return([{p1,5}]) quando reserved=3 → ErrInvalidOperation, nada muda
- [x] `TestInventory_ReturnProductRepetition` — Return([{p1,2}, {p1,3}]) soma para 5 antes de validar
- [x] `TestInventory_ReturnAtomicFail` — Return([{p1,3}, {p2,15}]) com p2 reserved=10 → tudo-ou-nada, nada muda
- [x] `TestInventory_ES04_DeductAfterPay` — Reserve([{p1,3}]) → available=7,reserved=3; Deduct([{p1,3}]) → available=7,reserved=0
- [x] `TestInventory_DeductInvalidQuantity` — Deduct com qty ≤ 0 → ErrInvalidQuantity
- [x] `TestInventory_DeductProductNotFound` — Deduct([{unknown, 1}]) → ErrProductNotFound
- [x] `TestInventory_DeductProductRepetition` — Deduct([{p1,2}, {p1,3}]) soma para 5 antes de validar
- [x] `TestInventory_DeductTooMuch` — Deduct([{p1,5}]) quando reserved=3 → ErrInvalidOperation

**Arquivo:** `internal/domain/stock/inventory.go` + `inventory_test.go`

---

## T004: Teste de concorrência — ES-05

**O quê:** Teste com 50 goroutines tentando comprar 1 unidade cada de estoque com 10 unidades. Deve haver exatamente 10 sucesso, 40 falhas com ErrInsufficientStock. Também testa leitura concorrente (getters com lock).

**Critério de aceite:**
- [x] `TestInventory_ES05_RaceCondition` — 50 goroutines, cada Reserve([{p1,1}]) com 10 disponível; 10 sucesso, 40 ErrInsufficientStock
- [x] `TestInventory_ES05_FinalState` — Após race, available=0, reserved=10
- [x] `TestInventory_ES05_ConcurrentReads` — Enquanto Reserve corre, 20 goroutines chamam Available/Reserved; sem panic
- [x] `go test -race` passa sem data races

**Arquivo:** `internal/domain/stock/stock_test.go`

---

## T005: Cobertura >= 80% e status implementada

**O quê:** Verificar cobertura do pacote domain/stock, confirmar >= 80%, atualizar spec.md para Status: implementada.

**Critério de aceite:**
- [x] `go test -cover ./internal/domain/stock` mostra >= 80% (94.9%)
- [x] Spec 003: Status mudado para `implementada`
- [x] `make check` passa (vet + lint + test -race -cover)

**Arquivo:** `docs/specs/003-estoque/spec.md`

---

## Ordem de execução

1. T001 — base com Replenish + getters
2. T002 — Reserve (validação tudo-ou-nada)
3. T003 — Return + Deduct (transições)
4. T004 — concorrência (ES-05)
5. T005 — cobertura + spec implementada

Cada tarefa finaliza com `make check` passando e diff explicado (comparando com OO).

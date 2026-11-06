# Tasks 002: Pedido (Order)

## T001: Tipos base Order, OrderItem — PE-01, PE-02

**O quê:** Criar structs Order e OrderItem. Implementar `NewOrder(items []OrderItem)` com validações de quantidade de itens e quantidade por item.

**Critério de aceite:**
- [x] `TestOrder_PE01_EmptyOrder` — 0 itens retorna `ErrEmptyOrder`
- [x] `TestOrder_PE01_TooManyItems` — 21 itens retorna `ErrTooManyItems`
- [x] `TestOrder_PE01_ValidItemCount` — 1-20 itens aceito
- [x] `TestOrder_PE02_InvalidQuantityZero` — item com qty 0 retorna `ErrInvalidQuantity`
- [x] `TestOrder_PE02_InvalidQuantityTooHigh` — item com qty 11 retorna `ErrInvalidQuantity`
- [x] `TestOrder_PE02_ValidQuantity` — item com qty 1-10 aceito
- [x] Order nasce com Status `pending`, ID UUID, CreatedAt

**Arquivo:** `internal/domain/order/order.go` + `order_item.go` + `order_test.go`

---

## T002: Cálculo de total — PE-03

**O quê:** Implementar cálculo de total como soma de (preço × quantidade) por item. Total calculado no NewOrder ou método separado CalculateTotal().

**Critério de aceite:**
- [x] `TestOrder_PE03_SingleItemTotal` — 1 item (preço 5000, qty 2) = total 10000
- [x] `TestOrder_PE03_MultipleItemsTotal` — 3 itens (30×1, 50×2, 20×3) = 30+100+60 = 190 em centavos
- [x] `TestOrder_PE03_TotalImmutable` — criar Order com item, depois mudar preço do Product, total do Order não muda
- [x] Total armazenado em Order.Total (int64, centavos)

**Arquivo:** `internal/domain/order/order.go` + `order_test.go`

---

## T003: Transições de status — PE-04, PE-05

**O quê:** Implementar métodos Pay() e Cancel() com validação de transições. Status: pending → paid (Pay) ou pending → canceled (Cancel).

**Critério de aceite:**
- [x] `TestOrder_PE04_PayPending` — Pay() em pending → status muda para `paid`, retorna sucesso
- [x] `TestOrder_PE04_InvalidTransitionPayPaid` — Pay() em paid retorna `ErrInvalidTransition`
- [x] `TestOrder_PE04_InvalidTransitionPayCanceled` — Pay() em canceled retorna `ErrInvalidTransition`
- [x] `TestOrder_PE05_CancelPending` — Cancel() em pending → status muda para `canceled`, sucesso
- [x] `TestOrder_PE05_InvalidTransitionCancelPaid` — Cancel() em paid retorna `ErrInvalidTransition`
- [x] Order.PaidAt é setado quando Pay() sucesso
- [x] type Status (não string): `"pending"`, `"paid"`, `"canceled"`

**Arquivo:** `internal/domain/order/order.go` + `order_test.go`

---

## T004: Expiração — PE-06

**O quê:** Implementar IsExpired() e Expire() com timeout de 30 minutos. IsExpired() retorna true se pending + (now - CreatedAt > 30min). Expire() muda status para canceled.

**Critério de aceite:**
- [x] `TestOrder_PE06_NotExpiredAt29Min` — pending criado 29 min atrás, IsExpired() = false
- [x] `TestOrder_PE06_ExpiredAt31Min` — pending criado 31 min atrás, IsExpired() = true
- [x] `TestOrder_PE06_NotExpiredPaid` — paid (não pending), IsExpired() = false (só pending expira)
- [x] `TestOrder_PE06_ExpireChangesStatus` — Expire() em pending expirado → status = `canceled`
- [x] `TestOrder_PE06_ExpireIdempotent` — Expire() em já canceled não error
- [x] Order.ExpiredAt é setado quando Expire() sucesso
- [x] Testes usam `time.Now()` ou time mockado (inject clock se necessário)

**Arquivo:** `internal/domain/order/order.go` + `order_test.go`

---

## Ordem de execução
1. T001 — testes vermelho, depois implementação (tipos + validações PE-01, PE-02)
2. T002 — testes vermelho, depois implementação (total PE-03)
3. T003 — testes vermelho, depois implementação (status PE-04, PE-05)
4. T004 — testes vermelho, depois implementação (expiração PE-06)

Cada tarefa finaliza com `make test && make lint` passando, e diff explicado (comparando com OO).

# Tarefas 004: Pagamento

## T001: Tipos base OrderPayments + NewOrderPayments (PG-01, PG-02)

**O quê:** Struct OrderPayments com orderID, total, attempts, mutex. Implementar NewOrderPayments(orderID, total) que valida total > 0.

**Critério de aceite:**
- [ ] `TestOrderPayments_PG01_InvalidTotal` — NewOrderPayments(id, 0) retorna ErrInvalidAmount
- [ ] `TestOrderPayments_PG01_ValidTotal` — NewOrderPayments(id, 10000) cria estrutura
- [ ] `TestOrderPayments_GettersInitial` — Approved() == zero-value, ok==false; Attempts() == []
- [ ] `TestOrderPayments_CopiesDefensive` — Attempts() devolve cópia (modificar não afeta interno)

**Arquivo:** `internal/domain/payment/payment.go` + `payment_test.go`

---

## T002: RecordAttempt com validação de valor (PG-01)

**O quê:** Método RecordAttempt(amount, status) que valida: amount == total, amount > 0, resultado válido.

**Critério de aceite:**
- [ ] `TestOrderPayments_PG01_AmountMismatch` — RecordAttempt(9900, declined) → ErrAmountMismatch
- [ ] `TestOrderPayments_PG01_InvalidAmount` — RecordAttempt(0, approved) → ErrInvalidAmount
- [ ] `TestOrderPayments_InvalidResult` — RecordAttempt(10000, "unknown") → ErrInvalidResult, não registra
- [ ] `TestOrderPayments_PG02_FirstApproved` — RecordAttempt(10000, approved) registra
- [ ] `TestOrderPayments_PG02_DeclinedThenApproved` — Decline, depois approve → ambos registram
- [ ] `TestOrderPayments_OnlyDeclinedNoApproved` — Decline×3 → Approved() retorna zero-value, ok==false

**Arquivo:** `internal/domain/payment/payment.go` + `payment_test.go`

---

## T003: Bloquear qualquer tentativa após aprovado (PG-02, ordem)

**O quê:** RecordAttempt retorna ErrAlreadyPaid para qualquer tentativa (aprovada ou recusada) se já há um aprovado. Validação: ErrAlreadyPaid primeiro (em qualquer caso), depois amount == total.

**Critério de aceite:**
- [ ] `TestOrderPayments_PG02_SecondApprovedFails` — Approve, approve → ErrAlreadyPaid
- [ ] `TestOrderPayments_PG02_DeclineAfterApproveFails` — Approve, decline → ErrAlreadyPaid (não OK)
- [ ] `TestOrderPayments_PG02_OrderValidations` — Aprovado + valor errado → ErrAlreadyPaid (não ErrAmountMismatch)

**Arquivo:** `internal/domain/payment/payment.go` + `payment_test.go`

---

## T004: Getters seguros (PG-04) + histórico em ordem

**O quê:** Approved() retorna snapshot do primeiro attempt aprovado. Attempts() retorna slice em ordem.

**Critério de aceite:**
- [ ] `TestOrderPayments_PG04_ApprovedGetter` — Approve, decline → Approved() retorna só o aprovado
- [ ] `TestOrderPayments_PG04_AttemptsOrder` — Decline, decline, approve → Attempts() em ordem
- [ ] `TestOrderPayments_GettersUseLock` — Getters adquirem lock (verificar -race)

**Arquivo:** `internal/domain/payment/payment.go` + `payment_test.go`

---

## T005: Concorrência PG-02 (ES-05-like)

**O quê:** Teste com 50 goroutines tentando approvar 1x cada, sincronizadas por canal; exatamente 1 sucesso, 49 ErrAlreadyPaid.

**Critério de aceite:**
- [ ] `TestOrderPayments_PG02_RaceCondition` — 50 goroutines (canal start), 1 sucesso, 49 falhas
- [ ] `go test -race -count=50` passa sem data races

**Arquivo:** `internal/domain/payment/payment_test.go`

---

## T006: Cobertura >= 80% + status implementada

**O quê:** Verificar cobertura, atualizar spec.md Status: implementada.

**Critério de aceite:**
- [ ] `go test -cover ./internal/domain/payment` mostra >= 80%
- [ ] Spec 004: Status mudado para `implementada`
- [ ] `make check` passa

**Arquivo:** `docs/specs/004-pagamento/spec.md`

---

## Ordem de execução

1. T001 — base com tipos e getters
2. T002 — RecordAttempt validação (PG-01: amount, PG-02: bloquear após 1º aprovado)
3. T003 — ErrAlreadyPaid para 2º aprovado
4. T004 — Getters e histórico (PG-04)
5. T005 — Concorrência com -race
6. T006 — Cobertura + spec implementada

Cada tarefa finaliza com `make check` passando e diff explicado (comparando com OO).

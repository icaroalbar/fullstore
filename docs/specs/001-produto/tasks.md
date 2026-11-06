# Tasks 001: Produto (Catálogo)

## T001: Tipos base Product — PR-01, PR-02

**O quê:** Criar a struct `Product` com campos Name, Price, Disabled. Implementar `NewProduct()` com validações de nome e preço.

**Critério de aceite:**
- ✅ `TestProduct_PR01_InvalidNameEmpty` — nome vazio retorna `ErrInvalidName`
- ✅ `TestProduct_PR01_InvalidNameTooShort` — nome 1-2 chars retorna erro
- ✅ `TestProduct_PR01_InvalidNameTooLong` — nome 101+ chars retorna erro
- ✅ `TestProduct_PR01_ValidNameMinMax` — nome 3-100 chars é aceito
- ✅ `TestProduct_PR02_InvalidPriceZero` — preço 0 retorna `ErrInvalidPrice`
- ✅ `TestProduct_PR02_InvalidPriceNegative` — preço < 0 retorna erro
- ✅ `TestProduct_PR02_ValidPrice` — preço > 0 é aceito

**Arquivo:** `internal/domain/product/product.go` + `product_test.go`

---

## T002: Status e transições — PR-03

**O quê:** Implementar métodos `Activate()` e `Deactivate()`, garantir desativação idempotente, validar que só ativa com preço > 0.

**Critério de aceite:**
- ✅ `TestProduct_PR03_StartsDisabled` — novo produto tem `Disabled = true`
- ✅ `TestProduct_PR03_ActivateWithValidPrice` — ativar com preço válido muda `Disabled` para `false`
- ✅ `TestProduct_PR03_ActivateRequiresValidPrice` — ativar sem preço válido retorna erro
- ✅ `TestProduct_PR03_DeactivateActive` — desativar produto active muda `Disabled` para `true`
- ✅ `TestProduct_PR03_DeactivateDisabledIdempotent` — desativar disabled novamente não erro

**Arquivo:** `internal/domain/product/product.go` + `product_test.go`

---

## T003: Método IsAvailable e cobertura — PR-04 (preparação)

**O quê:** Implementar `IsAvailable()` (retorna `!Disabled`) para prep de PR-04 (filtro no catálogo). Atingir cobertura 80% do domínio.

**Critério de aceite:**
- ✅ `TestProduct_PR04_DisabledNotAvailable` — produto disabled retorna `IsAvailable() == false`
- ✅ `TestProduct_PR04_ActiveAvailable` — produto active retorna `IsAvailable() == true`
- ✅ `TestProduct_PR05_PriceSnapshot1` — criar pedido com produto a R$ 50 snapshots preço
- ✅ `TestProduct_PR05_PriceSnapshot2` — mudar preço depois não afeta pedido (integration, Order aggregate)
- ✅ Cobertura mínima 80% (`go tool cover`)

**Arquivo:** `internal/domain/product/product.go` + `product_test.go`

---

## Ordem de execução
1. T001 — testes vermelho, depois implementação
2. T002 — testes vermelho, depois implementação
3. T003 — testes vermelho, depois implementação

Cada tarefa finaliza com `make test && make lint` passando, ou volta para o Icaro antes de prosseguir.

# Plan 001: Produto (Catálogo)

## Estratégia
Implementar a entidade `Product` no domínio com validações de nome, preço e transições de status (disabled/active). Os testes guiam o desenvolvimento com TDD: cada cenário da spec vira um teste unitário nomeado pelo ID da regra (PR-01, PR-02, etc.).

## Estrutura

### Domínio (`internal/domain/product/`)
- `product.go`: Struct `Product` com campos (ID, Name, Price, Disabled, CreatedAt)
- Erro package-level: `ErrInvalidName`, `ErrInvalidPrice`
- Métodos:
  - `NewProduct(name string, price int64) (*Product, error)` — cria disabled
  - `Activate() error` — ativa se preço > 0
  - `Deactivate()` — desativa (idempotente)
  - `IsAvailable() bool` — retorna `!Disabled`

### Casos de teste (17 cenários → 17 testes)
- `product_test.go` (domínio)
  - `TestProduct_PR01_ValidateName*` (4 testes)
  - `TestProduct_PR02_ValidatePrice*` (2 testes)
  - `TestProduct_PR03_StartsDisabled*` (2 testes)
  - `TestProduct_PR03_Activate*` (2 testes)
  - `TestProduct_PR04_Deactivate*` (2 testes)
  - `TestProduct_PR05_PriceImmutability*` (2 testes em Order, não aqui)

Nota: PR-04 (catálogo) e PR-05 (pedido) serão testados em outros domínios (use cases / agregado Order).

## Tarefas

1. **T001**: Criar tipos base (`Product` struct com validações PR-01, PR-02)
2. **T002**: Implementar transições de status (PR-03, desativação idempotente)
3. **T003**: Testes 1-16 passando (cobertura 80%)

## Dependências
- Nenhuma importação da infra/usecase (depguard valida)

## Pronto para TDD
Parar em cada tarefa, escrever teste vermelho, depois implementação.

# Plan 002: Pedido (Order)

## Estratégia
Implementar o agregado Order no domínio com validações de itens, cálculo de total, transições de status, e expiração. TDD: cada cenário da spec vira um teste nomeado pelo ID da regra (PE-01, PE-02, etc). Order depende de Product (Spec 001), portanto importa `product` package para validar disponibilidade.

## Estrutura

### Domínio (`internal/domain/order/`)
- `order.go`: Struct `Order` com campos (ID, Items, Status, Total, CreatedAt, ExpiredAt, PaidAt)
- `order_item.go`: Struct `OrderItem` como value object (ProductID, Quantity, Price)
- Erros package-level: `ErrEmptyOrder`, `ErrTooManyItems`, `ErrInvalidQuantity`, `ErrInvalidTransition`, `ErrProductUnavailable`
- Métodos:
  - `NewOrder(items []OrderItem) (*Order, error)` — cria com validações PE-01, PE-02
  - `CalculateTotal() int64` — soma PE-03 (ou calcula no NewOrder)
  - `Pay() error` — transição pending → paid (PE-04)
  - `Cancel() error` — transição pending → canceled (PE-05)
  - `IsExpired() bool` — verifica 30min (PE-06)
  - `Expire() error` — muda para canceled se expirado

### Casos de teste (20 cenários → 20 testes)
- `order_test.go` (domínio)
  - `TestOrder_PE01_Empty*` (3 testes: 0 items, 21 items, range 1-20)
  - `TestOrder_PE02_InvalidQuantity*` (3 testes: 0, 11, range 1-10)
  - `TestOrder_PE03_CalculateTotal*` (3 testes: 1 item, 3 items, snapshot imutável)
  - `TestOrder_PE04_Transitions*` (3 testes: pending→paid, invalid transitions)
  - `TestOrder_PE05_CancelPending*` (2 testes: cancel pending ok, cancel paid erro)
  - `TestOrder_PE06_Expiration*` (3 testes: 29min pending, 31min IsExpired, Expire muda status)

Nota: OrderItem valida ProductID existe/disponível? Talvez em usecase, não domínio (depende de repo). Por enquanto, assume Product válido.

## Tarefas

1. **T001**: Criar tipos base (Order, OrderItem) com validações PE-01, PE-02
2. **T002**: Cálculo de total (PE-03, imutabilidade)
3. **T003**: Transições de status (PE-04, PE-05)
4. **T004**: Expiração (PE-06, 30min timeout)

## Dependências
- `internal/domain/product` — importa para validar Product (Spec 001 já existe)
- Relógio: para PE-06, precisa de `time.Now()` injetável (ou fake em testes)

## Pronto para TDD
Parar em cada tarefa, escrever testes vermelho, depois implementação.

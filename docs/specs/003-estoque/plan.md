# Plano 003: Estoque

Spec: ./spec.md

## Onde mora o código

- **internal/domain/stock**: Agregado Stock com Reserve, Return, Deduct, Replenish e getters
- **internal/domain/stock/errors.go**: ErrInsufficientStock, ErrInvalidQuantity, ErrProductNotFound, ErrInvalidOperation
- **internal/usecase/order**: Integração com Order para chamar Stock.Reserve ao criar pedido
- **internal/usecase/payment**: Integração com Payment para chamar Stock.Deduct ao aprovar pagamento
- **internal/infra/stock**: Repositório Stock (em memória na etapa 1; Postgres na etapa 2)

## Tipos e interfaces

```go
// domain/stock/stock.go
type Stock struct {
  productID string
  available int32
  reserved  int32
  mu        sync.Mutex
}

func NewStock(productID string) *Stock
func (s *Stock) Reserve(items []Item) error         // ES-01, ES-02
func (s *Stock) Return(items []Item) error          // ES-03
func (s *Stock) Deduct(items []Item) error          // ES-04
func (s *Stock) Replenish(quantity int32) error     // Reposição
func (s *Stock) Available() int32
func (s *Stock) Reserved() int32

type Item struct {
  ProductID string
  Quantity  int32
}

// Erros
var (
  ErrInsufficientStock = errors.New(...)
  ErrInvalidQuantity   = errors.New(...)
  ErrProductNotFound   = errors.New(...)
  ErrInvalidOperation  = errors.New(...)
)
```

## Decisões

| Decisão | Alternativa descartada | Motivo |
| --- | --- | --- |
| **Mutex em memória para ES-05** | Pessimistic lock no banco (Postgres) | Etapa 1 é em memória; banco entra na etapa 2. Mutex é simples, testável com goroutines, suficiente agora. Banco será necessário quando múltiplas instâncias rodarem. |
| **Stock agregado separado de Product** | Stock dentro de Product | Stock é domínio próprio (reserva, baixa, reposição). Separação permite ler estoque sem carregar Product e vice-versa. |
| **Replenish para estoque inicial** | EstoqueInicial em NewProduct | Replenish é operação genérica: serve tanto para inicial quanto reposição. Mais flexível e explícito. |

### Comparação com OO

Em OO clássico, Stock seria:
```java
class Stock {
  private int productId;
  private int available;
  private int reserved;
  private Object lock = new Object(); // sync{ }

  public synchronized boolean reserve(List<Item> items) { ... }
  public synchronized void returnReserve(List<Item> items) { ... }
}
```

Em Go:
- `sync.Mutex` + `defer mu.Unlock()` = lock explícito no método, não decorador `synchronized`
- Métodos receptores em ponteiro `(s *Stock)` garantem que mutex funcione em instância compartilhada
- Goroutines testam concorrência sem threads; mais leve que Java

## Riscos

- **ES-05 com múltiplas instâncias**: Mutex em memória protege apenas uma instância. Se Stock rodar em múltiplas goroutines/servidores, pessimistic lock no banco é necessário (etapa 2).
- **Deadlock**: Se Reserve chamar Return dentro de si, mutex segura 2x. Cuidado com chamadas aninhadas.
- **Integração com Order**: Quando Reserve é chamado de Order.Create, a transação deve ser atômica (ordem + estoque). Isso é responsabilidade do use case, não do domínio.

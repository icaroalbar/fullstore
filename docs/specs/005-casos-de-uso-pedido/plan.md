# Plano 005: Casos de Uso — Pedido

Spec: ./spec.md

## Onde mora o código

- **internal/domain/order**: Agregado Order (já existe, será estendido com Expire)
- **internal/domain/product**: Agregado Product (já existe)
- **internal/domain/stock**: Agregado Inventory (já existe)
- **internal/usecase/order**: 3 casos de uso (CriarPedido, CancelarPedido, ExpirarPedidos) + portas + erros

## Tipos e interfaces

```go
// internal/usecase/ordering/order.go (alias: order = "internal/domain/order")

// Portas declaradas pelo caso de uso
type ProductCatalog interface {
  Get(productID string) (*product.Product, error) // ErrProductNotFound
}

type OrderRepository interface {
  Save(order *order.Order) error
  Get(orderID string) (*order.Order, error) // ErrOrderNotFound
  ListPending() ([]*order.Order, error)
}

type InventoryPort interface {
  Reserve(items []stock.Item) error    // tudo-ou-nada; stock.ErrInsufficientStock
  Return(items []stock.Item) error     // devolve, não idempotente
}

type Clock interface {
  Now() time.Time
}

// Erros do caso de uso
var (
  ErrProductNotFound    = errors.New("product not found")
  ErrProductUnavailable = errors.New("product unavailable")
  ErrOrderNotFound      = errors.New("order not found")
)
// ErrInsufficientStock vem de stock.ErrInsufficientStock (propaga do agregado)

// Use cases
type OrderUseCases struct {
  catalog  ProductCatalog
  orders   OrderRepository
  stock    InventoryPort
  clock    Clock
}

func NewOrderUseCases(
  catalog ProductCatalog,
  orders OrderRepository,
  stock InventoryPort,
  clock Clock,
) *OrderUseCases

// ItemInput: ProductID + Quantity; preço vem do catálogo
type ItemInput struct {
  ProductID string
  Quantity  int
}

func (uc *OrderUseCases) CriarPedido(items []ItemInput) (*order.Order, error)

func (uc *OrderUseCases) CancelarPedido(orderID string) error

func (uc *OrderUseCases) ExpirarPedidos() error  // usa clock injetável
```

## Decisões

| Decisão | Alternativa descartada | Motivo |
| --- | --- | --- |
| **Portas pequenas declaradas pelo usecase** | Injetar Inventory/Order diretos | Desacoplamento: caso de uso não conhece implementação; permite fakes nos testes |
| **Reserve/Return com lista inteira** | Item por item | Mantém tudo-ou-nada (ES-02); falha atomicamente |
| **Clock injetável** | time.Now() direto | Testável; ExpirarPedidos não fica acoplado a relógio real |
| **Ordem Cancel/Expire: memória → Save → Return** | Salvar antes de devolver | Se Save falha, não chama Return (sem efeito colateral órfão) |
| **Erros acumulados com errors.Join** | Parar no primeiro erro | ExpirarPedidos processa todos pedidos; erros reportados de uma vez |
| **Return falha: propaga erro com ID do pedido** | Log silencioso | Operacional rastreia qual pedido ficou com estoque preso |

## Comparação com OO

**Go vs Java:**
- **Portas (interfaces):** Em Go explícito (receiver, assinatura). Em Java, classe cliente tem dependência clara com interface.
- **Injeção de dependência:** Go: construtor; Java: constructor ou setter.
- **Erros:** Go: sentinel errors em var. Java: Exceptions (checked/unchecked).

## Riscos

- **Return falha após Cancel/Expire bem-sucedido:** Estoque fica reservado preso. Mitigação: teste T006 demonstra; erro propagado com ID do pedido para operacional investigar.
- **Order.Expire mudança em Spec 002:** Pré-requisito; testa antes de ativar Spec 005 (T001).
- **Concorrência em OrderRepository:** Etapa 1: memória (não thread-safe por padrão Go). Tarefa extra: teste com -race de dois CancelarPedido simultâneos — um ganha, outro falha com ErrInvalidTransition (transição bloqueada); um único Return chamado.
- **Repositório em memória devolve cópias:** Get e ListPending retornam slices copiadas, não referências — mudanças no cliente não afetam estado interno.

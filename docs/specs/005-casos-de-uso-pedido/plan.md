# Plano 005: Casos de Uso — Pedido

Spec: ./spec.md

## Onde mora o código

- **internal/domain/order**: Agregado Order (já existe, será estendido com Expire)
- **internal/domain/product**: Agregado Product (já existe)
- **internal/domain/stock**: Agregado Inventory (já existe)
- **internal/usecase/ordering**: 3 casos de uso (CriarPedido, CancelarPedido, ExpirarPedidos) + portas + erros

## Tipos e interfaces

```go
// internal/usecase/ordering/order.go (alias: order = "internal/domain/order")

// Portas declaradas pelo caso de uso
type ProductCatalog interface {
  Get(productID string) (*product.Product, error) // ErrProductNotFound
}

type OrderRepository interface {
  Save(order *order.Order) error  // CriarPedido: novo pedido
  Update(orderID string, fn func(*order.Order) error) error  // CancelarPedido, ExpirarPedidos: leitura-modificação-gravação atômica; ErrOrderNotFound
  ListPending() ([]*order.Order, error)  // ExpirarPedidos: lista cópias
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
  Quantity  int32
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
| **Update(orderID, fn) atômica** | Get + modificar em memória + Save | Cópias sozinhas + corrida: T1 lê, T2 lê, T1 modifica+salva, T2 modifica+salva → perde escrita de T1. Update serializa com lock/FOR UPDATE |
| **Return após Update ok** | Antes de Update | Se Update falha, não devolver (sem efeito colateral órfão) |
| **Erros reais em ExpirarPedidos** | Acumular ErrNotExpired também | ErrNotExpired = pedido novo (não erro); errors.Join só com erros reais |
| **Erros acumulados com errors.Join** | Parar no primeiro erro | ExpirarPedidos processa todos; erros de uma vez |
| **Return falha: propaga com ID** | Log silencioso | Operacional rastreia qual pedido ficou com estoque preso |

## Comparação com OO

**Go vs Java:**
- **Portas (interfaces):** Em Go explícito (receiver, assinatura). Em Java, classe cliente tem dependência clara com interface.
- **Injeção de dependência:** Go: construtor; Java: constructor ou setter.
- **Erros:** Go: sentinel errors em var. Java: Exceptions (checked/unchecked).

## Riscos

- **Return falha após Cancel/Expire bem-sucedido:** Estoque fica reservado preso. Mitigação: teste T006 demonstra; erro propagado com ID do pedido para operacional investigar.
- **Order.Expire mudança em Spec 002:** Pré-requisito; testa antes de ativar Spec 005 (T001).
- **Concorrência em OrderRepository:** Mitigação: Update(orderID, fn) serializa com mutex em memória (etapa 1) ou SELECT FOR UPDATE no banco (etapa 2+). Tarefa T007: teste com -race de dois CancelarPedido simultâneos — um ganha transição, outro falha com ErrInvalidTransition; um único Return chamado.
- **ListPending devolve cópias:** Slices copiadas, não referências — mudanças no cliente não afetam estado interno.

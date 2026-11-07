# Plano 003: Estoque

Spec: ./spec.md

## Onde mora o código

- **internal/domain/stock**: Agregado Inventory (mapa productID → {disponível, reservado}) com Reserve, Return, Deduct, Replenish, getters
- **internal/domain/stock/errors.go**: ErrInsufficientStock, ErrInvalidQuantity, ErrProductNotFound, ErrInvalidOperation

## Tipos e interfaces

```go
// domain/stock/inventory.go + inventory_test.go
type Inventory struct {
  stocks map[string]*productStock // productID -> {disponível, reservado}
  mu     sync.Mutex
}

type productStock struct {
  available int32
  reserved  int32
}

type Item struct {
  ProductID string
  Quantity  int32
}

func NewInventory() *Inventory
func (inv *Inventory) Reserve(items []Item) error           // ES-01, ES-02
func (inv *Inventory) Return(items []Item) error            // ES-03
func (inv *Inventory) Deduct(items []Item) error            // ES-04
func (inv *Inventory) Replenish(productID string, qty int32) error
func (inv *Inventory) Available(productID string) (int32, error)  // usa lock
func (inv *Inventory) Reserved(productID string) (int32, error)   // usa lock

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
| **Um Inventory com mapa e mutex único** | Múltiplos Stock, cada um com seu mutex | Mapa centralizado simplifica busca; mutex único evita deadlock. Escalas bem em etapa 1. Etapa 2+ migra para Postgres com lock pessimista. |
| **Inventory agregado separado de Product** | Estoque dentro de Product | Inventory é domínio próprio (reserva, baixa, reposição). Separação permite ler estoque sem carregar Product e vice-versa. |
| **Item no pacote stock, não order** | Item importado de order | Item é conceito genérico (productID + quantidade). Evita circular dependency e permite usar Item em outros contextos. |
| **Replenish para estoque inicial** | EstoqueInicial em NewInventory | Replenish é operação genérica: serve tanto para inicial quanto reposição. Mais flexível e explícito. |

### Comparação com OO

Em OO clássico, cada Stock seria:
```java
class Stock {
  private int available;
  private int reserved;

  public synchronized boolean reserve(int qty) { ... }  // synchronized = lock implícito
}
```

Em Go (Inventory com mapa):
- `sync.Mutex` + `defer mu.Unlock()` = lock explícito, não decorador
- Receptor ponteiro `(inv *Inventory)` + mutex garante concorrência segura
- Getters também usam lock (leitura é escrita potencial em mapa)
- Goroutines testam concorrência real sem threads; mais leve que Java

## Riscos

- **ES-05 com múltiplas instâncias**: Mutex em memória protege apenas uma instância. Se Inventory rodar em múltiplas goroutines/servidores, pessimistic lock no banco é necessário (etapa 2+).
- **Deadlock**: Chamadas aninhadas (Reserve → Return) tentam locks 2x. Mitigado: Return/Deduct são chamadas de fora, não recursivamente.
- **Integração com Order**: Quando Reserve é chamado de Order.Create, a transação (criar ordem + reservar estoque) deve ser atômica. Responsabilidade do use case.

# Plano 004: Pagamento

Spec: ./spec.md

## Onde mora o código

- **internal/domain/payment**: Agregado OrderPayments (validação, controle de aprovação, histórico)
- **internal/usecase/payment**: Orquestração (chamar gateway, decidir order.Pay())
- **internal/infra/payment**: Adaptador do gateway (falso para testes)

## Tipos e interfaces

```go
// domain/payment/payment.go
type PaymentStatus string // "approved" ou "declined"

type PaymentAttempt struct {
  Amount int64
  Status PaymentStatus
}

type OrderPayments struct {
  orderID  string
  total    int64
  attempts []PaymentAttempt
  mu       sync.Mutex
}

func NewOrderPayments(orderID string, total int64) (*OrderPayments, error)
func (op *OrderPayments) RecordAttempt(amount int64, result PaymentStatus) error
func (op *OrderPayments) Approved() *PaymentAttempt
func (op *OrderPayments) Attempts() []PaymentAttempt

// Erros
var (
  ErrAlreadyPaid = errors.New("...")
  ErrAmountMismatch = errors.New("...")
  ErrInvalidAmount = errors.New("...")
)
```

## Decisões

| Decisão | Alternativa descartada | Motivo |
| --- | --- | --- |
| **Ordem: ErrAlreadyPaid antes de ErrAmountMismatch** | Validar amount antes | Falha rápida: já aprovado é mais crítico; economiza processamento |
| **Mutex em memória (etapa 1)** | Pessimistic lock no banco | Monolito; etapa 2+ migra para banco |
| **OrderPayments sem orderID em RecordAttempt** | Receber cada vez | Imutável; reduz acoplamento |
| **Item do pacote, sem importar order** | Ref. a Order | Desacoplamento: payment só conhece IDs/valores |
| **Erros em var block no agregado** | Arquivo separado | Padrão: order.go, product.go, inventory.go |

### Comparação com OO

Java: `synchronized` implícito no monitor do objeto.
Go: `sync.Mutex` + `defer` explícito. Receiver ponteiro garante concorrência segura sem importar order = maior desacoplamento que Java.

## Riscos

- **PG-02 com 2 aprovações simultâneas:** Duas goroutines escrevem approved. Mitigado por mutex; etapa 3+ (K8s) precisa pessimistic lock. Teste com -race validará.
- **Concorrência entre validação e aplicação:** Mutex serializa; sem window.

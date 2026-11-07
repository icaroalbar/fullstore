# Plano 004: Pagamento

Spec: ./spec.md

## Onde mora o código

- **internal/domain/payment**: Agregado OrderPayments (validação, controle de aprovação, histórico)

## Tipos e interfaces

```go
// domain/payment/payment.go
type Status string

const (
  StatusApproved Status = "approved"
  StatusDeclined Status = "declined"
)

type Attempt struct {
  Amount int64
  Status Status
}

type OrderPayments struct {
  orderID  string
  total    int64
  attempts []Attempt
  mu       sync.Mutex
}

func NewOrderPayments(orderID string, total int64) (*OrderPayments, error)
func (op *OrderPayments) RecordAttempt(amount int64, result Status) error
func (op *OrderPayments) Approved() (Attempt, bool)
func (op *OrderPayments) Attempts() []Attempt

// Erros
var (
  ErrInvalidAmount = errors.New("...")
  ErrAmountMismatch = errors.New("...")
  ErrAlreadyPaid = errors.New("...")
  ErrInvalidResult = errors.New("...")
)
```

## Decisões

| Decisão | Alternativa descartada | Motivo |
| --- | --- | --- |
| **Ordem: ErrAlreadyPaid antes de ErrAmountMismatch** | Validar amount antes | Falha rápida: já aprovado é mais crítico; economiza processamento |
| **Mutex em memória (etapa 1)** | Pessimistic lock no banco | Monolito; etapa 2+ migra para banco com SELECT...FOR UPDATE |
| **OrderPayments sem orderID em RecordAttempt** | Receber cada vez | Imutável; reduz acoplamento |
| **Item do pacote, sem importar order** | Ref. a Order | Desacoplamento: payment só conhece IDs/valores |
| **Erros em var block no agregado** | Arquivo separado | Padrão: order.go, product.go, inventory.go |

### Comparação com OO

Java: `synchronized` implícito no monitor do objeto.
Go: `sync.Mutex` + `defer` explícito. Receiver ponteiro garante concorrência segura.

## Riscos

- **PG-02 com 2 aprovações simultâneas:** Duas goroutines escrevem approved. Mitigado por mutex; etapa 2+ (múltiplas instâncias) precisa pessimistic lock. Teste com -race validará.
- **Concorrência entre validação e aplicação:** Mutex serializa; sem window.

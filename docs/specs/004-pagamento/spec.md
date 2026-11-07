# Spec 004: Pagamento

Status: rascunho
Etapa: 1
Regras cobertas: PG-01, PG-02, PG-04

## Contexto
Um pedido pending precisa ser pago para virar paid. O cliente pode tentar pagar mais de uma vez (cartão recusado, tenta outro). Cada tentativa é registrada com o valor em centavos e um resultado: aprovada ou recusada. O domínio garante que o valor bate com o total do pedido e que o pedido nunca tem dois pagamentos aprovados. Falar com o gateway de pagamento não é papel do domínio: ele só recebe o resultado.

## Comportamento esperado

**Validação de valor (PG-01):**
- Cenário 01a: Dado um pedido com total 10000 centavos, quando registro um pagamento de 10000, então a tentativa é aceita para processamento.
- Cenário 01b: Dado um pedido com total 10000 centavos, quando registro um pagamento de 9900, então retorna `ErrAmountMismatch` e nada é registrado.
- Cenário 01c: Dado um pedido com total 10000 centavos, quando registro um pagamento de 0 ou negativo, então retorna `ErrInvalidAmount`.

**Máximo um pagamento aprovado por pedido (PG-02):**
- Cenário 02a: Dado um pedido sem pagamento aprovado, quando uma tentativa é aprovada, então o pedido passa a ter 1 pagamento aprovado.
- Cenário 02b: Dado um pedido que já tem um pagamento aprovado, quando outra tentativa é registrada (aprovada ou recusada), então retorna `ErrAlreadyPaid`.
- Cenário 02c: Dado um pedido com uma tentativa recusada, quando uma nova tentativa é aprovada, então é aceita (recusa não bloqueia).

**Pagamento recusado mantém pedido pending (PG-04):**
- Cenário 04a: Dado uma tentativa de pagamento, quando o gateway recusa, então a tentativa fica como recusada e o domínio não altera o pedido.
- Cenário 04b: Dado um pedido com tentativas recusadas, quando consulto o histórico, então todas as tentativas aparecem, na ordem em que foram feitas.

## Casos de erro
| Situação | Erro esperado |
| --- | --- |
| Valor do pagamento ≠ total do pedido | `ErrAmountMismatch` |
| Valor do pagamento ≤ 0 | `ErrInvalidAmount` |
| Registrar pagamento quando já há um aprovado | `ErrAlreadyPaid` |

## Fora do escopo
- PG-03 (idempotência): fica para cases de uso/infra
- PG-05 (gateway timeout): fica para casos de uso/infra
- Estorno (admin cancelar pedido pago): etapa 4+
- Gateway real: porta (interface) com falso para testes

## Dúvidas em aberto
- _(nenhuma — Contexto e cenários aprovados por Icaro)_

# Spec 004: Pagamento

Status: rascunho
Etapa: 1
Regras cobertas: PG-01, PG-02, PG-04

## Contexto
Um pedido pending precisa ser pago para virar paid. O cliente pode tentar pagar mais de uma vez (cartão recusado, tenta outro). Cada tentativa é registrada com o valor em centavos e um resultado: aprovada ou recusada. O domínio garante que o valor bate com o total do pedido e que o pedido nunca tem dois pagamentos aprovados. Falar com o gateway de pagamento não é papel do domínio: ele só recebe o resultado.

## Comportamento esperado

**Validação de valor (PG-01):**
- Cenário 01a: Dado um pedido com total 10000 centavos, quando registro uma tentativa de 10000, então a tentativa fica registrada.
- Cenário 01b: Dado um pedido com total 10000 centavos, quando registro uma tentativa de 9900, então retorna `ErrAmountMismatch` e nada é registrado.
- Cenário 01c: Dado um pedido com total 10000 centavos, quando registro uma tentativa de 0 ou negativo, então retorna `ErrInvalidAmount`.

**Máximo um pagamento aprovado por pedido (PG-02):**
- Cenário 02a: Dado um pedido sem pagamento aprovado, quando uma tentativa é aprovada, então o pedido passa a ter 1 pagamento aprovado.
- Cenário 02b: Dado um pedido que já tem um pagamento aprovado, quando outra tentativa é registrada (aprovada ou recusada), então retorna `ErrAlreadyPaid`.
- Cenário 02c: Dado um pedido com uma tentativa recusada, quando uma nova tentativa é aprovada, então é aceita (recusa não bloqueia).

**Pagamento recusado mantém pedido pending (PG-04):**
- Cenário 04a: Dado um pedido, quando registro uma tentativa recusada, então a tentativa fica registrada e o status do pedido não muda (continua pending).
- Cenário 04b: Dado um pedido com tentativas recusadas, quando consulto o histórico, então todas as tentativas aparecem, na ordem em que foram feitas.

## Casos de erro
| Situação | Erro esperado |
| --- | --- |
| Valor do pagamento ≠ total do pedido | `ErrAmountMismatch` |
| Valor do pagamento ≤ 0 | `ErrInvalidAmount` |
| Registrar pagamento quando já há um aprovado | `ErrAlreadyPaid` |

## Decisões de design

| Decisão | Motivo |
| --- | --- |
| **OrderPayments como agregado por pedido** | Centraliza todas as tentativas de pagamento de um pedido (ID do pedido, total em centavos recebido, lista de tentativas em ordem). Garante atomicamente PG-01 (valor == total) e PG-02 (máximo 1 aprovado). Não importa o pacote order. |
| **Uma só operação: registrar tentativa** | Cliente chama uma operação que aceita: pedido ID, valor, resultado (aprovado/recusado). Ela valida e armazena. |

## Fora do escopo
- Validar se o pedido está pending (responsabilidade de casos de uso)
- Quem chama order.Pay() para confirmar a mudança de status (casos de uso)
- PG-03 (idempotência): fica para casos de uso/infra
- PG-05 (gateway timeout): fica para casos de uso/infra
- Estorno (admin cancelar pedido pago): etapa 4+
- Gateway real: porta (interface) com falso para testes

## Dúvidas em aberto
- **PG-02 segura com 2 aprovações simultâneas?** Decisão de concorrência para o plano (como ES-05 do estoque).

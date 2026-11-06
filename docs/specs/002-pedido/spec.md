# Spec 002: Pedido (Order)

Status: aprovada
Etapa: 1
Regras cobertas: PE-01, PE-02, PE-03, PE-04, PE-05, PE-06

## Contexto
O Pedido (Order) é o aggregate root que agrupa itens (OrderItem), cada um referenciando um Produto ativo com quantidade e preço snapshot. O domínio calcula o total, valida transições de status (pending → paid/canceled), e garante que pedidos pending expirem em 30 minutos sem pagamento. Itens podem ser de 1 a 20 distintos, cada um com quantidade 1-10. Order ID é UUID (distribuído). OrderItem é value object (sem ID). Total é int64 centavos. Pedidos expirados ficam com status canceled, consultáveis no histórico.

## Comportamento esperado
**Criação de pedido (PE-01, PE-02, PE-03):**
- Cenário 1: Dado nada, quando crio pedido com 1 item válido (qtd 1-10), então o pedido nasce com status `pending` e ID único.
- Cenário 2: Dado nada, quando crio pedido com 20 itens distintos, então é aceito (max).
- Cenário 3: Dado nada, quando crio pedido vazio (0 itens), então retorna `ErrEmptyOrder`.
- Cenário 4: Dado nada, quando crio pedido com 21 itens, então retorna `ErrTooManyItems`.
- Cenário 5: Dado nada, quando crio pedido com item de quantidade 0, então retorna `ErrInvalidQuantity`.
- Cenário 6: Dado nada, quando crio pedido com item de quantidade 11, então retorna `ErrInvalidQuantity`.
- Cenário 7: Dado nada, quando crio pedido com item de quantidade 1-10, então é aceito.

**Cálculo de total (PE-03):**
- Cenário 8: Dado pedido com 1 item (preço R$ 50, qtd 2), então total = R$ 100 (calculado pelo domínio).
- Cenário 9: Dado pedido com 3 itens (R$ 30×1, R$ 50×2, R$ 20×3), então total = 30+100+60 = R$ 190.
- Cenário 10: Dado pedido criado, quando mudo preço do Produto, então total do pedido não muda (snapshot).

**Status e transições (PE-04, PE-05):**
- Cenário 11: Dado pedido `pending`, quando pago com sucesso, então status muda para `paid`.
- Cenário 12: Dado pedido `pending`, quando cancelado, então status muda para `canceled`.
- Cenário 13: Dado pedido `paid`, quando tento pagar novamente, então retorna `ErrInvalidTransition`.
- Cenário 14: Dado pedido `paid`, quando tento cancelar, então retorna `ErrInvalidTransition`.
- Cenário 15: Dado pedido `canceled`, quando tento pagar, então retorna `ErrInvalidTransition`.

**Cancelamento (PE-05):**
- Cenário 16: Dado pedido `pending`, quando cancelado, então status muda para `canceled` (sem erro).
- Cenário 17: Dado pedido `paid`, quando cancelado, então retorna `ErrInvalidTransition` (não pode cancelar pago).

**Expiração (PE-06):**
- Cenário 18: Dado pedido `pending` criado há 29 min, quando verifico status, então continua `pending`.
- Cenário 19: Dado pedido `pending` criado há 31 min, quando verifico IsExpired(), então retorna true.
- Cenário 20: Dado pedido `pending` criado há 31 min, quando chamo Expire(), então status muda para `canceled`.

## Casos de erro
| Situação | Erro esperado |
| --- | --- |
| Pedido vazio (0 itens) | `ErrEmptyOrder` |
| Pedido com 21+ itens | `ErrTooManyItems` |
| Item com quantidade 0 ou 11+ | `ErrInvalidQuantity` |
| Pagar pedido não pending | `ErrInvalidTransition` |
| Cancelar pedido não pending | `ErrInvalidTransition` |
| Pedido com Produto desativado | `ErrProductUnavailable` |

## Fora do escopo
- Cálculo de frete
- Cupons/descontos
- Múltiplas moedas
- Histórico de mudanças (eventos salvos em outra etapa)

## Dúvidas em aberto
_(Resolvidas — ver Contexto acima)_

## Decisões
- **Order ID:** UUID string (distribuído, rastreável, padrão mercado)
- **OrderItem:** Value object sem ID (ciclo de vida acoplado a Order)
- **Total:** int64 centavos (moeda única BRL, simplicidade, refatoração fácil)
- **Pedidos expirados:** Status canceled, consultáveis (auditoria, rastreabilidade, padrão mercado)

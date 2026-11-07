# Spec 001: Produto (Catálogo)

Status: implementada
Etapa: 1
Regras cobertas: PR-01, PR-02, PR-03, PR-04, PR-05, PR-06

## Contexto
O Produto é a entidade central do catálogo. Nasce desativado (draft) até ter um preço válido, e só pode ser comprado quando ativo. O domínio valida nome, preço e transições de status. Mudanças de preço só ocorrem via método ChangePrice (validado), e pedidos já criados mantêm o preço do momento da compra (imutabilidade).

## Comportamento esperado
**Criação:**
- Cenário 1: Dado nada, quando crio um produto com nome válido e preço > 0, então o produto nasce com status `disabled` (draft).
- Cenário 2: Dado um produto, quando tento criar outro com o mesmo nome, então é permitido (nomes não são únicos).

**Validação de nome (PR-01):**
- Cenário 3: Dado nada, quando crio produto com nome vazio, então recebo `ErrInvalidName`.
- Cenário 4: Dado nada, quando crio produto com nome de 1-2 caracteres, então recebo `ErrInvalidName`.
- Cenário 5: Dado nada, quando crio produto com nome de 101+ caracteres, então recebo `ErrInvalidName`.
- Cenário 6: Dado nada, quando crio produto com nome de 3-100 caracteres, então é aceito.

**Validação de preço (PR-02):**
- Cenário 7: Dado nada, quando crio produto com preço zero ou negativo, então recebo `ErrInvalidPrice`.
- Cenário 8: Dado nada, quando crio produto com preço > 0 em centavos, então é aceito.

**Ativação e desativação (PR-03, RF-01):**
- Cenário 9: Dado um produto disabled, quando ativo com preço válido, então o flag `disabled` muda para `false`.
- Cenário 10: Dado um produto, quando tento ativar sem preço válido, então recebo `ErrInvalidPrice`.
- Cenário 11: Dado um produto active (`disabled = false`), quando desativo, então `disabled` muda para `true`.
- Cenário 12: Dado um produto disabled, quando desativo novamente, então é idempotente (sem erro).

**Indisponibilidade no catálogo (PR-04):**
- Cenário 13: Dado um produto disabled, quando listo o catálogo, então ele não aparece.
- Cenário 14: Dado um produto active, quando listo o catálogo, então aparece com nome e preço.
- Cenário 15: Dado um produto disabled, quando tento criar um pedido com ele, então recebo `ErrProductUnavailable`.

**Imutabilidade de preço no pedido (PR-05):**
- Cenário 16: Dado um pedido com 1 item a R$ 50, quando mudo o preço do produto para R$ 100, então o pedido continua com R$ 50.
- Cenário 17: Dado um pedido criado, então seu preço unitário é snapshot do momento (não segue mudanças futuras do produto).

**Alterar preço com validação (PR-06):**
- Cenário 18: Dado um produto com preço R$ 50, quando chamo ChangePrice(100), então o preço muda para R$ 100 e retorna sucesso.
- Cenário 19: Dado um produto com preço R$ 50, quando chamo ChangePrice(0), então retorna `ErrInvalidPrice` e o preço continua R$ 50.
- Cenário 20: Dado um produto com preço R$ 50, quando chamo ChangePrice(-50), então retorna `ErrInvalidPrice` e o preço continua R$ 50.

## Casos de erro
| Situação | Erro esperado |
| --- | --- |
| Nome vazio ou fora de 3-100 caracteres | `ErrInvalidName` |
| Preço ≤ 0 (NewProduct ou ChangePrice) | `ErrInvalidPrice` |
| Tentar usar produto disabled em pedido | `ErrProductUnavailable` |
| Tentar ativar sem preço válido | `ErrInvalidPrice` |

## Notas de implementação

**PR-04 (Indisponibilidade):** O domínio implementa `IsAvailable()` (verificar se produto está ativo). O erro `ErrProductUnavailable` é responsabilidade da camada de casos de uso (quando tentar reservar estoque de um produto disabled). O domínio não lança esse erro — apenas fornece o getter.

## Fora do escopo
- Multimoeda
- Frontend

## Dúvidas em aberto
- _(nenhuma — PRD e feedback do Icaro resolvem tudo)_

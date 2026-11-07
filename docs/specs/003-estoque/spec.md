# Spec 003: Estoque

Status: rascunho
Etapa: 1
Regras cobertas: ES-01, ES-02, ES-03, ES-04, ES-05

## Contexto

Quando um cliente compra, o estoque deve ser reservado na hora para evitar sobrevenda. Cada produto tem duas quantidades: disponível (pode vender agora) e reservado (em pedidos pending). Quando o pedido é pago, a reserva vira baixa definitiva. Se o pedido expira ou é cancelado, a reserva volta.

## Comportamento esperado

- **ES-01**: Dado produto com 10 disponível, quando crio pedido com 3 unidades, então disponível = 7, reservado = 3
- **ES-02a**: Dado produto com 2 disponível, quando tento pedido com 5 unidades, então retorna ErrInsufficientStock e estoque permanece 2/0
- **ES-02b**: Dado Notebook (10 disp), Mouse (2 disp), quando tento pedido com Notebook×3 + Mouse×5, então nada é reservado
- **ES-03**: Dado pedido com Notebook reservado 3, quando cancelo, então Notebook volta a 10 disponível, 0 reservado
- **ES-04**: Dado pedido com Notebook reservado 3, quando pagamento aprovado, então reservado = 0 (3 vendidas de vez)
- **ES-05**: Dado produto com 10 unidades, quando 50 goroutines tentam comprar simultaneamente, então exatamente 10 aceitas, 40 com ErrInsufficientStock

## Casos de erro

| Situação | Erro esperado |
| --- | --- |
| Reserve: item sem estoque suficiente | ErrInsufficientStock |
| Return: devolver mais do que foi reservado | ErrInvalidOperation |
| Deduct: baixar mais do que foi reservado | ErrInvalidOperation |
| Replenish: quantidade <= 0 | ErrInvalidQuantity |
| Qualquer operação com produto não encontrado | ErrProductNotFound |

## Decisões de design

- **Estrutura**: Cada produto tem dois campos: `disponível` (int32) e `reservado` (int32)
- **Estoque inicial**: NewProduct recebe quantidade inicial como parâmetro
- **Reposição**: Método Replenish(productID, quantidade) adiciona ao disponível
- **Tudo ou nada**: Reserve valida todos os itens antes de mudar qualquer quantidade
  - Quando um pedido tem múltiplos produtos, Reserve deve falhar atomicamente se qualquer item não tiver estoque
  - Responsabilidade: quem chama Reserve (use case Order) deve passar todos os itens de uma vez
- **Erro de estoque**: ErrInsufficientStock quando não há quantidade suficiente

## Operações de Stock

- **Reserve(itens []Item) error** — reserva para um pedido (ES-01, ES-02)
- **Return(itens []Item) error** — devolve uma reserva (ES-03)
- **Deduct(itens []Item) error** — baixa reservado de vez (ES-04)
- **Replenish(productID, quantidade int32) error** — adiciona ao disponível
- **Getters**: Available(productID), Reserved(productID)

## Fora do escopo

- Integração com casos de uso (Reserve chamado ao criar pedido, Return ao cancelar, Deduct ao pagar)
- Frontend de gerenciamento de estoque
- Histórico de movimentações
- Alertas de estoque baixo
- Múltiplos almoxarifados ou localização

## Dúvidas em aberto

- **Técnica de concorrência para ES-05**: Spec define apenas a regra (nunca negativo). Plano e código decidem: mutex em memória (etapa 1) ou pessimistic lock no banco (etapa 2+).

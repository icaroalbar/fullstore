# Spec 003: Estoque

Status: implementada
Etapa: 1
Regras cobertas: ES-01, ES-02, ES-03, ES-04, ES-05

## Contexto

Quando um cliente compra, o estoque deve ser reservado na hora para evitar sobrevenda. Inventory é um agregado centralizado que rastreia estoque de vários produtos via um mapa (productID → {disponível, reservado}). Estoque inicial entra via Replenish. Quando o pedido é pago, a reserva vira baixa definitiva. Se o pedido expira ou é cancelado, a reserva volta.

## Comportamento esperado

- **ES-01**: Dado estoque com 10 disponível, quando reservo pedido com 3 unidades, então disponível = 7, reservado = 3
- **ES-02a**: Dado estoque com 2 disponível, quando tento reservar 5 unidades, então retorna ErrInsufficientStock e estoque permanece 2/0
- **ES-02b**: Dado Notebook (10 disp), Mouse (2 disp), quando tento reservar pedido com Notebook×3 + Mouse×5, então nada é reservado
- **ES-03a**: Dado pedido com Notebook reservado 3, quando cancelo o pedido, então devolvo a reserva e Notebook volta a 10 disponível, 0 reservado
- **ES-03b**: Dado pedido com Notebook reservado 3, quando o pedido expira, então devolvo a reserva e Notebook volta a 10 disponível, 0 reservado
- **ES-04**: Dado pedido com Notebook reservado 3, quando o pagamento é aprovado, então baixo a reserva de vez e Notebook fica 7 disponível, 0 reservado
- **ES-05**: Dado estoque com 10 unidades, quando 50 requisições simultâneas cada uma tenta comprar 1 unidade, então exatamente 10 são aceitas e reservado = 10, disponível = 0

## Casos de erro

| Situação | Erro esperado |
| --- | --- |
| Reserve: item sem estoque suficiente (soma de quantidade) | ErrInsufficientStock |
| Reserve: quantidade <= 0 em qualquer item | ErrInvalidQuantity |
| Reserve: produto não encontrado no inventário | ErrProductNotFound |
| Return: devolver mais do que foi reservado | ErrInvalidOperation |
| Return: quantidade <= 0 em qualquer item | ErrInvalidQuantity |
| Return: produto não encontrado no inventário | ErrProductNotFound |
| Deduct: baixar mais do que foi reservado | ErrInvalidOperation |
| Deduct: quantidade <= 0 em qualquer item | ErrInvalidQuantity |
| Deduct: produto não encontrado no inventário | ErrProductNotFound |
| Replenish: quantidade <= 0 | ErrInvalidQuantity |

## Decisões de design

- **Agregado Stock**: Separado de Product; rastreia um produto por ID e mantém disponível + reservado
- **Estoque inicial**: Entra via Replenish (não em NewProduct)
- **Tudo ou nada**: Reserve valida todos os itens do pedido antes de mudar qualquer quantidade
  - Se um item faltar, nenhum é reservado
  - Responsabilidade: use case Order passa todos os itens de uma vez
- **Produto repetido em qualquer operação**: Reserve, Return, Deduct — se o mesmo produto aparece múltiplas vezes, soma as quantidades antes de validar
- **Validações**: quantidade <= 0 retorna erro antes de qualquer mudança

## Operações de Stock

- **Reserve**: Reserva itens de um pedido (valida tudo antes de mudar)
- **Return**: Devolve uma reserva (quando pedido é cancelado ou expira)
- **Deduct**: Baixa reservado de vez (quando pagamento é aprovado)
- **Replenish**: Adiciona ao disponível (reposição, estoque inicial)
- **Getters**: Consulta de disponível e reservado por produto

## Fora do escopo

- Integração com casos de uso (Reserve chamado ao criar pedido, Return ao cancelar, Deduct ao pagar)
- Frontend de gerenciamento de estoque
- Histórico de movimentações
- Alertas de estoque baixo
- Múltiplos almoxarifados ou localização

## Dúvidas em aberto

(nenhuma)

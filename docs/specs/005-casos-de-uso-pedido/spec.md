# Spec 005: Casos de Uso — Pedido

Status: aprovada
Etapa: 1
Regras cobertas: ES-01, ES-02, ES-03, PR-04, PE-06

## Contexto
Os agregados Produto, Pedido e Estoque já aplicam suas próprias regras, mas nada os une. Estes casos de uso orquestram: criar um pedido só com produtos ativos, reservando o estoque; cancelar ou expirar devolvendo a reserva. Nenhuma regra nova de negócio: eles só garantem a ordem das chamadas e que, se algo falhar no meio, não sobra estoque reservado nem pedido pela metade.

## Comportamento esperado

### CriarPedido
- **01:** Dado produtos ativos com estoque, quando crio o pedido, então ele fica pending, com o preço do momento, e o estoque reservado.
- **02:** Dado um item de produto desativado, quando crio o pedido, então retorna ErrProductUnavailable (erro do caso de uso) e nada é reservado nem gravado (PR-04).
- **03:** Dado um item sem estoque suficiente, quando crio o pedido, então retorna ErrInsufficientStock e nenhum item é reservado (ES-02).
- **04:** Dado que o pedido foi reservado mas gravar falha, quando crio, então a reserva é devolvida e o erro sobe.
- **04b:** Dado um pedido inválido (vazio, >20 itens, quantidade 0/11, item duplicado), quando crio, então retorna ErrEmptyOrder/ErrTooManyItems/ErrInvalidQuantity/ErrDuplicateItem e nada é reservado.
- **05:** Dado um produto que não existe no catálogo, quando crio o pedido, então retorna ErrProductNotFound (erro do caso de uso) e nada muda.

### CancelarPedido
- **06:** Dado um pedido pending, quando cancelo, então fica canceled e a reserva volta (ES-03).
- **07:** Dado um pedido já pago, quando cancelo, então retorna ErrInvalidTransition e a reserva não volta (PE-05).
- **08:** Dado um pedido já cancelado, quando cancelo de novo, então retorna ErrInvalidTransition e a reserva não é devolvida duas vezes.
- **08b:** Dado um pedido que não existe, quando cancelo, então retorna ErrOrderNotFound e nada muda.

### ExpirarPedidos
- **09:** Dado pedidos pending com mais de 30 min, quando rodo a expiração, então eles ficam canceled (com expiredAt marcado) e as reservas voltam; pedidos mais novos não mudam (PE-06).
- **09b:** Dado um pedido já cancelado que chega à expiração, quando rodo, então retorna ErrInvalidTransition (mudança em spec 002) e a reserva não é devolvida.
- **10:** Dado que a devolução de um pedido falha, quando rodo a expiração, então os demais continuam sendo processados e o erro é reportado.

## Casos de erro
| Situação | Erro esperado |
| --- | --- |
| Produto desativado no pedido | ErrProductUnavailable (caso de uso) |
| Item sem estoque suficiente | ErrInsufficientStock |
| Produto não existe no catálogo | ErrProductNotFound (caso de uso) |
| Pedido inválido (vazio, >20, qty 0/11, duplicado) | ErrEmptyOrder, ErrTooManyItems, ErrInvalidQuantity, ErrDuplicateItem |
| Cancelar pedido pago | ErrInvalidTransition |
| Cancelar pedido já cancelado | ErrInvalidTransition |
| Cancelar pedido inexistente | ErrOrderNotFound |
| Expirar pedido já cancelado | ErrInvalidTransition (mudança spec 002) |
| Gravar pedido depois de reservar | Rollback de reserva, erro sobe |

## Fora do escopo
- Pagamento (fica para Spec 006)
- Administrador cancelar pedido pago (etapa 4)
- Frontend
- Expiração automática (aplicação injeta relógio e chama)
- **Pré-requisito (Spec 002):** Order.Expire precisa retornar ErrInvalidTransition para pedido não pending (implementar com teste antes de Spec 005)

## Decisões de design

| Decisão | Motivo |
| --- | --- |
| **Porta InventoryPort com Reserve/Return de lista** | Preserva tudo-ou-nada do estoque: sem reserva parcial (ES-02) |
| **Ordem CriarPedido: catálogo → NewOrder → Reserve → gravar** | Validações baratas antes, efeitos colaterais por último; só gravação compensa (Return) |
| **ExpirarPedidos recebe `agora` (clock)** | Caso de uso stateless; gatilho (job/HTTP) fora de escopo |
| **Order.Expire marca canceled com expiredAt** | Usa o mesmo status canceled mas registra quando expirou (PE-06) |
| **Order.Expire retorna ErrInvalidTransition se não pending** | Impede Return duplicado; Expire fica idempotente com segurança (pré-requisito spec 002) |
| **CancelarPedido valida Order.Cancel() antes de Return** | Evita Return duplicado; Inventory não é idempotente |

## Dúvidas em aberto
- **Falha de Return após Cancel/Expire bem-sucedido:** Se o Order.Cancel() ou Order.Expire() der certo mas o Inventory.Return() falhar, a reserva fica presa no agregado. Como tratar? (Compensação? Log e fila de reprocessamento?) — o plano decide.

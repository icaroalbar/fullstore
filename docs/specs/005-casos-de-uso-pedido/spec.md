# Spec 005: Casos de Uso — Pedido

Status: rascunho
Etapa: 1
Regras cobertas: ES-01, ES-02, ES-03, PR-04, PE-06

## Contexto
Os agregados Produto, Pedido e Estoque já aplicam suas próprias regras, mas nada os une. Estes casos de uso orquestram: criar um pedido só com produtos ativos, reservando o estoque; cancelar ou expirar devolvendo a reserva. Nenhuma regra nova de negócio: eles só garantem a ordem das chamadas e que, se algo falhar no meio, não sobra estoque reservado nem pedido pela metade.

## Comportamento esperado

### CriarPedido
- **01:** Dado produtos ativos com estoque, quando crio o pedido, então ele fica pending, com o preço do momento, e o estoque reservado.
- **02:** Dado um item de produto desativado, quando crio o pedido, então retorna ErrProductUnavailable e nada é reservado nem gravado (PR-04).
- **03:** Dado um item sem estoque suficiente, quando crio o pedido, então retorna ErrInsufficientStock e nenhum item é reservado (ES-02).
- **04:** Dado que o pedido foi reservado mas gravar falha, quando crio, então a reserva é devolvida e o erro sobe.
- **05:** Dado um produto que não existe no catálogo, então retorna erro e nada muda.

### CancelarPedido
- **06:** Dado um pedido pending, quando cancelo, então fica canceled e a reserva volta (ES-03).
- **07:** Dado um pedido já pago, quando cancelo, então retorna erro e a reserva não volta (PE-05).
- **08:** Dado um pedido já cancelado, quando cancelo de novo, então retorna erro e a reserva não é devolvida duas vezes.

### ExpirarPedidos
- **09:** Dado pedidos pending com mais de 30 min, quando rodo a expiração, então eles ficam expirados e as reservas voltam; pedidos mais novos não mudam (PE-06).
- **10:** Dado que a devolução de um pedido falha, quando rodo a expiração, então os demais continuam sendo processados e o erro é reportado.

## Casos de erro
| Situação | Erro esperado |
| --- | --- |
| Produto desativado no pedido | ErrProductUnavailable |
| Item sem estoque suficiente | ErrInsufficientStock |
| Produto não existe | ErrProductNotFound |
| Cancelar pedido pago | ErrInvalidTransition |
| Cancelar pedido já cancelado | ErrInvalidTransition |
| Gravar pedido depois de reservar | Rollback de reserva, erro sobe |

## Fora do escopo
- Pagamento (fica para Spec 006)
- Administrador cancelar pedido pago (etapa 4)
- Frontend
- Expiração automática (aplicação injeta relógio e chama)

## Decisões de design

| Decisão | Motivo |
| --- | --- |
| **Porta InventoryPort com Reserve/Return de lista** | Preserva tudo-ou-nada do estoque: sem reserva parcial (ES-02) |
| **Ordem CriarPedido: catálogo → NewOrder → Reserve → gravar** | Validações baratas antes, efeitos colaterais por último; só gravação compensa (Return) |
| **ExpirarPedidos recebe `agora` (clock)** | Caso de uso stateless; gatilho (job/HTTP) fora de escopo |
| **ExpirarPedidos usa Expire (não Cancel)** | Transição diferente: pending → expired (não canceled) |
| **CancelarPedido valida Order.Cancel() antes de Return** | Evita Return duplicado; Inventory não é idempotente |

## Dúvidas em aberto
- **Falha de Return após Cancel/Expire bem-sucedido:** Se o Order.Cancel() ou Order.Expire() der certo mas o Inventory.Return() falhar, a reserva fica presa no agregado. Como tratar? (Compensação? Log e fila de reprocessamento?) — o plano decide.

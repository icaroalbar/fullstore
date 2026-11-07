# Spec F01: Catálogo de produtos (front end)

Status: aprovada
Etapa: 1
Regras cobertas: RF-02, FE-01, FE-02, FE-03, FE-05

## Contexto
O cliente precisa ver o que a loja vende antes de fazer um pedido (RF-02). Esta é a primeira tela do front end da FullStore, em React, e também cria o layout que as próximas telas vão reaproveitar. A camada HTTP do Go ainda não existe, então a tela conversa com uma API simulada que segue o mesmo contrato que o backend vai implementar depois. Quando o HTTP do Go ficar pronto, basta apontar o front para ele, sem mudar a tela.

## Comportamento esperado
- **FE-01a:** Dado que a API devolve 3 produtos ativos, quando abro o catálogo, então vejo os 3, cada um com nome e preço.
- **FE-01b:** Dado um produto com preço 123456 centavos, quando ele aparece no catálogo, então o preço exibido é "R$ 1.234,56".
- **FE-01c:** Dado um produto com preço 990 centavos, quando ele aparece, então o preço exibido é "R$ 9,90".
- **FE-01d:** Dado um produto com imagem, quando ele aparece, então mostro a imagem; sem imagem, mostro um marcador neutro no mesmo tamanho.
- **FE-02a:** Dado um produto com 0 unidades disponíveis, quando ele aparece, então mostra o selo "Esgotado".
- **FE-02b:** Dado um produto com 5 unidades disponíveis, quando ele aparece, então mostra "Em estoque" e não mostra "Esgotado".
- **FE-03a:** Dado que a API ainda não respondeu, quando abro o catálogo, então vejo um indicador de carregamento.
- **FE-03b:** Dado que a API devolve uma lista vazia, quando abro o catálogo, então vejo "Nenhum produto disponível no momento".
- **FE-03c:** Dado que a API falha, quando abro o catálogo, então vejo "Não foi possível carregar os produtos" e um botão "Tentar de novo"; ao clicar, a tela busca de novo.
- **FE-05a:** Dado qualquer página da loja, quando ela abre, então vejo o cabeçalho com o logo FullStore e o rodapé.
- **FE-05b:** Dado que rolo a página, quando desço, então o cabeçalho continua visível no topo.
- **FE-05c:** Dado que passo o mouse sobre um produto, quando o cursor entra no card, então ele se destaca com animação; com `prefers-reduced-motion` ligado, ele se destaca sem animação.

## Casos de erro
| Situação | O que o cliente vê |
| --- | --- |
| API fora do ar ou erro 5xx | Mensagem de erro + "Tentar de novo" (FE-03c) |
| Lista vazia | Mensagem de catálogo vazio (FE-03b) |

## Decisões
- O front confia no backend para as regras: quem filtra produto desativado é a API (PR-04), e o front só exibe o que recebe.
- Produto esgotado aparece no catálogo com selo, em vez de sumir.
- O preço chega em centavos (int64) e o front só formata; nunca faz conta com dinheiro.
- A tela funciona em celular (uma coluna) e em desktop (grade).
- Estrutura de loja inspirada na Amazon e acabamento escuro com destaque violeta (ver "Identidade visual" em `docs/front-end.md`).
- Imagem do produto: `image_url` é opcional no contrato. Agora só a API simulada preenche; o Produto do Go ganha o campo numa spec do backend depois (regra PR-07).
- Visual com componentes do shadcn/ui: cada produto é um Card, o selo é um Badge, o carregamento usa Skeleton e o erro usa Alert com Button. Detalhes de stack em `docs/front-end.md`.

## Fora do escopo
- Carrinho, botão de comprar, página de detalhe do produto
- Busca, faixa de categorias e carrinho (o cabeçalho reserva o espaço, mas eles entram nas specs F03 e F04)
- Filtros e paginação
- Banner e prateleiras da home
- Login e telas de administrador
- Ligação com a API Go real (entra quando a camada HTTP existir)

## Dúvidas em aberto
(nenhuma)

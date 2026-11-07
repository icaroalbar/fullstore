# Front end da FullStore

Guia do front. O CLAUDE.md aponta para cá quando a tarefa mexe em `web/`.

## Stack
| Peça | Escolha | Por quê |
| --- | --- | --- |
| Base | React + Vite + TypeScript | SPA simples, sem servidor próprio; o backend é a API Go |
| Estilo | Tailwind CSS v4 | Classes utilitárias; é o que o shadcn usa |
| Componentes | shadcn/ui (base `radix`) | O código dos componentes é copiado para `web/src/components/ui/` e passa a ser nosso: dá pra ler e mudar |
| Ícones | lucide-react | Padrão do shadcn |
| Testes | Vitest + Testing Library | Mesma ideia do `go test`: teste nomeado com o ID da regra |
| API simulada | MSW (Mock Service Worker) | Intercepta o `fetch` e responde conforme `docs/api/openapi.yaml` |
| Pacotes | npm | Já vem com o Node; nada extra para aprender |

Comparando com OO: um componente React é parecido com uma classe de view. As props são os parâmetros do construtor, o estado é o atributo privado, e o retorno é o `render()`. O shadcn é uma biblioteca de componentes que você não importa de um pacote: o código-fonte vem para dentro do projeto.

## O que precisa no desktop
1. **Node LTS sem sudo**, via fnm em `~/.local`:
   `curl -fsSL https://fnm.vercel.app/install | bash -s -- --install-dir ~/.local/fnm --skip-shell`, depois adicionar o fnm ao `~/.zshrc` e rodar `fnm install --lts`.
2. **Disco:** o `web/node_modules` ocupa uns 300 MB. Com ~1,9 GB livres cabe, mas vale aumentar o disco antes de crescer.
3. Extensões do VS Code: Tailwind CSS IntelliSense e ESLint.

## Como o projeto nasce (tarefa 1 do plano da F01)
Na raiz do repo:
```bash
npx shadcn@latest init -t vite -b radix -n web --no-monorepo
cd web
npx shadcn@latest add card badge button skeleton alert
```
- `-t vite -n web` cria o projeto Vite + React + TS + Tailwind v4 na pasta `web/`.
- `-b radix` escolhe o Radix como base dos componentes (o mais usado e documentado).
- `--no-monorepo` evita que o shadcn monte um workspace com vários pacotes; o "monorepo" aqui é só o Go e o `web/` no mesmo repo.
- Antes de qualquer `add`, rodar com `--dry-run` para ver o que muda.

## Identidade visual
Estrutura inspirada na Amazon, acabamento inspirado no Felipe Rocha (Full Stack Club), sem copiar layout nem marca de nenhum dos dois.

**Estrutura (Amazon):**
- Cabeçalho fixo com logo, barra de busca larga no centro e carrinho com contador à direita.
- Faixa de categorias logo abaixo do cabeçalho.
- Banner de destaque no topo da home e prateleiras horizontais ("Mais vendidos", "Ofertas").
- Grade de produtos densa: imagem grande, nome em até 2 linhas, preço em destaque, selo de estoque.
- Rodapé com links agrupados.

**Acabamento (Felipe Rocha):**
- Tema escuro como padrão, fundo quase preto e cards um tom acima.
- Cor de destaque violeta para botões, links e preço; verde para "Em estoque", vermelho para "Esgotado" e erros.
- Fonte Inter (ou Geist), títulos em negrito, bordas finas, cantos arredondados (`--radius` ~0.75rem).

**Vida (movimento):**
- Card do produto sobe levemente e ganha borda violeta no hover; a imagem dá um zoom suave.
- Produtos entram com fade + deslize, em cascata.
- Skeleton com brilho enquanto carrega, toast ao adicionar ao carrinho, contador do carrinho "pula".
- Começa com as animações do `tw-animate-css` (já vem com o shadcn); a biblioteca `motion` entra só quando precisar de algo mais elaborado.
- Respeitar `prefers-reduced-motion`: quem desliga animações no sistema vê tudo parado.

Tudo de cor e raio vive nas variáveis CSS do shadcn (`--background`, `--primary`, `--radius` em `src/index.css`), então trocar o tema não mexe nos componentes. Opcional: desenhar a home no pen.dev antes de codar.

## Roteiro de telas
| Spec | Tela |
| --- | --- |
| F01 | Layout da loja (cabeçalho, rodapé) + catálogo em grade |
| F02 | Página do produto |
| F03 | Busca e categorias |
| F04 | Carrinho |
| F05 | Checkout e pagamento (depende do HTTP do Go) |
| F06 | Meus pedidos |

## Estrutura
```
web/
├── components.json          # config do shadcn
├── src/
│   ├── components/ui/       # gerado pelo shadcn (não editar sem motivo)
│   ├── features/catalog/    # telas e componentes da FullStore, por funcionalidade
│   ├── lib/
│   │   ├── api.ts           # cliente HTTP (fetch) para a API
│   │   ├── money.ts         # formatar centavos em BRL
│   │   └── utils.ts         # cn() do shadcn
│   ├── mocks/               # handlers do MSW, seguindo o openapi.yaml
│   └── main.tsx
└── package.json
docs/api/openapi.yaml        # contrato entre front e back
```
A pasta `features/` faz o papel das camadas do Go: a tela não chama `fetch` direto, ela usa `lib/api.ts`. Quando a API Go existir, só a URL muda.

## Comandos
| Comando | Faz |
| --- | --- |
| `make web-dev` | `npm --prefix web run dev` (com MSW ligado) |
| `make web-test` | `npm --prefix web test -- --run` |
| `make web-check` | lint + typecheck + testes do front |
| `make check` | continua só Go; o front tem seu próprio check |

## CI
Novo job `web` no `ci.yml`, que roda só quando muda algo em `web/` ou `docs/api/`: `setup-node` com cache do npm, `npm ci`, `npm run lint`, `npm run typecheck`, `npm test -- --run`.

## Convenções
- Testes: `it("FE-01b formata 123456 centavos como R$ 1.234,56")`.
- Dinheiro sempre chega em centavos; só `lib/money.ts` formata (`Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" })`).
- O front não decide regra de negócio (FE-04).
- Commits: `feat(web): FE-02 selo de esgotado`.

## Mudanças no harness
- CLAUDE.md: seção "Front end" apontando para este guia, e os comandos `make web-*`.
- Hook `go-check.sh`: ignorar arquivos fora de `.go` (já deveria). Um hook equivalente para `web/` (prettier/eslint) fica para depois.

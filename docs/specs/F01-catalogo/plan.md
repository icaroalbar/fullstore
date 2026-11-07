# Plan F01: Catálogo de produtos

## Onde mora o código
```
web/
├── src/
│   ├── components/ui/          # shadcn/ui (gerado, não editar)
│   ├── features/catalog/
│   │   ├── Catalog.tsx         # tela principal (FE-01, FE-02, FE-03)
│   │   ├── ProductCard.tsx     # componente de produto
│   │   └── __tests__/
│   │       └── Catalog.test.ts
│   ├── layouts/
│   │   ├── MainLayout.tsx      # cabeçalho + rodapé (FE-05)
│   │   └── __tests__/
│   │       └── MainLayout.test.ts
│   ├── lib/
│   │   ├── api.ts             # cliente HTTP (GET /products)
│   │   ├── money.ts           # formatar centavos em BRL
│   │   └── __tests__/
│   │       ├── api.test.ts
│   │       └── money.test.ts
│   ├── mocks/
│   │   ├── handlers.ts        # handlers do MSW (GET /products)
│   │   └── server.ts          # setup do MSW
│   └── main.tsx
└── package.json
```

## Tipos TypeScript

```typescript
// lib/api.ts
interface CatalogProduct {
  id: string;
  name: string;
  price_cents: number;
  available: number;
  image_url?: string;
}

interface FetchProductsResponse {
  products: CatalogProduct[];
}

// Função (só assinatura)
export async function fetchProducts(): Promise<CatalogProduct[]>;

// features/catalog/Catalog.tsx
interface CatalogProps {
  // nenhuma prop obrigatória; tudo vem de fetch interno
}

// features/catalog/ProductCard.tsx
interface ProductCardProps {
  product: CatalogProduct;
}

// layouts/MainLayout.tsx
interface MainLayoutProps {
  children: React.ReactNode;
}

// lib/money.ts
export function formatPrice(cents: number): string;
```

## Handlers MSW (só assinaturas)

```typescript
// mocks/handlers.ts
import { http, HttpResponse } from 'msw';

export const handlers = [
  http.get('/api/products', async ({ request, params, cookies }) => {
    // GET /api/products → 200 com lista de produtos
  }),
];

// Resposta simulada:
// {
//   "products": [
//     {
//       "id": "123e4567-e89b-12d3-a456-426614174000",
//       "name": "Laptop",
//       "price_cents": 299900,
//       "available": 5,
//       "image_url": "https://example.com/laptop.jpg"
//     },
//     ...
//   ]
// }
```

## Decisões

### 1. Camadas do front end (FE-04)
- **features/catalog/Catalog.tsx** (tela) → busca via `lib/api.ts` → `/api/products` (nesta etapa, MSW; depois, Go)
- Tela **nunca** chama `fetch` direto. Toda lógica de HTTP fica em `lib/api.ts`.
- Regras de negócio (PR-04, FE-04) vêm da API: o front confia.

**Comparando com OO:** 
- `Catalog` é como uma classe de view que depende injetar uma interface de "repositório" (aqui, `api.ts`).
- Em Go, teríamos uma interface `ProductRepository` que `Catalog` (handler HTTP) recebe; em React, é a função `fetchProducts` que fica desacoplada.
- MSW simula a camada de dados para não precisar do backend real.

### 2. Formatação de dinheiro (FE-01b, FE-01c)
- `lib/money.ts` usa `Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" })`.
- Testes: `it("FE-01b formata 123456 centavos como R$ 1.234,56")`.
- Never faça operações com centavos no front; apenas formatação.

### 3. Layout compartilhado (FE-05)
- `MainLayout` (Header + Rodapé fixo) envolve todas as páginas.
- Header fixo com `position: sticky` (Tailwind `sticky top-0`).
- Componentes do shadcn: `Card` para produto, `Badge` para "Esgotado"/"Em estoque", `Button` para "Tentar de novo", `Alert` para erro, `Skeleton` para carregamento.

### 4. Tratamento de erro e carregamento (FE-03)
- Estados: `loading`, `error`, `empty`, `success`.
- `try/catch` em `fetchProducts`; error vira "Não foi possível carregar os produtos".
- "Tentar de novo" refaz a busca (refetch).

### 5. Imagen e placeholder (FE-01d)
- Se `image_url` existe e é válida (formato uri), mostra `<img>`.
- Caso contrário, mostra um `<div>` cinza com ícone genérico (lucide-react).

### 6. Produto esgotado (FE-02)
- `available === 0` → `<Badge variant="destructive">Esgotado</Badge>`.
- `available > 0` → `<Badge variant="secondary">Em estoque</Badge>` (ou nada).

### 7. Animações (FE-05c)
- Hover: Card sobe (Tailwind `group-hover:translate-y-[-4px]`) e borda ganha cor (violeta).
- Respeita `prefers-reduced-motion`: wrapper verifica `window.matchMedia("(prefers-reduced-motion: reduce)").matches` e desativa `transitions`.
- Produtos entram com fade (opacity 0 → 1) + slide (usar `framer-motion` se necessário, mas começa com `tw-animate-fade-in`).

### 8. Responsividade
- Mobile: 1 coluna (grid-cols-1).
- Tablet/Desktop: 2–4 colunas (grid-cols-2 md:grid-cols-3 lg:grid-cols-4).

## Alternativas descartadas

1. **Carregar dados na App.tsx root vs. em Catalog.tsx**
   - **Escolha:** Em `Catalog.tsx` (hook `useEffect`).
   - **Por quê:** Specs posteriores (F02 produto, F03 busca) terão dados diferentes. Descentralizar deixa cada tela independente.

2. **Usar Redux/Zustand vs. React hooks**
   - **Escolha:** React hooks (`useState`, `useEffect`).
   - **Por quê:** Só tem um ponto de fetch nesta etapa. Redux entra quando houver complexidade (carrinho compartilhado, etc).

3. **API base URL hardcoded vs. env**
   - **Escolha:** Env (`VITE_API_URL`), padrão `/api`.
   - **Por quê:** Facilita trocar de MSW → backend real sem recompilar.

## Riscos

1. **MSW não interceptar requests corretamente** → Testes falham com "fetch error".
   - **Mitigação:** Rodar `make web-test` com debug, conferir se handlers estão registrados.

2. **Tipos TypeScript não baterem com o OpenAPI** → Frontend espera um campo que não vem.
   - **Mitigação:** `docs/api/openapi.yaml` é a fonte da verdade; tipos gerados de lá (ferramentas como `openapi-typescript` depois, agora manual).

3. **Imagem do backend vier em base64 ou path relativo, não URL** → Placeholder aparece onde não deveria.
   - **Mitigação:** Validar `image_url` como `http://` ou `https://`; ignorar formato inválido. Regra PR-07 fica no backend depois.

4. **Performance com muitos produtos** → Skeleton demora para sair.
   - **Mitigação:** MSW simula resposta rápida; servidor real pode ser lento. Nesta etapa, é OK. Paginação/virtualização entram em specs posteriores.

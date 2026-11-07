# Tasks F01: Catálogo de produtos

## T001: Setup do projeto Vite + shadcn + make targets (30 min)
**Regras:** nenhuma regra diretamente, mas infraestrutura obrigatória para as tarefas seguintes.

Na raiz do repo:
```bash
npx shadcn@latest init -t vite -b radix -n web --no-monorepo
```

Depois:
```bash
cd web
npx shadcn@latest add card badge button skeleton alert
```

Teste antes de cada `add` com `--dry-run`.

Criar em `Makefile` (raiz):
```makefile
web-dev:
	npm --prefix web run dev

web-test:
	npm --prefix web test -- --run

web-check: web-lint web-typecheck web-test

web-lint:
	npm --prefix web run lint

web-typecheck:
	npm --prefix web run typecheck
```

Criar job `web` em `.github/workflows/ci.yml` que roda só quando muda algo em `web/` ou `docs/api/`:
```yaml
web:
  runs-on: ubuntu-latest
  if: |
    github.event_name == 'pull_request' &&
    (contains(github.event.pull_request.paths, 'web/**') || 
     contains(github.event.pull_request.paths, 'docs/api/**'))
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-node@v4
      with:
        node-version: lts/*
        cache: npm
        cache-dependency-path: web/package-lock.json
    - run: npm ci --prefix web
    - run: npm run lint --prefix web
    - run: npm run typecheck --prefix web
    - run: npm run test -- --run --prefix web
```

Instalar MSW:
```bash
npm install -D msw
npx msw init web/public
```

**Teste que prova:** `make web-dev` abre sem erro; `make web-test` roda sem testes ainda (0 passed); `make web-check` passa (lint + typecheck + test).

---

## T002: lib/money.ts - formatar centavos em BRL (30 min)
**Regras:** FE-01b, FE-01c

Criar `web/src/lib/money.ts`:

```typescript
export function formatPrice(cents: number): string {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(cents / 100);
}
```

Criar `web/src/lib/__tests__/money.test.ts`:

```typescript
import { describe, it, expect } from "vitest";
import { formatPrice } from "../money";

describe("money", () => {
  it("FE-01b formata 123456 centavos como R$ 1.234,56", () => {
    expect(formatPrice(123456)).toBe("R$ 1.234,56");
  });

  it("FE-01c formata 990 centavos como R$ 9,90", () => {
    expect(formatPrice(990)).toBe("R$ 9,90");
  });

  it("FE-01 formata 100 centavos como R$ 1,00", () => {
    expect(formatPrice(100)).toBe("R$ 1,00");
  });
});
```

**Teste que prova:** `npm run test -- --run` passa com 3 testes verdes.

---

## T003: lib/api.ts + MSW handlers - GET /api/products (45 min)
**Regras:** FE-01a (lista de produtos), RF-02 (backend simula)

Criar `web/src/lib/api.ts`:

```typescript
export interface CatalogProduct {
  id: string;
  name: string;
  price_cents: number;
  available: number;
  image_url?: string;
}

export async function fetchProducts(): Promise<CatalogProduct[]> {
  const response = await fetch(`${import.meta.env.VITE_API_URL}/products`);
  if (!response.ok) {
    throw new Error("Failed to fetch products");
  }
  const data = await response.json();
  return data.products;
}
```

Criar `web/.env.local`:
```
VITE_API_URL=http://localhost:3000/api
```

Criar `web/src/mocks/handlers.ts`:

```typescript
import { http, HttpResponse } from "msw";
import type { CatalogProduct } from "../lib/api";

const mockProducts: CatalogProduct[] = [
  {
    id: "123e4567-e89b-12d3-a456-426614174000",
    name: "Laptop",
    price_cents: 299900,
    available: 5,
    image_url: "https://via.placeholder.com/300x300?text=Laptop",
  },
  {
    id: "223e4567-e89b-12d3-a456-426614174001",
    name: "Mouse",
    price_cents: 4990,
    available: 0,
  },
  {
    id: "323e4567-e89b-12d3-a456-426614174002",
    name: "Teclado",
    price_cents: 12990,
    available: 10,
    image_url: "https://via.placeholder.com/300x300?text=Teclado",
  },
];

export const handlers = [
  http.get("/api/products", () => {
    return HttpResponse.json({ products: mockProducts });
  }),
];
```

Criar `web/src/mocks/server.ts`:

```typescript
import { setupServer } from "msw/node";
import { handlers } from "./handlers";

export const server = setupServer(...handlers);
```

Atualizar `web/vitest.config.ts` (ou criar se não existir):

```typescript
import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: ["./src/mocks/setup.ts"],
  },
});
```

Criar `web/src/mocks/setup.ts`:

```typescript
import { beforeAll, afterEach, afterAll } from "vitest";
import { server } from "./server";

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
```

Criar `web/src/lib/__tests__/api.test.ts`:

```typescript
import { describe, it, expect } from "vitest";
import { fetchProducts } from "../api";

describe("api", () => {
  it("FE-01a busca 3 produtos da API", async () => {
    const products = await fetchProducts();
    expect(products).toHaveLength(3);
    expect(products[0].name).toBe("Laptop");
    expect(products[0].price_cents).toBe(299900);
  });
});
```

**Teste que prova:** `npm run test -- --run` passa com todos os testes (incluindo money.test.ts).

---

## T004: MainLayout - cabeçalho fixo + rodapé (30 min)
**Regras:** FE-05a, FE-05b

Criar `web/src/layouts/MainLayout.tsx`:

```typescript
import React from "react";

export interface MainLayoutProps {
  children: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-background border-b">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-2xl font-bold">FullStore</h1>
          <div className="flex items-center gap-4">
            {/* Placeholder para barra de busca (F03) */}
            <button className="p-2">🛒</button>
          </div>
        </div>
      </header>

      {/* Main */}
      <main className="flex-1">{children}</main>

      {/* Footer */}
      <footer className="bg-background border-t">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <p className="text-sm text-muted-foreground">
            © 2026 FullStore. Todos os direitos reservados.
          </p>
        </div>
      </footer>
    </div>
  );
}
```

Criar `web/src/layouts/__tests__/MainLayout.test.tsx`:

```typescript
import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MainLayout } from "../MainLayout";

describe("MainLayout", () => {
  it("FE-05a exibe o cabeçalho com logo FullStore", () => {
    render(
      <MainLayout>
        <div>Conteúdo</div>
      </MainLayout>
    );
    expect(screen.getByText("FullStore")).toBeInTheDocument();
  });

  it("FE-05b cabeçalho fica fixo ao rolar (sticky)", () => {
    const { container } = render(
      <MainLayout>
        <div>Conteúdo</div>
      </MainLayout>
    );
    const header = container.querySelector("header");
    expect(header).toHaveClass("sticky");
  });
});
```

**Teste que prova:** `npm run test -- --run` passa com os testes do MainLayout.

---

## T005: ProductCard - nome, preço, imagem, esgotado (30 min)
**Regras:** FE-01, FE-02, FE-01d

Criar `web/src/features/catalog/ProductCard.tsx`:

```typescript
import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Package } from "lucide-react";
import type { CatalogProduct } from "@/lib/api";
import { formatPrice } from "@/lib/money";

interface ProductCardProps {
  product: CatalogProduct;
}

export function ProductCard({ product }: ProductCardProps) {
  const isOutOfStock = product.available === 0;

  return (
    <Card className="overflow-hidden hover:border-primary transition-all duration-200">
      <CardContent className="p-0">
        {/* Imagem */}
        <div className="aspect-square bg-muted flex items-center justify-center overflow-hidden">
          {product.image_url ? (
            <img
              src={product.image_url}
              alt={product.name}
              className="w-full h-full object-cover hover:scale-105 transition-transform duration-200"
            />
          ) : (
            <Package className="w-12 h-12 text-muted-foreground" />
          )}
        </div>

        {/* Info */}
        <div className="p-4">
          <h3 className="font-semibold line-clamp-2">{product.name}</h3>
          <p className="text-lg font-bold text-primary mt-2">
            {formatPrice(product.price_cents)}
          </p>

          {/* Status */}
          <div className="mt-2">
            {isOutOfStock ? (
              <Badge variant="destructive">Esgotado</Badge>
            ) : (
              <Badge variant="secondary">Em estoque</Badge>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
```

Criar `web/src/features/catalog/__tests__/ProductCard.test.tsx`:

```typescript
import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { ProductCard } from "../ProductCard";
import type { CatalogProduct } from "@/lib/api";

describe("ProductCard", () => {
  const mockProduct: CatalogProduct = {
    id: "123",
    name: "Laptop",
    price_cents: 123456,
    available: 5,
    image_url: "https://example.com/laptop.jpg",
  };

  it("FE-01 exibe nome e preço do produto", () => {
    render(<ProductCard product={mockProduct} />);
    expect(screen.getByText("Laptop")).toBeInTheDocument();
    expect(screen.getByText("R$ 1.234,56")).toBeInTheDocument();
  });

  it("FE-01d exibe imagem quando fornecida", () => {
    render(<ProductCard product={mockProduct} />);
    const img = screen.getByAltText("Laptop") as HTMLImageElement;
    expect(img).toBeInTheDocument();
    expect(img.src).toBe("https://example.com/laptop.jpg");
  });

  it("FE-01d exibe ícone quando sem imagem", () => {
    const productNoImage: CatalogProduct = {
      ...mockProduct,
      image_url: undefined,
    };
    render(<ProductCard product={productNoImage} />);
    expect(screen.queryByAltText("Laptop")).not.toBeInTheDocument();
  });

  it("FE-02a mostra selo Esgotado quando disponível é 0", () => {
    const soldOut: CatalogProduct = { ...mockProduct, available: 0 };
    render(<ProductCard product={soldOut} />);
    expect(screen.getByText("Esgotado")).toBeInTheDocument();
  });

  it("FE-02b mostra Em estoque quando disponível > 0", () => {
    render(<ProductCard product={mockProduct} />);
    expect(screen.getByText("Em estoque")).toBeInTheDocument();
  });
});
```

**Teste que prova:** `npm run test -- --run` passa com os testes do ProductCard (todos verdes).

---

## T006: Catalog - lista, carregamento, erro, refetch (45 min)
**Regras:** FE-01a, FE-03, FE-05c

Criar `web/src/features/catalog/Catalog.tsx`:

```typescript
import React, { useEffect, useState } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { AlertCircle, RotateCcw } from "lucide-react";
import { fetchProducts, type CatalogProduct } from "@/lib/api";
import { ProductCard } from "./ProductCard";

export function Catalog() {
  const [products, setProducts] = useState<CatalogProduct[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadProducts = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchProducts();
      if (data.length === 0) {
        setError("empty");
      } else {
        setProducts(data);
      }
    } catch (err) {
      setError("Não foi possível carregar os produtos");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProducts();
  }, []);

  const prefersReducedMotion = window.matchMedia(
    "(prefers-reduced-motion: reduce)"
  ).matches;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      {/* FE-03a Carregamento */}
      {loading && (
        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <Skeleton key={i} className="aspect-square" />
          ))}
        </div>
      )}

      {/* FE-03c Erro */}
      {error && error !== "empty" && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
          <Button
            variant="outline"
            size="sm"
            onClick={loadProducts}
            className="mt-4"
          >
            <RotateCcw className="w-4 h-4 mr-2" />
            Tentar de novo
          </Button>
        </Alert>
      )}

      {/* FE-03b Vazio */}
      {!loading && error === "empty" && (
        <Alert>
          <AlertDescription>Nenhum produto disponível no momento</AlertDescription>
        </Alert>
      )}

      {/* FE-01a, FE-05c Lista com animações */}
      {!loading && !error && products.length > 0 && (
        <div
          className={`grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4 ${
            !prefersReducedMotion ? "animate-fade-in" : ""
          }`}
        >
          {products.map((product, idx) => (
            <div
              key={product.id}
              className={!prefersReducedMotion ? "animate-slide-up" : ""}
              style={{
                animationDelay: `${idx * 50}ms`,
              }}
            >
              <ProductCard product={product} />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
```

Adicionar animações em `web/src/index.css`:

```css
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@layer utilities {
  @media (prefers-reduced-motion: no-preference) {
    .animate-fade-in {
      animation: fadeIn 0.3s ease-in-out;
    }

    .animate-slide-up {
      animation: slideUp 0.4s ease-out forwards;
    }
  }
}
```

Criar `web/src/features/catalog/__tests__/Catalog.test.tsx`:

```typescript
import { describe, it, expect, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Catalog } from "../Catalog";
import { server } from "@/mocks/server";
import { http, HttpResponse } from "msw";

describe("Catalog", () => {
  it("FE-01a exibe produtos do catálogo", async () => {
    render(<Catalog />);
    expect(screen.getByText(/Carregando/i)).toBeInTheDocument(); // Skeleton ou "Carregando"
    await waitFor(() => {
      expect(screen.getByText("Laptop")).toBeInTheDocument();
      expect(screen.getByText("Mouse")).toBeInTheDocument();
      expect(screen.getByText("Teclado")).toBeInTheDocument();
    });
  });

  it("FE-03a exibe indicador de carregamento enquanto busca", async () => {
    render(<Catalog />);
    const skeletons = screen.getAllByRole("presentation"); // Skeleton renderiza assim
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it("FE-03b exibe mensagem quando lista está vazia", async () => {
    server.use(
      http.get("/api/products", () =>
        HttpResponse.json({ products: [] })
      )
    );
    render(<Catalog />);
    await waitFor(() => {
      expect(
        screen.getByText("Nenhum produto disponível no momento")
      ).toBeInTheDocument();
    });
  });

  it("FE-03c exibe erro e botão ao falhar", async () => {
    server.use(
      http.get("/api/products", () => HttpResponse.error())
    );
    render(<Catalog />);
    await waitFor(() => {
      expect(
        screen.getByText("Não foi possível carregar os produtos")
      ).toBeInTheDocument();
      expect(screen.getByText("Tentar de novo")).toBeInTheDocument();
    });
  });

  it("FE-03c refetch ao clicar Tentar de novo", async () => {
    const user = userEvent.setup();
    server.use(
      http.get("/api/products", () => HttpResponse.error())
    );
    render(<Catalog />);
    await waitFor(() =>
      expect(screen.getByText("Tentar de novo")).toBeInTheDocument()
    );

    server.use(
      http.get("/api/products", () =>
        HttpResponse.json({
          products: [
            {
              id: "123",
              name: "Laptop",
              price_cents: 299900,
              available: 5,
            },
          ],
        })
      )
    );

    await user.click(screen.getByText("Tentar de novo"));
    await waitFor(() => {
      expect(screen.getByText("Laptop")).toBeInTheDocument();
    });
  });

  it("FE-05c respeita prefers-reduced-motion", () => {
    // Mock matchMedia
    Object.defineProperty(window, "matchMedia", {
      writable: true,
      value: () => ({ matches: true }),
    });
    const { container } = render(<Catalog />);
    // Sem animações, elemento tem classe normal (sem animate-*)
    expect(container.querySelector(".animate-fade-in")).not.toBeInTheDocument();
  });
});
```

Atualizar `web/src/main.tsx` para usar `MainLayout` + `Catalog`:

```typescript
import React from "react";
import ReactDOM from "react-dom/client";
import { MainLayout } from "./layouts/MainLayout";
import { Catalog } from "./features/catalog/Catalog";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <MainLayout>
      <Catalog />
    </MainLayout>
  </React.StrictMode>
);
```

**Teste que prova:** `npm run test -- --run` passa com todos os testes (10+ verdes); `make web-check` passa (lint + typecheck + testes).

---

## Cobertura de regras

| ID | Tarefa | Teste |
| --- | --- | --- |
| RF-02 | T003 | `it("FE-01a busca 3 produtos...")` |
| FE-01 | T002, T005, T006 | `it("FE-01 exibe nome e preço...")` |
| FE-02 | T005 | `it("FE-02a mostra selo Esgotado...")` |
| FE-03 | T006 | `it("FE-03a exibe indicador...")`, `it("FE-03b exibe mensagem vazia...")`, `it("FE-03c exibe erro...")` |
| FE-05 | T004, T006 | `it("FE-05a exibe cabeçalho...")`, `it("FE-05c respeita prefers-reduced-motion...")` |

Todas as regras têm pelo menos uma tarefa.

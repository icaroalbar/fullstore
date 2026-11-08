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

  const loadProducts = React.useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchProducts();
      if (data.length === 0) {
        setError("empty");
      } else {
        setProducts(data);
      }
    } catch {
      setError("Não foi possível carregar os produtos");
    } finally {
      setLoading(false);
    }
  }, []);

  // eslint-disable-next-line react-hooks/exhaustive-deps, react-hooks/set-state-in-effect
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

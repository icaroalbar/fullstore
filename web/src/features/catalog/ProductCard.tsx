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

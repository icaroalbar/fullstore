import { http, HttpResponse } from "msw";
import type { CatalogProduct } from "@/lib/api";

const mockProducts: CatalogProduct[] = [
  {
    id: "123e4567-e89b-12d3-a456-426614174000",
    name: "Laptop",
    price_cents: 299900,
    available: 5,
    image_url: "/mock/laptop.svg",
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
    image_url: "/mock/teclado.svg",
  },
];

export const handlers = [
  http.get("*/api/products", () => {
    return HttpResponse.json({ products: mockProducts });
  }),
];

import { describe, it, expect } from "vitest";
import { fetchProducts } from "../api";
import { server } from "@/mocks/server";
import { http, HttpResponse } from "msw";

describe("api", () => {
  it("FE-01a busca 3 produtos (nomes, preços, Mouse esgotado sem imagem)", async () => {
    const products = await fetchProducts();
    expect(products).toHaveLength(3);
    expect(products[0]).toMatchObject({
      name: "Laptop",
      price_cents: 299900,
      available: 5,
      image_url: "/mock/laptop.svg",
    });
    expect(products[1]).toMatchObject({
      name: "Mouse",
      price_cents: 4990,
      available: 0,
    });
    expect(products[1].image_url).toBeUndefined();
  });

  it("FE-01a respeita contrato: todo produto tem id, name, price_cents, available", async () => {
    const products = await fetchProducts();
    products.forEach((p) => {
      expect(p.id).toBeDefined();
      expect(p.name).toBeDefined();
      expect(p.price_cents).toBeDefined();
      expect(p.available).toBeDefined();
    });
  });

  it("FE-03c erro na API (500) lança erro em fetchProducts", async () => {
    server.use(
      http.get("*/api/products", () => new HttpResponse(null, { status: 500 }))
    );
    await expect(fetchProducts()).rejects.toThrow();
  });

  it("erro de rede (sem resposta) lança erro", async () => {
    server.use(http.get("*/api/products", () => HttpResponse.error()));
    await expect(fetchProducts()).rejects.toThrow();
  });
});

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
    expect(screen.getByText(/R\$\s1\.234,56/)).toBeInTheDocument();
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

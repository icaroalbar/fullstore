export interface CatalogProduct {
  id: string;
  name: string;
  price_cents: number;
  available: number;
  image_url?: string;
}

export async function fetchProducts(): Promise<CatalogProduct[]> {
  const API_URL = import.meta.env.VITE_API_URL ?? "/api";
  const response = await fetch(`${API_URL}/products`);
  if (!response.ok) {
    throw new Error("Failed to fetch products");
  }
  const data = await response.json();
  return data.products;
}

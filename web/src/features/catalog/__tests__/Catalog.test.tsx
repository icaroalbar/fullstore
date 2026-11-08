import { describe, it, expect } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Catalog } from "../Catalog";
import { server } from "@/mocks/server";
import { http, HttpResponse } from "msw";

describe("Catalog", () => {
  it("FE-01a exibe produtos do catálogo", async () => {
    render(<Catalog />);
    await waitFor(() => {
      expect(screen.getByText("Laptop")).toBeInTheDocument();
      expect(screen.getByText("Mouse")).toBeInTheDocument();
      expect(screen.getByText("Teclado")).toBeInTheDocument();
    });
  });

  it("FE-03a exibe indicador de carregamento enquanto busca", async () => {
    render(<Catalog />);
    const skeletons = document.querySelectorAll('[data-slot="skeleton"]');
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
    Object.defineProperty(window, "matchMedia", {
      writable: true,
      value: () => ({ matches: true }),
    });
    const { container } = render(<Catalog />);
    expect(container.querySelector(".animate-fade-in")).not.toBeInTheDocument();
  });
});

import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MainLayout } from "../MainLayout";

describe("MainLayout", () => {
  it("FE-05a exibe cabeçalho com logo FullStore e rodapé", () => {
    render(
      <MainLayout>
        <div>Conteúdo teste</div>
      </MainLayout>
    );
    expect(screen.getByRole("banner")).toBeInTheDocument();
    expect(screen.getByText("FullStore")).toBeInTheDocument();
    expect(screen.getByRole("contentinfo")).toBeInTheDocument();
  });

  it("FE-05b cabeçalho é fixo (sticky top-0)", () => {
    const { container } = render(
      <MainLayout>
        <div>Conteúdo</div>
      </MainLayout>
    );
    const header = container.querySelector("header");
    expect(header).toHaveClass("sticky");
    expect(header).toHaveClass("top-0");
  });

  it("conteúdo (children) aparece dentro do main", () => {
    render(
      <MainLayout>
        <div data-testid="conteudo-teste">Conteúdo</div>
      </MainLayout>
    );
    const main = screen.getByRole("main");
    expect(main).toContainElement(screen.getByTestId("conteudo-teste"));
  });

  it("botão do carrinho tem nome acessível", () => {
    render(
      <MainLayout>
        <div>Conteúdo</div>
      </MainLayout>
    );
    const carrinhoBtn = screen.getByRole("button", { name: /carrinho/i });
    expect(carrinhoBtn).toBeInTheDocument();
    expect(carrinhoBtn.querySelector("svg")).toBeInTheDocument();
  });

  it("campo de busca existe e está desabilitado (F03)", () => {
    render(
      <MainLayout>
        <div>Conteúdo</div>
      </MainLayout>
    );
    const searchInputs = screen.getAllByPlaceholderText(/buscar/i);
    searchInputs.forEach((input) => {
      expect(input).toBeDisabled();
    });
  });
});

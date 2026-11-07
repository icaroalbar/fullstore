import { describe, it, expect } from "vitest";
import { formatPrice } from "../money";

describe("money", () => {
  it("FE-01b formata 123456 centavos como R$ 1.234,56", () => {
    // Intl.NumberFormat("pt-BR", currency) usa espaço não-separável (U+00A0) entre "R$" e número.
    // O formatPrice retorna a string exatamente como Intl.NumberFormat gera. Não normalizamos.
    const result = formatPrice(123456);
    expect(result).toMatch(/R.\s1\.234,56/); // Regex com \s (qualquer espaço) para aceitar U+00A0
    expect(result).toContain("1.234,56");
  });

  it("FE-01c formata 990 centavos como R$ 9,90", () => {
    const result = formatPrice(990);
    expect(result).toMatch(/R.\s9,90/);
    expect(result).toContain("9,90");
  });

  it("formata 100 centavos como R$ 1,00", () => {
    const result = formatPrice(100);
    expect(result).toMatch(/R.\s1,00/);
    expect(result).toContain("1,00");
  });

  it("formata 0 centavos como R$ 0,00", () => {
    const result = formatPrice(0);
    expect(result).toMatch(/R.\s0,00/);
    expect(result).toContain("0,00");
  });

  it("formata 100000000 centavos (R$ 1 milhão) sem perder precisão", () => {
    const result = formatPrice(100000000);
    expect(result).toContain("1.000.000,00");
  });
});

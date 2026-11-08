import React, { useState } from "react";
import { ShoppingCart, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";

export interface MainLayoutProps {
  children: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  const [cartCount] = useState(0);

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-background/95 border-b border-border backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="max-w-7xl mx-auto px-4 py-4">
          {/* Desktop: 1 linha com logo, busca, carrinho */}
          <div className="flex items-center gap-4">
            {/* Logo */}
            <h1 className="text-2xl font-bold text-primary flex-shrink-0">
              FullStore
            </h1>

            {/* Busca (desktop) */}
            <div className="hidden md:flex flex-1 max-w-md">
              <div className="relative w-full">
                <Search className="absolute left-2 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  placeholder="Buscar produtos..."
                  disabled
                  className="pl-8"
                  aria-label="Buscar"
                />
              </div>
            </div>

            {/* Carrinho (desktop) */}
            <Button
              variant="ghost"
              size="icon"
              className="relative flex-shrink-0"
              aria-label="Carrinho"
            >
              <ShoppingCart className="w-5 h-5" />
              {cartCount > 0 && (
                <Badge className="absolute -top-2 -right-2" variant="default">
                  {cartCount}
                </Badge>
              )}
            </Button>
          </div>

          {/* Busca (mobile) */}
          <div className="md:hidden mt-3">
            <div className="relative w-full">
              <Search className="absolute left-2 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
              <Input
                placeholder="Buscar..."
                disabled
                className="pl-8"
                aria-label="Buscar"
              />
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1">
        <div className="max-w-7xl mx-auto w-full">{children}</div>
      </main>

      {/* Footer */}
      <footer className="bg-background border-t border-border mt-12">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
            {/* Ajuda */}
            <div>
              <h3 className="font-semibold mb-3">Ajuda</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Perguntas Frequentes
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Contato
                  </a>
                </li>
              </ul>
            </div>

            {/* Sobre */}
            <div>
              <h3 className="font-semibold mb-3">Sobre</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Sobre a FullStore
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Blog
                  </a>
                </li>
              </ul>
            </div>

            {/* Políticas */}
            <div>
              <h3 className="font-semibold mb-3">Políticas</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Privacidade
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Termos de Uso
                  </a>
                </li>
              </ul>
            </div>
          </div>

          {/* Copyright */}
          <div className="border-t border-border pt-6 text-center text-sm text-muted-foreground">
            <p>© 2026 FullStore. Todos os direitos reservados.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}

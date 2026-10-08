import * as React from "react";
import Link from "next/link";
import { LogOut } from "lucide-react";

import { Avatar, AvatarFallback, iniciais } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { RelpLogo } from "@/components/relp/RelpLogo";
import { SearchBar } from "@/components/relp/SearchBar";
import { cn } from "@/lib/utils";

export type AppHeaderProps = {
  /** `loading` = ainda perguntando `GET /auth/me`; evita piscar "Entrar" para quem já está logado. */
  state: "loading" | "guest" | "user";
  user?: { nome: string; email: string } | null;
  onSearch?: (termo: string) => void;
  onLogout?: () => void;
  loggingOut?: boolean;
  className?: string;
};

/** Cabeçalho único do app: logo, navegação, busca e área de sessão. Puro (sem rede): estados via props. */
export function AppHeader({ state, user, onSearch, onLogout, loggingOut, className }: AppHeaderProps) {
  return (
    <header className={cn("sticky top-0 z-40 border-b border-border bg-card/95 backdrop-blur", className)}>
      <div className="mx-auto flex max-w-5xl flex-wrap items-center gap-x-4 gap-y-2 px-4 py-3">
        <Link
          href="/"
          aria-label="Relp — página inicial"
          className="rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
        >
          <RelpLogo />
        </Link>

        <nav aria-label="Principal" className="hidden sm:block">
          <Link
            href="/"
            className="rounded-md px-2 py-1 text-sm font-medium text-muted-foreground transition-colors hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            Trilhas
          </Link>
        </nav>

        {/* Busca: linha própria no mobile, centralizada a partir de md. */}
        <div className="order-last w-full md:order-none md:ml-2 md:w-auto md:max-w-xs md:flex-1">
          <SearchBar onSubmit={onSearch} />
        </div>

        <div className="ml-auto flex items-center gap-2">
          {state === "loading" && (
            <div role="status" aria-label="Verificando sessão" className="flex items-center gap-2">
              <Skeleton className="h-9 w-9 rounded-full" />
              <Skeleton className="hidden h-9 w-16 sm:block" />
            </div>
          )}

          {state === "guest" && (
            <>
              <Button asChild variant="ghost" size="sm" className="hidden sm:inline-flex">
                <Link href="/registro">Criar conta</Link>
              </Button>
              <Button asChild size="sm">
                <Link href="/login">Entrar</Link>
              </Button>
            </>
          )}

          {state === "user" && user && (
            <>
              <div className="flex items-center gap-2">
                <Avatar>
                  <AvatarFallback aria-hidden="true">{iniciais(user.nome, user.email)}</AvatarFallback>
                </Avatar>
                <span className="hidden max-w-[10rem] truncate text-sm font-medium text-foreground md:inline">
                  {user.nome}
                </span>
              </div>
              <Button variant="ghost" size="sm" onClick={onLogout} loading={loggingOut}>
                <LogOut />
                Sair
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}

"use client";

import * as React from "react";
import { Search } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

type SearchBarProps = {
  defaultValue?: string;
  placeholder?: string;
  /** Nome acessível do campo (diferencie quando houver duas buscas na mesma página). */
  label?: string;
  /** Chamado com o termo já sem espaços nas pontas. Termos vazios não disparam. */
  onSubmit?: (termo: string) => void;
  className?: string;
};

/** Barra de busca (header e página de resultados). Sem lógica de rede: quem usa decide o destino. */
export function SearchBar({
  defaultValue = "",
  placeholder = "Buscar artigos…",
  label = "Buscar artigos",
  onSubmit,
  className,
}: SearchBarProps) {
  const [valor, setValor] = React.useState(defaultValue);

  React.useEffect(() => {
    setValor(defaultValue);
  }, [defaultValue]);

  function enviar(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const termo = valor.trim();
    if (termo) onSubmit?.(termo);
  }

  return (
    <form role="search" onSubmit={enviar} className={cn("relative flex w-full items-center", className)}>
      <Input
        type="search"
        name="q"
        value={valor}
        onChange={(e) => setValor(e.target.value)}
        placeholder={placeholder}
        aria-label={label}
        autoComplete="off"
        className="pr-11"
      />
      <Button
        type="submit"
        variant="ghost"
        size="icon"
        className="absolute right-0 top-0 h-10 w-10 text-muted-foreground hover:text-primary"
        aria-label="Buscar"
      >
        <Search />
      </Button>
    </form>
  );
}

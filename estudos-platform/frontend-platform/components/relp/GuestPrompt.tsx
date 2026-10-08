import * as React from "react";
import Link from "next/link";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";

export type GuestPromptProps = {
  titulo: string;
  descricao: string;
  /** Caminho interno para voltar depois de entrar (ex.: /artigos/slug). */
  next?: string;
  className?: string;
};

/** Convite discreto para visitantes: login/cadastro mantendo o retorno ao conteúdo atual. */
export function GuestPrompt({ titulo, descricao, next, className }: GuestPromptProps) {
  const query = next ? `?next=${encodeURIComponent(next)}` : "";
  return (
    <Alert variant="info" className={className}>
      <AlertTitle>{titulo}</AlertTitle>
      <AlertDescription>
        <p>{descricao}</p>
        <div className="mt-3 flex flex-wrap gap-2">
          <Button asChild size="sm">
            <Link href={`/login${query}`}>Entrar</Link>
          </Button>
          <Button asChild size="sm" variant="outline">
            <Link href={`/registro${query}`}>Criar conta</Link>
          </Button>
        </div>
      </AlertDescription>
    </Alert>
  );
}

"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { WifiOff } from "lucide-react";

import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";

/** Erro de carregamento (API fora do ar, 5xx). Sem stack na tela; "Tentar de novo" refaz a busca no servidor. */
export default function ErrorPage({ reset }: { error: Error & { digest?: string }; reset: () => void }) {
  const router = useRouter();
  const [pendente, startTransition] = React.useTransition();

  function tentarDeNovo() {
    startTransition(() => {
      router.refresh();
      reset();
    });
  }

  return (
    <EmptyState
      tone="danger"
      icon={<WifiOff />}
      title="Não conseguimos carregar esta página"
      description="O servidor pode estar indisponível no momento. Verifique sua conexão e tente novamente."
      action={
        <Button onClick={tentarDeNovo} loading={pendente}>
          Tentar de novo
        </Button>
      }
    />
  );
}

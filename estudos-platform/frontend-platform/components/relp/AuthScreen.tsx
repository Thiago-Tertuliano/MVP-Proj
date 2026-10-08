"use client";

import * as React from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { toast } from "@/components/ui/sonner";
import { Skeleton } from "@/components/ui/skeleton";
import { AuthForm, type AuthValues } from "@/components/relp/AuthForm";
import { erroLogin, erroRegistro, type FormAlertData } from "@/lib/auth-errors";
import { destinoSeguro } from "@/lib/navigation";
import { useSession } from "@/lib/session";

/** Página de login/registro: liga o formulário à sessão, trata erros e redireciona (?next= seguro). */
export function AuthScreen({ mode }: { mode: "login" | "registro" }) {
  const router = useRouter();
  const params = useSearchParams();
  const { status, user, login, registrar } = useSession();
  const [erro, setErro] = React.useState<FormAlertData | null>(null);
  const [redirecionando, setRedirecionando] = React.useState(false);

  const next = destinoSeguro(params.get("next"));
  const nextQuery = params.get("next") ? `?next=${encodeURIComponent(params.get("next") as string)}` : "";

  // Quem já está logado não precisa ver o formulário.
  React.useEffect(() => {
    if (status === "user" && !redirecionando) router.replace(next);
  }, [status, next, router, redirecionando]);

  async function enviar(values: AuthValues) {
    setErro(null);
    try {
      if (mode === "login") {
        const u = await login(values.email.trim(), values.senha);
        toast.success(`Olá, ${u.nome.split(" ")[0]}!`);
      } else {
        const u = await registrar((values.nome ?? "").trim(), values.email.trim(), values.senha);
        toast.success(`Conta criada. Bem-vindo(a), ${u.nome.split(" ")[0]}!`);
      }
      setRedirecionando(true);
      router.replace(next);
      router.refresh();
    } catch (err) {
      setErro(mode === "login" ? erroLogin(err) : erroRegistro(err));
    }
  }

  if (status === "loading" || status === "user") {
    return (
      <div className="space-y-4" role="status" aria-label="Carregando">
        <Skeleton className="h-8 w-40" />
        <Skeleton className="h-72 w-full" />
        {user && <span className="sr-only">Redirecionando…</span>}
      </div>
    );
  }

  return (
    <AuthForm
      mode={mode}
      onSubmit={enviar}
      serverError={erro}
      alternateHref={`${mode === "login" ? "/registro" : "/login"}${nextQuery}`}
    />
  );
}

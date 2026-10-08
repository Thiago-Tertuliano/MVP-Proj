import type { Metadata } from "next";
import { Suspense } from "react";

import { AuthScreen } from "@/components/relp/AuthScreen";

export const metadata: Metadata = { title: "Entrar" };

export default function LoginPage() {
  return (
    <div className="mx-auto max-w-md space-y-6">
      <header className="space-y-2 text-center">
        <h1 className="text-2xl font-bold tracking-tight text-foreground">Entrar</h1>
        <p className="text-sm text-muted-foreground">Acesse sua conta para salvar progresso e anotações.</p>
      </header>
      <Suspense>
        <AuthScreen mode="login" />
      </Suspense>
    </div>
  );
}

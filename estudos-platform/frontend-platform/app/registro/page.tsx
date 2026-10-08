import type { Metadata } from "next";
import { Suspense } from "react";

import { AuthScreen } from "@/components/relp/AuthScreen";

export const metadata: Metadata = { title: "Criar conta" };

export default function RegistroPage() {
  return (
    <div className="mx-auto max-w-md space-y-6">
      <header className="space-y-2 text-center">
        <h1 className="text-2xl font-bold tracking-tight text-foreground">Criar conta</h1>
        <p className="text-sm text-muted-foreground">É grátis. Leva menos de um minuto.</p>
      </header>
      <Suspense>
        <AuthScreen mode="registro" />
      </Suspense>
    </div>
  );
}

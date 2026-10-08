"use client";

import * as React from "react";
import { useRouter } from "next/navigation";

import { toast } from "@/components/ui/sonner";
import { AppHeader } from "@/components/relp/AppHeader";
import { useSession } from "@/lib/session";

/** Header ligado à sessão: busca navega para /busca?q=, "Sair" encerra a sessão e volta à home. */
export function SessionHeader() {
  const { status, user, logout } = useSession();
  const router = useRouter();
  const [saindo, setSaindo] = React.useState(false);

  async function sair() {
    setSaindo(true);
    try {
      await logout();
      toast.success("Você saiu da sua conta.");
      router.push("/");
      router.refresh();
    } finally {
      setSaindo(false);
    }
  }

  return (
    <AppHeader
      state={status}
      user={user}
      onSearch={(termo) => router.push(`/busca?q=${encodeURIComponent(termo)}`)}
      onLogout={sair}
      loggingOut={saindo}
    />
  );
}

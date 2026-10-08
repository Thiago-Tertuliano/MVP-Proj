"use client";

import * as React from "react";

import { api, ApiError, onSessionExpired } from "@/lib/api";
import { limparProgresso } from "@/lib/progress-store";
import type { AuthResposta, Usuario } from "@/lib/types";
import { toast } from "@/components/ui/sonner";

export type SessionStatus = "loading" | "guest" | "user";

type SessionValue = {
  status: SessionStatus;
  user: Usuario | null;
  login: (email: string, senha: string) => Promise<Usuario>;
  registrar: (nome: string, email: string, senha: string) => Promise<Usuario>;
  logout: () => Promise<void>;
};

const SessionContext = React.createContext<SessionValue | null>(null);

/**
 * Dica NÃO sensível ("já entrei neste navegador"). O JWT continua só em cookie HttpOnly.
 * Serve para só tentar refresh silencioso de quem realmente teve sessão — visitantes não
 * gastam o rate limit de /auth/*.
 */
const HINT_KEY = "relp:sessao";

function lerHint(): boolean {
  try {
    return window.localStorage.getItem(HINT_KEY) === "1";
  } catch {
    return false;
  }
}
function gravarHint(ativo: boolean) {
  try {
    if (ativo) window.localStorage.setItem(HINT_KEY, "1");
    else window.localStorage.removeItem(HINT_KEY);
  } catch {
    /* modo privado / storage bloqueado: segue sem dica */
  }
}

export function SessionProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = React.useState<{ status: SessionStatus; user: Usuario | null }>({
    status: "loading",
    user: null,
  });
  const statusRef = React.useRef<SessionStatus>("loading");
  statusRef.current = state.status;

  // GET /auth/me ao carregar o app: define se o header mostra "Entrar" ou o avatar.
  React.useEffect(() => {
    let ativo = true;
    api<Usuario>("/auth/me", { refresh: lerHint() })
      .then((user) => {
        if (!ativo) return;
        gravarHint(true);
        setState({ status: "user", user });
      })
      .catch(() => {
        if (!ativo) return;
        gravarHint(false);
        setState({ status: "guest", user: null });
      });
    return () => {
      ativo = false;
    };
  }, []);

  // Refresh falhou no meio da navegação: volta a visitante sem tela quebrada.
  React.useEffect(() => {
    return onSessionExpired(() => {
      if (statusRef.current !== "user") return;
      gravarHint(false);
      limparProgresso();
      setState({ status: "guest", user: null });
      toast.warning("Sua sessão expirou.", { description: "Entre novamente para continuar de onde parou." });
    });
  }, []);

  const value = React.useMemo<SessionValue>(
    () => ({
      ...state,
      async login(email, senha) {
        const resp = await api<AuthResposta>("/auth/login", { method: "POST", body: { email, senha } });
        gravarHint(true);
        setState({ status: "user", user: resp.usuario });
        return resp.usuario;
      },
      async registrar(nome, email, senha) {
        const resp = await api<AuthResposta>("/auth/registrar", { method: "POST", body: { nome, email, senha } });
        gravarHint(true);
        setState({ status: "user", user: resp.usuario });
        return resp.usuario;
      },
      async logout() {
        try {
          await api("/auth/logout", { method: "POST" });
        } catch (err) {
          // Mesmo se o servidor não confirmar, o front deixa de tratar o usuário como logado.
          if (!(err instanceof ApiError) || !err.isUnauthorized) {
            toast.warning("Não foi possível confirmar o logout no servidor.", {
              description: "Você saiu neste navegador. Se estiver em computador compartilhado, tente de novo.",
            });
          }
        } finally {
          gravarHint(false);
          limparProgresso();
          setState({ status: "guest", user: null });
        }
      },
    }),
    [state],
  );

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
}

export function useSession(): SessionValue {
  const ctx = React.useContext(SessionContext);
  if (!ctx) throw new Error("useSession deve ser usado dentro de <SessionProvider>");
  return ctx;
}

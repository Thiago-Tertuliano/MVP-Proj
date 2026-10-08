"use client";

import * as React from "react";

import { api } from "@/lib/api";
import type { Anotacao } from "@/lib/types";

export type NotaStatus = "idle" | "saving" | "saved" | "error";

export const ANOTACAO_DEBOUNCE_MS = 1000;

/**
 * Anotação do aluno em um artigo: GET ao abrir, PUT com debounce enquanto digita.
 * Última edição vence; ao sair da página o que estiver pendente é enviado na hora.
 */
export function useAnotacao(artigoId: string, habilitado: boolean) {
  const [texto, setTexto] = React.useState("");
  const [carregando, setCarregando] = React.useState(habilitado);
  const [erroCarga, setErroCarga] = React.useState(false);
  const [status, setStatus] = React.useState<NotaStatus>("idle");

  const textoRef = React.useRef("");
  const timer = React.useRef<ReturnType<typeof setTimeout> | null>(null);
  const pendente = React.useRef(false);
  const versao = React.useRef(0);

  const enviar = React.useCallback(
    async (conteudo: string) => {
      const minha = ++versao.current;
      pendente.current = false;
      setStatus("saving");
      try {
        await api<Anotacao>(`/artigos/${artigoId}/anotacoes`, {
          method: "PUT",
          body: { conteudo: { texto: conteudo } },
        });
        if (minha === versao.current) setStatus("saved");
      } catch {
        if (minha === versao.current) setStatus("error");
      }
    },
    [artigoId],
  );

  // Carga inicial.
  React.useEffect(() => {
    if (!habilitado) {
      setCarregando(false);
      return;
    }
    const ac = new AbortController();
    setCarregando(true);
    setErroCarga(false);
    api<Anotacao>(`/artigos/${artigoId}/anotacoes`, { signal: ac.signal })
      .then((a) => {
        const t = typeof a.conteudo?.texto === "string" ? a.conteudo.texto : "";
        // Se o aluno já começou a digitar antes da carga terminar, não sobrescreve.
        if (!pendente.current && versao.current === 0) {
          textoRef.current = t;
          setTexto(t);
        }
        setCarregando(false);
      })
      .catch(() => {
        if (ac.signal.aborted) return;
        setErroCarga(true);
        setCarregando(false);
      });
    return () => ac.abort();
  }, [artigoId, habilitado]);

  // Ao sair/trocar de artigo: envia o que ficou pendente.
  React.useEffect(() => {
    return () => {
      if (timer.current) clearTimeout(timer.current);
      if (pendente.current && habilitado) {
        pendente.current = false;
        void api(`/artigos/${artigoId}/anotacoes`, {
          method: "PUT",
          body: { conteudo: { texto: textoRef.current } },
        }).catch(() => undefined);
      }
    };
  }, [artigoId, habilitado]);

  const alterar = React.useCallback(
    (novo: string) => {
      textoRef.current = novo;
      setTexto(novo);
      pendente.current = true;
      setStatus("idle");
      if (timer.current) clearTimeout(timer.current);
      timer.current = setTimeout(() => void enviar(textoRef.current), ANOTACAO_DEBOUNCE_MS);
    },
    [enviar],
  );

  const tentarNovamente = React.useCallback(() => {
    if (timer.current) clearTimeout(timer.current);
    void enviar(textoRef.current);
  }, [enviar]);

  return { texto, alterar, status, carregando, erroCarga, tentarNovamente };
}

"use client";

import * as React from "react";
import { AlertCircle, Check, Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export type NoteEditorProps = {
  value: string;
  onChange: (valor: string) => void;
  status?: "idle" | "saving" | "saved" | "error";
  /** Carregando a anotação existente. */
  loading?: boolean;
  /** Falha ao carregar (ainda dá para escrever e salvar). */
  loadError?: boolean;
  onRetry?: () => void;
  disabled?: boolean;
  className?: string;
};

/** Área de anotações com estado de salvamento (idle / salvando / salvo / erro) e reenvio. */
export function NoteEditor({
  value,
  onChange,
  status = "idle",
  loading,
  loadError,
  onRetry,
  disabled,
  className,
}: NoteEditorProps) {
  const id = React.useId();
  const statusId = `${id}-status`;

  return (
    <div className={cn("space-y-2", className)}>
      <div className="flex items-center justify-between gap-2">
        <Label htmlFor={id}>Minhas anotações</Label>
        <p id={statusId} role="status" aria-live="polite" className="flex items-center gap-1.5 text-xs">
          {status === "saving" && (
            <>
              <Loader2 className="h-3 w-3 animate-spin text-muted-foreground" aria-hidden="true" />
              <span className="text-muted-foreground">Salvando…</span>
            </>
          )}
          {status === "saved" && (
            <>
              <Check className="h-3.5 w-3.5 text-done-fg" aria-hidden="true" />
              <span className="text-done-fg">Salvo</span>
            </>
          )}
          {status === "error" && (
            <>
              <AlertCircle className="h-3.5 w-3.5 text-danger-fg" aria-hidden="true" />
              <span className="text-danger-fg">Não foi possível salvar.</span>
              {onRetry && (
                <Button type="button" variant="link" size="sm" className="h-auto p-0 text-xs" onClick={onRetry}>
                  Tentar de novo
                </Button>
              )}
            </>
          )}
        </p>
      </div>

      {loading ? (
        <Skeleton className="h-32 w-full" />
      ) : (
        <Textarea
          id={id}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          disabled={disabled}
          rows={5}
          placeholder="Resumos, dúvidas, pontos-chave desta leitura… Salvamos automaticamente."
          aria-describedby={statusId}
          aria-invalid={status === "error" || undefined}
        />
      )}

      {loadError && (
        <p className="text-xs text-warning-fg">
          Não conseguimos carregar suas anotações anteriores. O que você escrever aqui será salvo normalmente.
        </p>
      )}
    </div>
  );
}

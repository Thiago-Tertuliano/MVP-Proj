"use client";

import * as React from "react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { ApiError } from "@/lib/api";
import type { QuizQuestaoRoadmap, QuizResultado } from "@/lib/roadmaps/types";
import { cn } from "@/lib/utils";

export type ChefeQuizProps = {
  questoes: QuizQuestaoRoadmap[];
  /** Envia TODAS as respostas de uma vez; a correção acontece no servidor. */
  onEnviar: (respostas: Record<string, string>) => Promise<QuizResultado>;
  onFechar?: () => void;
  className?: string;
};

function mensagemDeErro(err: unknown): string {
  if (err instanceof ApiError) return err.message;
  return "Não foi possível enviar suas respostas. Tente novamente.";
}

/**
 * Quiz do chefe. O gabarito nunca chega ao navegador antes da aprovação: o servidor devolve só
 * acertos/total e, se passou, a correção comentada. Reprovado = tentar de novo, sem limite.
 */
export function ChefeQuiz({ questoes, onEnviar, onFechar, className }: ChefeQuizProps) {
  const baseId = React.useId();
  const [respostas, setRespostas] = React.useState<Record<string, string>>({});
  const [enviando, setEnviando] = React.useState(false);
  const [resultado, setResultado] = React.useState<QuizResultado | null>(null);
  const [erro, setErro] = React.useState<string | null>(null);

  const completas = questoes.every((q) => !!respostas[q.id]);
  const correcaoPorQuestao = new Map((resultado?.correcao ?? []).map((c) => [c.questao_id, c]));

  async function enviar() {
    setEnviando(true);
    setErro(null);
    try {
      setResultado(await onEnviar(respostas));
    } catch (err) {
      setErro(mensagemDeErro(err));
    } finally {
      setEnviando(false);
    }
  }

  function tentarDeNovo() {
    setResultado(null);
    setRespostas({});
    setErro(null);
  }

  const travado = enviando || !!resultado?.aprovado;

  return (
    <div className={cn("space-y-6", className)}>
      {questoes.map((q, i) => {
        const correcao = correcaoPorQuestao.get(q.id);
        return (
          <fieldset key={q.id} className="space-y-3" disabled={travado}>
            <legend className="text-sm font-medium text-foreground">
              <span className="mr-1.5 text-muted-foreground">{i + 1}.</span>
              {q.enunciado}
            </legend>
            <RadioGroup
              value={respostas[q.id] ?? ""}
              onValueChange={(v) => setRespostas((r) => ({ ...r, [q.id]: v }))}
              aria-label={q.enunciado}
            >
              {q.opcoes.map((op) => {
                const itemId = `${baseId}-${q.id}-${op.id}`;
                const marcada = respostas[q.id] === op.id;
                return (
                  <div
                    key={op.id}
                    className={cn(
                      "flex items-start gap-3 rounded-md border border-border px-3 py-2 transition-colors",
                      marcada && !resultado && "border-primary bg-primary-muted",
                      marcada && correcao?.correta && "border-done bg-done-muted",
                      marcada && correcao && !correcao.correta && "border-danger bg-danger-muted",
                    )}
                  >
                    <RadioGroupItem value={op.id} id={itemId} className="mt-0.5" />
                    <Label htmlFor={itemId} className="flex-1 cursor-pointer text-sm font-normal leading-snug">
                      {op.texto}
                    </Label>
                  </div>
                );
              })}
            </RadioGroup>
            {correcao && (
              <p
                className={cn("text-sm", correcao.correta ? "text-done-fg" : "text-danger-fg")}
                role="status"
              >
                <strong>{correcao.correta ? "Correto." : "Incorreto."}</strong>
                {correcao.explicacao ? ` ${correcao.explicacao}` : ""}
              </p>
            )}
          </fieldset>
        );
      })}

      {erro && (
        <Alert variant="danger">
          <AlertTitle>Não foi possível corrigir</AlertTitle>
          <AlertDescription>{erro}</AlertDescription>
        </Alert>
      )}

      {resultado && (
        <Alert variant={resultado.aprovado ? "success" : "danger"}>
          <AlertTitle>{resultado.aprovado ? "Chefe derrotado!" : "O chefe resistiu"}</AlertTitle>
          <AlertDescription>
            Você acertou {resultado.acertos} de {resultado.total}.{" "}
            {resultado.aprovado
              ? "Veja a correção acima."
              : `São necessários ${resultado.nota_minima_pct}% para vencer. Revise o conteúdo e tente de novo.`}
          </AlertDescription>
        </Alert>
      )}

      <div className="flex flex-wrap items-center gap-2">
        {!resultado && (
          <Button type="button" onClick={enviar} loading={enviando} disabled={!completas}>
            Enviar respostas
          </Button>
        )}
        {resultado && !resultado.aprovado && (
          <Button type="button" onClick={tentarDeNovo}>
            Tentar de novo
          </Button>
        )}
        {resultado?.aprovado && onFechar && (
          <Button type="button" variant="success" onClick={onFechar}>
            Continuar
          </Button>
        )}
        {!resultado && !completas && (
          <span className="text-xs text-muted-foreground">Responda todas as questões para enviar.</span>
        )}
      </div>
    </div>
  );
}

"use client";

import * as React from "react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import type { QuizQuestao } from "@/lib/types";
import { cn } from "@/lib/utils";

type QuestaoProps = { questao: QuizQuestao; indice: number };

function Questao({ questao, indice }: QuestaoProps) {
  const baseId = React.useId();
  const [escolha, setEscolha] = React.useState<string>("");
  const [verificada, setVerificada] = React.useState(false);

  const acertou = verificada && escolha === questao.correta;

  function reiniciar() {
    setEscolha("");
    setVerificada(false);
  }

  return (
    <fieldset className="space-y-3">
      <legend className="text-sm font-medium text-foreground">
        <span className="mr-1.5 text-muted-foreground">{indice}.</span>
        {questao.enunciado}
      </legend>

      <RadioGroup
        value={escolha}
        onValueChange={(v) => {
          setEscolha(v);
          setVerificada(false);
        }}
        aria-label={questao.enunciado}
      >
        {questao.opcoes.map((op) => {
          const itemId = `${baseId}-${op.id}`;
          const marcadaCerta = verificada && op.id === questao.correta;
          const marcadaErrada = verificada && op.id === escolha && op.id !== questao.correta;
          return (
            <div
              key={op.id}
              className={cn(
                "flex items-start gap-3 rounded-md border border-border px-3 py-2 transition-colors",
                escolha === op.id && !verificada && "border-primary bg-primary-muted",
                marcadaCerta && "border-done bg-done-muted",
                marcadaErrada && "border-danger bg-danger-muted",
              )}
            >
              <RadioGroupItem value={op.id} id={itemId} className="mt-0.5" />
              <Label htmlFor={itemId} className="flex-1 cursor-pointer text-sm font-normal leading-snug">
                {op.texto}
                {marcadaCerta && <span className="sr-only"> (resposta correta)</span>}
                {marcadaErrada && <span className="sr-only"> (sua resposta, incorreta)</span>}
              </Label>
            </div>
          );
        })}
      </RadioGroup>

      <div className="flex flex-wrap items-center gap-2">
        <Button type="button" size="sm" onClick={() => setVerificada(true)} disabled={!escolha || verificada}>
          Verificar
        </Button>
        {verificada && (
          <Button type="button" size="sm" variant="ghost" onClick={reiniciar}>
            Tentar de novo
          </Button>
        )}
      </div>

      {verificada && (
        <Alert variant={acertou ? "success" : "danger"}>
          <AlertTitle>{acertou ? "Correto!" : "Ainda não foi dessa vez"}</AlertTitle>
          {questao.explicacao && <AlertDescription>{questao.explicacao}</AlertDescription>}
        </Alert>
      )}
    </fieldset>
  );
}

export type QuizPanelProps = {
  questoes: QuizQuestao[];
  className?: string;
};

/** Quiz de autoavaliação. A correção é local (o gabarito vem em `metadados.quiz`): não gera nota nem progresso. */
export function QuizPanel({ questoes, className }: QuizPanelProps) {
  if (!questoes?.length) return null;
  return (
    <aside aria-labelledby="quiz-titulo" className={cn("space-y-5 rounded-xl border border-border bg-card p-5", className)}>
      <div className="space-y-1">
        <h2 id="quiz-titulo" className="text-base font-semibold text-foreground">
          Teste o que aprendeu
        </h2>
        <p className="text-sm text-muted-foreground">Autoavaliação: não vale nota e não altera seu progresso.</p>
      </div>
      <div className="space-y-6">
        {questoes.map((q, i) => (
          <Questao key={q.id} questao={q} indice={i + 1} />
        ))}
      </div>
    </aside>
  );
}

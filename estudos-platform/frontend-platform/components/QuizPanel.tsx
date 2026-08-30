"use client";

import { useState } from "react";
import type { QuizQuestao } from "@/lib/types";

type Props = { questoes: QuizQuestao[] };

export function QuizPanel({ questoes }: Props) {
  const [respostas, setRespostas] = useState<Record<string, string>>({});

  return (
    <aside className="card space-y-4">
      <div>
        <h3 className="font-semibold">Checklist de estudo</h3>
        <p className="mt-1 text-sm text-zinc-400">
          Não é prova — marque o que você já fez. Correção no servidor vem depois.
        </p>
      </div>
      <ul className="space-y-4">
        {questoes.map((q) => (
          <li key={q.id} className="space-y-2">
            <p className="text-sm font-medium">{q.enunciado}</p>
            <div className="flex flex-wrap gap-2">
              {q.opcoes.map((op) => {
                const selected = respostas[q.id] === op.id;
                return (
                  <button
                    key={op.id}
                    type="button"
                    onClick={() => setRespostas((prev) => ({ ...prev, [q.id]: op.id }))}
                    className={`rounded-lg border px-3 py-1.5 text-xs transition ${
                      selected
                        ? "border-accent bg-accent-muted text-accent"
                        : "border-surface-border text-zinc-400 hover:border-zinc-500"
                    }`}
                  >
                    {op.texto}
                  </button>
                );
              })}
            </div>
          </li>
        ))}
      </ul>
    </aside>
  );
}

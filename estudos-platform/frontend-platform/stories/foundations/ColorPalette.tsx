"use client";

import * as React from "react";

type Swatch = {
  /** Nome do token CSS sem `--` (ex.: "relp-done"). */
  token: string;
  /** Classe Tailwind de fundo que usa esse token. */
  bg: string;
  nota?: string;
};

type Grupo = { titulo: string; descricao: string; itens: Swatch[] };

/** Fonte única da paleta Relp: tudo vem de `app/globals.css`. Se um token mudar lá, muda aqui. */
export const GRUPOS: Grupo[] = [
  {
    titulo: "Marca e superfícies",
    descricao: "Azul Relp para ações e links. Elevação = borda + superfície, nunca sombra solta.",
    itens: [
      { token: "primary", bg: "bg-primary", nota: "Ação principal, links, foco" },
      { token: "primary-hover", bg: "bg-primary-hover" },
      { token: "primary-muted", bg: "bg-primary-muted", nota: "Fundo de destaque suave" },
      { token: "background", bg: "bg-background border border-border" },
      { token: "card", bg: "bg-card border border-border" },
      { token: "muted", bg: "bg-muted border border-border" },
      { token: "border", bg: "bg-border" },
    ],
  },
  {
    titulo: "Educação",
    descricao: "Verde = concluído / progresso. Verde-água = aprendizado (dica, quiz). Nunca use o azul para progresso.",
    itens: [
      { token: "relp-done", bg: "bg-done", nota: "Lido, barra de progresso" },
      { token: "relp-done-muted", bg: "bg-done-muted" },
      { token: "relp-learn", bg: "bg-learn", nota: "Dica / aprendizado" },
      { token: "relp-learn-muted", bg: "bg-learn-muted" },
      { token: "relp-progress-track", bg: "bg-progress-track" },
    ],
  },
  {
    titulo: "Status de conteúdo",
    descricao: "Badge de ciclo de vida do artigo.",
    itens: [
      { token: "relp-status-draft", bg: "bg-status-draft", nota: "Rascunho" },
      { token: "relp-status-review", bg: "bg-status-review", nota: "Em revisão" },
      { token: "relp-status-published", bg: "bg-status-published", nota: "Publicado" },
      { token: "relp-status-archived", bg: "bg-status-archived", nota: "Arquivado" },
      { token: "relp-status-locked", bg: "bg-status-locked", nota: "Bloqueado / futuro" },
    ],
  },
  {
    titulo: "Feedback",
    descricao: "Mensagens do sistema. Notificação (roxo) não é erro.",
    itens: [
      { token: "relp-danger", bg: "bg-danger", nota: "Erro" },
      { token: "relp-danger-muted", bg: "bg-danger-muted" },
      { token: "relp-warning", bg: "bg-warning", nota: "Atenção" },
      { token: "relp-warning-muted", bg: "bg-warning-muted" },
      { token: "relp-info", bg: "bg-info", nota: "Informação" },
      { token: "relp-info-muted", bg: "bg-info-muted" },
      { token: "relp-notify", bg: "bg-notify", nota: "Notificação" },
    ],
  },
];

/** Pares texto/fundo usados nos componentes (contraste AA verificado por addon-a11y). */
const PARES = [
  { rotulo: "Texto principal / fundo", classe: "bg-background text-foreground" },
  { rotulo: "Texto secundário / card", classe: "bg-card text-muted-foreground" },
  { rotulo: "Primário / muted", classe: "bg-primary-muted text-primary" },
  { rotulo: "Concluído", classe: "bg-done-muted text-done-fg" },
  { rotulo: "Erro", classe: "bg-danger-muted text-danger-fg" },
  { rotulo: "Atenção", classe: "bg-warning-muted text-warning-fg" },
  { rotulo: "Info", classe: "bg-info-muted text-info-fg" },
  { rotulo: "Dica / aprendizado", classe: "bg-learn-muted text-learn" },
];

function useValorDoToken(token: string): string {
  const [valor, setValor] = React.useState("");
  React.useEffect(() => {
    const ler = () => setValor(getComputedStyle(document.documentElement).getPropertyValue(`--${token}`).trim());
    ler();
    // O addon de tema troca a classe `dark` no <html>: relê quando muda.
    const obs = new MutationObserver(ler);
    obs.observe(document.documentElement, { attributes: true, attributeFilter: ["class"] });
    return () => obs.disconnect();
  }, [token]);
  return valor;
}

function Amostra({ token, bg, nota }: Swatch) {
  const valor = useValorDoToken(token);
  return (
    <li className="space-y-2">
      <div className={`h-16 rounded-lg ${bg}`} aria-hidden="true" />
      <div>
        <p className="font-mono text-sm font-medium text-foreground">--{token}</p>
        <p className="font-mono text-xs text-muted-foreground">{valor ? `hsl(${valor})` : "—"}</p>
        {nota && <p className="text-xs text-muted-foreground">{nota}</p>}
      </div>
    </li>
  );
}

export function ColorPalette() {
  return (
    <div className="max-w-4xl space-y-10">
      {GRUPOS.map((g) => (
        <section key={g.titulo} aria-labelledby={`g-${g.titulo}`} className="space-y-3">
          <div>
            <h2 id={`g-${g.titulo}`} className="text-lg font-semibold text-foreground">
              {g.titulo}
            </h2>
            <p className="text-sm text-muted-foreground">{g.descricao}</p>
          </div>
          <ul className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
            {g.itens.map((s) => (
              <Amostra key={s.token} {...s} />
            ))}
          </ul>
        </section>
      ))}

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-foreground">Pares de texto</h2>
        <ul className="grid gap-3 sm:grid-cols-2">
          {PARES.map((p) => (
            <li key={p.rotulo} className={`rounded-lg border border-border px-4 py-3 text-sm font-medium ${p.classe}`}>
              {p.rotulo}
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}

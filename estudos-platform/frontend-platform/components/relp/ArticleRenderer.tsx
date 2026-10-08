import * as React from "react";

import { CalloutBlock } from "@/components/relp/CalloutBlock";
import { CodeBlock } from "@/components/relp/CodeBlock";
import type { Block } from "@/lib/types";
import { cn } from "@/lib/utils";

/** Só http(s): qualquer outro esquema (javascript:, data:) vira texto simples. */
function urlSegura(url: string): string | null {
  try {
    const u = new URL(url);
    return u.protocol === "http:" || u.protocol === "https:" ? u.toString() : null;
  } catch {
    return null;
  }
}

const LINK_CLASSE =
  "break-words text-primary underline underline-offset-2 hover:text-primary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded-sm";

function Link({ url, children }: { url: string; children?: React.ReactNode }) {
  const href = urlSegura(url);
  if (!href) return <span>{children ?? url}</span>;
  return (
    <a href={href} target="_blank" rel="noopener noreferrer" className={LINK_CLASSE}>
      {children ?? url}
      <span className="sr-only"> (abre em nova aba)</span>
    </a>
  );
}

/** `código` e **negrito** dentro de parágrafos. React escapa o texto: nada de HTML vindo da API. */
function Inline({ text }: { text: string }) {
  const partes = text.split(/(`[^`]+`|\*\*[^*]+\*\*)/g).filter((p) => p !== "");
  return (
    <>
      {partes.map((p, i) => {
        if (p.length > 2 && p.startsWith("`") && p.endsWith("`")) {
          return (
            <code key={i} className="rounded bg-muted px-1.5 py-0.5 font-mono text-[0.875em] text-foreground">
              {p.slice(1, -1)}
            </code>
          );
        }
        if (p.length > 4 && p.startsWith("**") && p.endsWith("**")) {
          return (
            <strong key={i} className="font-semibold text-foreground">
              {p.slice(2, -2)}
            </strong>
          );
        }
        return <React.Fragment key={i}>{p}</React.Fragment>;
      })}
    </>
  );
}

type Item =
  | { kind: "block"; block: Block }
  | { kind: "bullets"; itens: string[] };

/** Parágrafos consecutivos iniciados por "- " viram uma lista de verdade (acessível). */
function agrupar(blocks: Block[]): Item[] {
  const out: Item[] = [];
  const BULLET = /^\s*[-•]\s+/;

  const pushBullet = (texto: string) => {
    const prev = out[out.length - 1];
    if (prev?.kind === "bullets") prev.itens.push(texto);
    else out.push({ kind: "bullets", itens: [texto] });
  };

  for (const block of blocks) {
    if (block.type !== "p") {
      out.push({ kind: "block", block });
      continue;
    }
    // Um parágrafo da API pode trazer várias linhas ("texto\n- item\n- item"): cada linha é tratada à parte.
    for (const linha of block.text.split(/\r?\n/)) {
      if (!linha.trim()) continue;
      // O content-job às vezes junta os itens numa linha só ("- a. - b. - c."): linha que COMEÇA com "- "
      // é lista, e cada " - " interno separa um item.
      if (BULLET.test(linha)) {
        linha
          .replace(BULLET, "")
          .split(/\s+[-•]\s+/)
          .filter(Boolean)
          .forEach(pushBullet);
      }
      else out.push({ kind: "block", block: { type: "p", text: linha } });    }
  }
  return out;
}

export type ArticleRendererProps = {
  blocks: Block[];
  className?: string;
};

/** Renderiza `conteudo.blocks` da API. Tipos desconhecidos são ignorados (não quebram a página). */
export function ArticleRenderer({ blocks, className }: ArticleRendererProps) {
  const itens = agrupar(blocks ?? []);

  return (
    <div className={cn("space-y-4 text-base leading-relaxed", className)}>
      {itens.map((item, i) => {
        if (item.kind === "bullets") {
          return (
            <ul key={i} className="list-disc space-y-1.5 pl-6 text-muted-foreground marker:text-primary">
              {item.itens.map((t, j) => (
                <li key={j}>
                  <Inline text={t} />
                </li>
              ))}
            </ul>
          );
        }

        const b = item.block;
        switch (b.type) {
          case "h":
            return b.level <= 2 ? (
              <h2 key={i} className="pt-4 text-2xl font-semibold tracking-tight text-foreground">
                <Inline text={b.text} />
              </h2>
            ) : (
              <h3 key={i} className="pt-2 text-lg font-semibold text-foreground">
                <Inline text={b.text} />
              </h3>
            );

          case "p": {
            // Parágrafo que é só uma URL vira link clicável.
            if (/^https?:\/\/\S+$/.test(b.text.trim())) {
              return (
                <p key={i}>
                  <Link url={b.text.trim()} />
                </p>
              );
            }
            return (
              <p key={i} className="text-muted-foreground">
                <Inline text={b.text} />
              </p>
            );
          }

          case "code":
            return <CodeBlock key={i} code={b.text} lang={b.lang} />;

          case "list": {
            const Tag = b.ordered ? "ol" : "ul";
            return (
              <Tag
                key={i}
                className={cn(
                  "space-y-1.5 pl-6 text-muted-foreground marker:text-primary",
                  b.ordered ? "list-decimal" : "list-disc",
                )}
              >
                {(b.items ?? []).map((t, j) => (
                  <li key={j}>
                    <Inline text={t} />
                  </li>
                ))}
              </Tag>
            );
          }

          case "link":
            return (
              <p key={i}>
                <Link url={b.url}>{b.text}</Link>
              </p>
            );

          case "callout":
            return (
              <CalloutBlock key={i} variant={b.variant} title={b.title}>
                <Inline text={b.text} />
              </CalloutBlock>
            );

          default:
            return null;
        }
      })}
    </div>
  );
}

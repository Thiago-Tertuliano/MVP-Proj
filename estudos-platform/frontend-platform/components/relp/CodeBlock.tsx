"use client";

import * as React from "react";
import { Check, Copy } from "lucide-react";

import { Button } from "@/components/ui/button";
import { tokenizar, type TokenType } from "@/lib/highlight";
import { cn } from "@/lib/utils";

const COR: Record<TokenType, string> = {
  plain: "",
  keyword: "font-semibold text-primary",
  string: "text-learn",
  comment: "italic text-muted-foreground",
  number: "text-warning-fg",
};

async function copiar(texto: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(texto);
    return true;
  } catch {
    // Fallback para contextos sem Clipboard API (http, iframes antigos).
    try {
      const ta = document.createElement("textarea");
      ta.value = texto;
      ta.setAttribute("readonly", "");
      ta.style.position = "fixed";
      ta.style.opacity = "0";
      document.body.appendChild(ta);
      ta.select();
      const ok = document.execCommand("copy");
      document.body.removeChild(ta);
      return ok;
    } catch {
      return false;
    }
  }
}

export type CodeBlockProps = {
  code: string;
  lang?: string;
  className?: string;
};

/** Bloco de código: realce leve, rótulo da linguagem e botão "Copiar" acessível. */
export function CodeBlock({ code, lang, className }: CodeBlockProps) {
  const [copiado, setCopiado] = React.useState(false);
  const timer = React.useRef<ReturnType<typeof setTimeout> | null>(null);
  const tokens = React.useMemo(() => tokenizar(code, lang), [code, lang]);

  React.useEffect(
    () => () => {
      if (timer.current) clearTimeout(timer.current);
    },
    [],
  );

  async function aoCopiar() {
    if (await copiar(code)) {
      setCopiado(true);
      if (timer.current) clearTimeout(timer.current);
      timer.current = setTimeout(() => setCopiado(false), 2000);
    }
  }

  return (
    <figure className={cn("overflow-hidden rounded-lg border border-border bg-muted", className)}>
      <figcaption className="flex items-center justify-between border-b border-border px-3 py-1.5">
        <span className="font-mono text-xs uppercase tracking-wide text-muted-foreground">{lang || "código"}</span>
        <Button type="button" variant="ghost" size="sm" className="h-7 px-2 text-xs" onClick={aoCopiar}>
          {copiado ? <Check className="text-done-fg" /> : <Copy />}
          <span aria-live="polite">{copiado ? "Copiado!" : "Copiar"}</span>
          <span className="sr-only"> código</span>
        </Button>
      </figcaption>
      {/* tabIndex: permite rolar o código longo só com o teclado. */}
      <pre tabIndex={0} className="overflow-x-auto p-4 font-mono text-sm leading-relaxed text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring">
        <code>
          {tokens.map((t, i) =>
            t.type === "plain" ? (
              <React.Fragment key={i}>{t.text}</React.Fragment>
            ) : (
              <span key={i} className={COR[t.type]}>
                {t.text}
              </span>
            ),
          )}
        </code>
      </pre>
    </figure>
  );
}

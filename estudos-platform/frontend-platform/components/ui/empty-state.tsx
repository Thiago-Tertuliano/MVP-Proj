import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { Inbox } from "lucide-react";

import { cn } from "@/lib/utils";

const emptyStateVariants = cva(
  "flex flex-col items-center justify-center gap-3 rounded-xl border border-dashed px-6 py-12 text-center",
  {
    variants: {
      tone: {
        neutral: "border-border bg-card",
        danger: "border-danger/40 bg-danger-muted",
      },
    },
    defaultVariants: { tone: "neutral" },
  },
);

export interface EmptyStateProps
  extends Omit<React.HTMLAttributes<HTMLDivElement>, "title">,
    VariantProps<typeof emptyStateVariants> {
  title: string;
  description?: React.ReactNode;
  /** Ícone lucide; padrão `Inbox`. */
  icon?: React.ReactNode;
  /** Botão/link de ação (ex.: "Tentar novamente"). */
  action?: React.ReactNode;
  /** Nível de heading do título (a11y). Padrão h2. */
  headingLevel?: "h1" | "h2" | "h3";
}

/** Estado vazio genérico: "Nenhum artigo encontrado", "Sem trilhas", erro de API etc. */
function EmptyState({
  title,
  description,
  icon,
  action,
  tone,
  className,
  headingLevel: Heading = "h2",
  ...props
}: EmptyStateProps) {
  return (
    <div className={cn(emptyStateVariants({ tone }), className)} {...props}>
      <div
        className={cn(
          "flex h-12 w-12 items-center justify-center rounded-full [&_svg]:h-6 [&_svg]:w-6",
          tone === "danger" ? "bg-danger/15 text-danger-fg" : "bg-muted text-muted-foreground",
        )}
        aria-hidden="true"
      >
        {icon ?? <Inbox />}
      </div>
      <div className="space-y-1">
        <Heading className="text-base font-semibold text-foreground">{title}</Heading>
        {description && <p className="mx-auto max-w-md text-sm text-muted-foreground">{description}</p>}
      </div>
      {action && <div className="pt-1">{action}</div>}
    </div>
  );
}

export { EmptyState };

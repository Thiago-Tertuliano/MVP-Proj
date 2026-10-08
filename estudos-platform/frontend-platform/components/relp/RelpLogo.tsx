import * as React from "react";

import { cn } from "@/lib/utils";

type RelpLogoProps = React.HTMLAttributes<HTMLSpanElement> & {
  /** `mark` = só o ícone; `full` = ícone + wordmark. */
  variant?: "full" | "mark";
  /** Tamanho do ícone em px. */
  size?: number;
};

/** Ícone Relp: três nós conectados (um roadmap) sobre o azul da marca. Só tokens, funciona em light/dark. */
function RelpMark({ size = 32, className }: { size?: number; className?: string }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 32 32"
      role="img"
      aria-label="Relp"
      className={cn("shrink-0", className)}
    >
      <rect width="32" height="32" rx="7" className="fill-primary" />
      <path
        d="M9 21.5 15 11l8 3.5"
        className="stroke-primary-foreground/70"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      />
      <circle cx="9" cy="21.5" r="3" className="fill-primary-foreground" />
      <circle cx="15" cy="11" r="3" className="fill-primary-foreground" />
      <circle cx="23" cy="14.5" r="3" className="fill-done" />
    </svg>
  );
}

function RelpLogo({ variant = "full", size = 32, className, ...props }: RelpLogoProps) {
  if (variant === "mark") {
    return (
      <span className={cn("inline-flex", className)} {...props}>
        <RelpMark size={size} />
      </span>
    );
  }
  return (
    <span className={cn("inline-flex items-center gap-2", className)} {...props}>
      <RelpMark size={size} />
      <span className="text-lg font-bold tracking-tight text-foreground" aria-hidden="true">
        Relp<span className="text-done">!</span>
      </span>
    </span>
  );
}

export { RelpLogo, RelpMark };

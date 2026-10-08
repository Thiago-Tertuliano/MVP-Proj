import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default: "border-transparent bg-primary-muted text-primary",
        secondary: "border-transparent bg-secondary text-secondary-foreground",
        outline: "border-border text-foreground",
        /** Status de conteúdo */
        publicado: "border-transparent bg-done-muted text-done-fg",
        rascunho: "border-transparent bg-warning-muted text-warning-fg",
        revisao: "border-transparent bg-info-muted text-info-fg",
        arquivado: "border-transparent bg-muted text-muted-foreground",
        /** Educação */
        done: "border-transparent bg-done-muted text-done-fg",
        learn: "border-transparent bg-learn-muted text-learn",
        /** Feedback */
        danger: "border-transparent bg-danger-muted text-danger-fg",
        warning: "border-transparent bg-warning-muted text-warning-fg",
        info: "border-transparent bg-info-muted text-info-fg",
        /** Notificação (não é erro) */
        notify: "border-transparent bg-notify text-primary-foreground",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  },
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return <span className={cn(badgeVariants({ variant }), className)} {...props} />;
}

export { Badge, badgeVariants };

"use client";

import * as React from "react";
import * as ProgressPrimitive from "@radix-ui/react-progress";

import { cn } from "@/lib/utils";

type ProgressProps = React.ComponentPropsWithoutRef<typeof ProgressPrimitive.Root> & {
  /** Classe do preenchimento. Padrão: verde educação (`bg-done`), nunca o azul da marca. */
  indicatorClassName?: string;
};

const Progress = React.forwardRef<React.ElementRef<typeof ProgressPrimitive.Root>, ProgressProps>(
  ({ className, value, indicatorClassName, ...props }, ref) => {
    const pct = Math.min(100, Math.max(0, value ?? 0));
    return (
      <ProgressPrimitive.Root
        ref={ref}
        value={pct}
        className={cn("relative h-2 w-full overflow-hidden rounded-full bg-progress-track", className)}
        {...props}
      >
        <ProgressPrimitive.Indicator
          className={cn("h-full w-full flex-1 rounded-full bg-done transition-all", indicatorClassName)}
          style={{ transform: `translateX(-${100 - pct}%)` }}
        />
      </ProgressPrimitive.Root>
    );
  },
);
Progress.displayName = ProgressPrimitive.Root.displayName;

export { Progress };

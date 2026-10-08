"use client";

import { Toaster as Sonner, toast } from "sonner";

type ToasterProps = React.ComponentProps<typeof Sonner>;

/**
 * Toaster Relp: sonner estilizado só com tokens (sem hex).
 * success = verde educação, error = danger, warning = aviso, info = info.
 */
function Toaster({ ...props }: ToasterProps) {
  return (
    <Sonner
      className="toaster group"
      position="bottom-right"
      closeButton
      toastOptions={{
        classNames: {
          toast:
            "group toast !rounded-lg !border !border-border !bg-card !text-card-foreground !shadow-none",
          title: "!font-semibold",
          description: "!text-muted-foreground",
          actionButton: "!bg-primary !text-primary-foreground",
          cancelButton: "!bg-muted !text-muted-foreground",
          closeButton: "!border-border !bg-card !text-foreground",
          success: "!border-done/50 [&_[data-icon]]:!text-done",
          error: "!border-danger/50 [&_[data-icon]]:!text-danger",
          warning: "!border-warning/50 [&_[data-icon]]:!text-warning",
          info: "!border-info/50 [&_[data-icon]]:!text-info",
        },
      }}
      {...props}
    />
  );
}

export { Toaster, toast };

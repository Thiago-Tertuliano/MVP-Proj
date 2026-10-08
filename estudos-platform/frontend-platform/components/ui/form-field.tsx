"use client";

import * as React from "react";

import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

/** Texto de apoio abaixo do campo. */
function HelperText({ className, ...props }: React.HTMLAttributes<HTMLParagraphElement>) {
  return <p className={cn("text-xs text-muted-foreground", className)} {...props} />;
}

/** Mensagem de erro do campo — anunciada por leitores de tela. */
function FieldError({ className, children, ...props }: React.HTMLAttributes<HTMLParagraphElement>) {
  if (!children) return null;
  return (
    <p role="alert" className={cn("text-xs font-medium text-danger-fg", className)} {...props}>
      {children}
    </p>
  );
}

export type FormFieldControlProps = {
  id: string;
  "aria-describedby"?: string;
  "aria-invalid"?: boolean;
  "aria-required"?: boolean;
};

type FormFieldProps = {
  label: string;
  helper?: React.ReactNode;
  error?: React.ReactNode;
  required?: boolean;
  className?: string;
  /** Recebe os atributos de acessibilidade (id, aria-*) para repassar ao controle. */
  children: (control: FormFieldControlProps) => React.ReactNode;
};

/**
 * Compõe Label + controle + HelperText + FieldError, ligando `id` e `aria-describedby`
 * automaticamente. Funciona com Input, Textarea, Select (trigger), etc.
 */
function FormField({ label, helper, error, required, className, children }: FormFieldProps) {
  const id = React.useId();
  const helperId = `${id}-helper`;
  const errorId = `${id}-error`;
  const describedBy = [error ? errorId : null, helper ? helperId : null].filter(Boolean).join(" ");

  return (
    <div className={cn("space-y-1.5", className)}>
      <Label htmlFor={id}>
        {label}
        {required && (
          <span aria-hidden="true" className="ml-0.5 text-danger-fg">
            *
          </span>
        )}
      </Label>
      {children({
        id,
        "aria-describedby": describedBy || undefined,
        "aria-invalid": error ? true : undefined,
        "aria-required": required || undefined,
      })}
      {helper && !error && <HelperText id={helperId}>{helper}</HelperText>}
      <FieldError id={errorId}>{error}</FieldError>
    </div>
  );
}

export { FormField, HelperText, FieldError };

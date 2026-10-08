"use client";

import * as React from "react";
import Link from "next/link";
import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff } from "lucide-react";
import { useForm } from "react-hook-form";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { FormField } from "@/components/ui/form-field";
import { Input } from "@/components/ui/input";
import type { FormAlertData } from "@/lib/auth-errors";
import { loginSchema, registroSchema } from "@/lib/validation";

export type AuthValues = { nome?: string; email: string; senha: string };

export type AuthFormProps = {
  mode: "login" | "registro";
  /** Chamado só com dados válidos. Quem usa trata a rede e devolve o erro em `serverError`. */
  onSubmit: (values: AuthValues) => void | Promise<void>;
  serverError?: FormAlertData | null;
  /** Estado de envio controlado de fora (opcional; além do interno do react-hook-form). */
  submitting?: boolean;
  /** Link "Criar conta" / "Já tenho conta". */
  alternateHref?: string;
  defaultValues?: Partial<AuthValues>;
};

const TEXTOS = {
  login: {
    submit: "Entrar",
    alternate: { pergunta: "Ainda não tem conta?", link: "Criar conta" },
    senhaHelper: undefined,
    autocompleteSenha: "current-password",
  },
  registro: {
    submit: "Criar conta",
    alternate: { pergunta: "Já tem conta?", link: "Entrar" },
    senhaHelper: "Mínimo de 8 caracteres.",
    autocompleteSenha: "new-password",
  },
} as const;

export function AuthForm({
  mode,
  onSubmit,
  serverError,
  submitting,
  alternateHref,
  defaultValues,
}: AuthFormProps) {
  const [verSenha, setVerSenha] = React.useState(false);
  const t = TEXTOS[mode];

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<AuthValues>({
    resolver: zodResolver(mode === "login" ? loginSchema : registroSchema) as never,
    defaultValues: { nome: "", email: "", senha: "", ...defaultValues },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const enviando = isSubmitting || !!submitting;

  return (
    <form
      onSubmit={handleSubmit(async (values) => {
        await onSubmit(values);
      })}
      noValidate
      className="space-y-5 rounded-xl border border-border bg-card p-6"
    >
      {serverError && (
        <Alert variant={serverError.tone}>
          <AlertTitle>{serverError.titulo}</AlertTitle>
          <AlertDescription>{serverError.mensagem}</AlertDescription>
        </Alert>
      )}

      {mode === "registro" && (
        <FormField label="Nome" required error={errors.nome?.message}>
          {(control) => (
            <Input type="text" autoComplete="name" {...control} {...register("nome")} />
          )}
        </FormField>
      )}

      <FormField label="E-mail" required error={errors.email?.message}>
        {(control) => (
          <Input
            type="email"
            inputMode="email"
            autoComplete="email"
            placeholder="nome@exemplo.com"
            {...control}
            {...register("email")}
          />
        )}
      </FormField>

      <FormField label="Senha" required helper={t.senhaHelper} error={errors.senha?.message}>
        {(control) => (
          <div className="relative">
            <Input
              type={verSenha ? "text" : "password"}
              autoComplete={t.autocompleteSenha}
                className="pr-11"
              {...control}
              {...register("senha")}
            />
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="absolute right-0 top-0 h-10 w-10 text-muted-foreground"
              onClick={() => setVerSenha((v) => !v)}
              aria-label={verSenha ? "Ocultar senha" : "Mostrar senha"}
              aria-pressed={verSenha}
              >
              {verSenha ? <EyeOff /> : <Eye />}
            </Button>
          </div>
        )}
      </FormField>

      <Button type="submit" className="w-full" loading={enviando}>
        {enviando ? (mode === "login" ? "Entrando…" : "Criando conta…") : t.submit}
      </Button>

      {alternateHref && (
        <p className="text-center text-sm text-muted-foreground">
          {t.alternate.pergunta}{" "}
          <Link
            href={alternateHref}
            className="rounded-sm font-medium text-primary underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            {t.alternate.link}
          </Link>
        </p>
      )}
    </form>
  );
}

import type { Meta, StoryObj } from "@storybook/react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";
import { Progress } from "@/components/ui/progress";
import { Skeleton, Spinner } from "@/components/ui/skeleton";
import { toast, Toaster } from "@/components/ui/sonner";

/** Feedback ao usuário: Badge, Alert, Toast (sonner), Progress, Skeleton/Spinner e EmptyState. */
const meta = {
  title: "UI/Feedback",
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

export const Badges: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <Badge>Padrão</Badge>
      <Badge variant="secondary">Secundário</Badge>
      <Badge variant="outline">Outline</Badge>
      <Badge variant="publicado">Publicado</Badge>
      <Badge variant="rascunho">Rascunho</Badge>
      <Badge variant="revisao">Em revisão</Badge>
      <Badge variant="arquivado">Arquivado</Badge>
      <Badge variant="done">Lido</Badge>
      <Badge variant="learn">Aprendizado</Badge>
      <Badge variant="danger">Erro</Badge>
      <Badge variant="warning">Atenção</Badge>
      <Badge variant="info">Info</Badge>
      <Badge variant="notify">3</Badge>
    </div>
  ),
};

export const Alertas: Story = {
  render: () => (
    <div className="max-w-xl space-y-3">
      <Alert variant="info">
        <AlertTitle>Informação</AlertTitle>
        <AlertDescription>Novos artigos entram toda semana.</AlertDescription>
      </Alert>
      <Alert variant="success">
        <AlertTitle>Conta criada</AlertTitle>
        <AlertDescription>Bem-vindo(a) ao Relp!</AlertDescription>
      </Alert>
      <Alert variant="warning">
        <AlertTitle>Muitas tentativas</AlertTitle>
        <AlertDescription>Por segurança, aguarde 1 minuto.</AlertDescription>
      </Alert>
      <Alert variant="danger">
        <AlertTitle>Não foi possível entrar</AlertTitle>
        <AlertDescription>E-mail ou senha incorretos.</AlertDescription>
      </Alert>
    </div>
  ),
};

export const AlertasEscuro: Story = { ...Alertas, globals: { theme: "dark" } };

export const Toasts: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <Toaster />
      <Button variant="outline" onClick={() => toast.success("Artigo marcado como lido.")}>
        Sucesso
      </Button>
      <Button variant="outline" onClick={() => toast.error("Não foi possível salvar seu progresso.")}>
        Erro
      </Button>
      <Button variant="outline" onClick={() => toast.warning("Sua sessão expirou.")}>
        Aviso
      </Button>
      <Button variant="outline" onClick={() => toast("Anotação salva", { description: "Há poucos segundos" })}>
        Neutro
      </Button>
    </div>
  ),
};

export const Progresso: Story = {
  render: () => (
    <div className="max-w-sm space-y-4">
      {[0, 33, 66, 100].map((v) => (
        <div key={v} className="space-y-1">
          <p className="text-xs text-muted-foreground">{v}%</p>
          <Progress value={v} aria-label={`Progresso ${v}%`} />
        </div>
      ))}
    </div>
  ),
};

export const Carregando: Story = {
  render: () => (
    <div className="max-w-sm space-y-4">
      <Skeleton className="h-6 w-48" />
      <Skeleton className="h-20 w-full" />
      <Spinner />
    </div>
  ),
};

export const EstadoVazio: Story = {
  render: () => (
    <EmptyState
      title="Nenhum artigo encontrado"
      description="Tente outras palavras ou explore as trilhas."
      action={<Button variant="outline">Ver trilhas</Button>}
    />
  ),
};

export const EstadoDeErro: Story = {
  render: () => (
    <EmptyState
      tone="danger"
      title="Não conseguimos carregar esta página"
      description="O servidor pode estar indisponível. Tente novamente."
      action={<Button>Tentar de novo</Button>}
    />
  ),
};

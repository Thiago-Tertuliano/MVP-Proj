import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";

import { ArticleRenderer } from "@/components/relp/ArticleRenderer";
import { CalloutBlock } from "@/components/relp/CalloutBlock";
import { CodeBlock } from "@/components/relp/CodeBlock";
import { GuestPrompt } from "@/components/relp/GuestPrompt";
import { NoteEditor } from "@/components/relp/NoteEditor";
import { QuizPanel } from "@/components/relp/QuizPanel";
import { blocosExemplo, questoesExemplo } from "../../stories/fixtures";

const meta = {
  title: "Relp/Artigo",
  parameters: { layout: "padded" },
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

export const Conteudo: Story = {
  name: "ArticleRenderer",
  render: () => (
    <div className="max-w-2xl">
      <ArticleRenderer blocks={blocosExemplo} />
    </div>
  ),
};

export const ConteudoEscuro: Story = { ...Conteudo, name: "ArticleRenderer / escuro", globals: { theme: "dark" } };

export const Codigo: Story = {
  name: "CodeBlock",
  render: () => (
    <div className="max-w-2xl space-y-4">
      <CodeBlock lang="go" code={'package main\n\n// Olá\nfunc main() {\n\tprintln("oi", 1)\n}'} />
      <CodeBlock lang="sql" code={"SELECT id, titulo FROM artigos WHERE status = 'publicado' LIMIT 10;"} />
      <CodeBlock code="linguagem desconhecida: texto simples" />
    </div>
  ),
};

export const Destaques: Story = {
  name: "CalloutBlock",
  render: () => (
    <div className="max-w-2xl space-y-3">
      <CalloutBlock variant="info">Pacotes de teste terminam em _test.</CalloutBlock>
      <CalloutBlock variant="warning">Nomes minúsculos não são exportados.</CalloutBlock>
      <CalloutBlock variant="tip">Rode go vet antes do PR.</CalloutBlock>
    </div>
  ),
};

export const Anotacoes: Story = {
  name: "NoteEditor / estados",
  render: () => {
    const estados = ["idle", "saving", "saved", "error"] as const;
    return (
      <div className="grid max-w-3xl gap-6 md:grid-cols-2">
        {estados.map((status) => (
          <NoteEditor key={status} value="Resumo da aula…" onChange={() => undefined} status={status} onRetry={() => undefined} />
        ))}
        <NoteEditor value="" onChange={() => undefined} loading />
        <NoteEditor value="" onChange={() => undefined} loadError />
      </div>
    );
  },
};

function QuizInterativo() {
  return (
    <div className="max-w-md">
      <QuizPanel questoes={questoesExemplo} />
    </div>
  );
}

export const Quiz: Story = {
  name: "QuizPanel",
  render: () => <QuizInterativo />,
};

/** Fluxo completo: escolher errado → ver correção → tentar de novo → acertar. */
export const QuizInteracao: Story = {
  name: "QuizPanel / interação",
  render: () => <QuizInterativo />,
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const verificar = canvas.getByRole("button", { name: "Verificar" });
    await expect(verificar).toBeDisabled();

    await userEvent.click(canvas.getByRole("radio", { name: /Hello world/ }));
    await userEvent.click(verificar);
    await expect(await canvas.findByText("Ainda não foi dessa vez")).toBeInTheDocument();

    await userEvent.click(canvas.getByRole("button", { name: "Tentar de novo" }));
    await userEvent.click(canvas.getByRole("radio", { name: /Neste repositório/ }));
    await userEvent.click(canvas.getByRole("button", { name: "Verificar" }));
    await expect(await canvas.findByText("Correto!")).toBeInTheDocument();
  },
};

export const ConviteVisitante: Story = {
  name: "GuestPrompt",
  render: () => (
    <div className="max-w-md">
      <GuestPrompt
        titulo="Salve seu progresso"
        descricao="Entre para marcar artigos como lidos e guardar anotações."
        next="/artigos/pacotes-em-go"
      />
    </div>
  ),
};

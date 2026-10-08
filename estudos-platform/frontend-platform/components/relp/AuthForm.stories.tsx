import type { Meta, StoryObj } from "@storybook/react";
import { expect, fn, userEvent, waitFor, within } from "@storybook/test";

import { AuthForm } from "@/components/relp/AuthForm";

const meta = {
  title: "Relp/AuthForm",
  component: AuthForm,
  parameters: { layout: "centered" },
  decorators: [
    (Story) => (
      <div className="w-[26rem] max-w-full">
        <Story />
      </div>
    ),
  ],
  args: { mode: "login", onSubmit: fn(), alternateHref: "/registro" },
} satisfies Meta<typeof AuthForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Login: Story = {};

export const Registro: Story = { args: { mode: "registro", alternateHref: "/login" } };

export const Enviando: Story = { args: { submitting: true } };

export const ErroCredenciais: Story = {
  args: { serverError: { tone: "danger", titulo: "Não foi possível entrar", mensagem: "E-mail ou senha incorretos." } },
};

export const RateLimit: Story = {
  args: {
    serverError: { tone: "warning", titulo: "Muitas tentativas", mensagem: "Por segurança, aguarde 1 minuto antes de tentar de novo." },
  },
};

export const EmailJaCadastrado: Story = {
  args: {
    mode: "registro",
    serverError: { tone: "danger", titulo: "E-mail já cadastrado", mensagem: "Já existe uma conta com esse e-mail." },
  },
};

export const Escuro: Story = { globals: { theme: "dark" } };

export const ValidacaoClient: Story = {
  name: "Validação (campos vazios)",
  play: async ({ canvasElement, args }) => {
    const canvas = within(canvasElement);
    await userEvent.click(canvas.getByRole("button", { name: "Entrar" }));
    await expect(await canvas.findByText("Informe seu e-mail.")).toBeInTheDocument();
    await expect(await canvas.findByText("Informe sua senha.")).toBeInTheDocument();
    await expect(args.onSubmit).not.toHaveBeenCalled();
  },
};

export const EnvioValido: Story = {
  name: "Envio válido + mostrar senha",
  play: async ({ canvasElement, args }) => {
    const canvas = within(canvasElement);
    await userEvent.type(canvas.getByLabelText(/E-mail/), "aluno@exemplo.com");
    await userEvent.type(canvas.getByLabelText(/Senha/), "senha1234");

    await userEvent.click(canvas.getByRole("button", { name: "Mostrar senha" }));
    await expect(canvas.getByLabelText(/Senha/)).toHaveAttribute("type", "text");

    await userEvent.click(canvas.getByRole("button", { name: "Entrar" }));
    await waitFor(() => expect(args.onSubmit).toHaveBeenCalledTimes(1));
    await expect(args.onSubmit).toHaveBeenCalledWith(
      expect.objectContaining({ email: "aluno@exemplo.com", senha: "senha1234" }),
    );
  },
};

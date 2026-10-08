import type { Meta, StoryObj } from "@storybook/react";

import { AppHeader } from "@/components/relp/AppHeader";

const meta = {
  title: "Relp/AppHeader",
  component: AppHeader,
  parameters: { layout: "fullscreen" },
  args: { onSearch: () => undefined, onLogout: () => undefined },
} satisfies Meta<typeof AppHeader>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Visitante: Story = { args: { state: "guest" } };
export const Logado: Story = { args: { state: "user", user: { nome: "Thiago Matos", email: "thiago@axellion.dev" } } };
export const VerificandoSessao: Story = { args: { state: "loading" } };
export const Saindo: Story = {
  args: { state: "user", user: { nome: "Thiago Matos", email: "thiago@axellion.dev" }, loggingOut: true },
};
export const Mobile: Story = {
  args: { state: "user", user: { nome: "Thiago Matos", email: "thiago@axellion.dev" } },
  parameters: { viewport: { defaultViewport: "mobile1" } },
};
export const Escuro: Story = { args: { state: "guest" }, globals: { theme: "dark" } };

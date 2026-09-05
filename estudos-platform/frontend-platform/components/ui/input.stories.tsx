import type { Meta, StoryObj } from "@storybook/react";

import { Input } from "@/components/ui/input";

const meta = {
  title: "UI/Input",
  component: Input,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  argTypes: {
    type: {
      control: "select",
      options: ["text", "email", "password", "search"],
    },
    disabled: { control: "boolean" },
  },
  args: {
    type: "text",
    placeholder: "Digite aqui…",
    disabled: false,
  },
  decorators: [
    (Story) => (
      <div className="w-80">
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Email: Story = {
  args: {
    type: "email",
    placeholder: "voce@relp.app",
    defaultValue: "autor.seed@estudos.local",
  },
};

export const Password: Story = {
  args: {
    type: "password",
    placeholder: "Senha",
    defaultValue: "senha1234",
  },
};

export const Disabled: Story = {
  args: {
    disabled: true,
    placeholder: "Desabilitado",
  },
};

export const Invalid: Story = {
  args: {
    "aria-invalid": true,
    defaultValue: "valor-invalido",
    className: "border-destructive focus-visible:ring-destructive",
  },
};

export const WithLabel: Story = {
  render: (args) => (
    <label className="flex w-80 flex-col gap-1.5 text-sm">
      <span className="font-medium text-foreground">E-mail</span>
      <Input {...args} type="email" placeholder="voce@relp.app" />
      <span className="text-xs text-muted-foreground">Usamos o e-mail só para login.</span>
    </label>
  ),
};

import type { Meta, StoryObj } from "@storybook/react";

import { ColorPalette } from "./ColorPalette";

const meta = {
  title: "Foundations/Paleta de cores",
  component: ColorPalette,
  parameters: { layout: "padded" },
} satisfies Meta<typeof ColorPalette>;

export default meta;
type Story = StoryObj<typeof meta>;

/** Tema claro (padrão). */
export const Claro: Story = {};

/** Tema escuro — os mesmos tokens com valores `.dark`. */
export const Escuro: Story = {
  globals: { theme: "dark" },
};

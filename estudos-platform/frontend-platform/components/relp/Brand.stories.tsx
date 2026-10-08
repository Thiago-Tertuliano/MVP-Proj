import type { Meta, StoryObj } from "@storybook/react";

import { RelpLogo } from "@/components/relp/RelpLogo";
import { StatusBadge } from "@/components/relp/StatusBadge";

/** Marca e selos pequenos: RelpLogo e StatusBadge. */
const meta = {
  title: "Relp/Marca e status",
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

export const Logo: Story = {
  render: () => (
    <div className="flex items-center gap-8">
      <RelpLogo />
      <RelpLogo variant="mark" size={48} />
    </div>
  ),
};

export const LogoEscuro: Story = { ...Logo, globals: { theme: "dark" } };

export const Status: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <StatusBadge status="rascunho" />
      <StatusBadge status="revisao" />
      <StatusBadge status="publicado" />
      <StatusBadge status="arquivado" />
    </div>
  ),
};

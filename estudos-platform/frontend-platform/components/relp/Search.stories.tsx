import type { Meta, StoryObj } from "@storybook/react";

import { SearchBar } from "@/components/relp/SearchBar";
import { SearchResults } from "@/components/relp/SearchResults";

const meta = {
  title: "Relp/Busca",
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

export const Barra: Story = {
  render: () => (
    <div className="max-w-md">
      <SearchBar onSubmit={() => undefined} />
    </div>
  ),
};

export const BarraHover: Story = { ...Barra, parameters: { pseudo: { hover: true } } };
export const BarraPreenchida: Story = {
  render: () => (
    <div className="max-w-md">
      <SearchBar defaultValue="interfaces" onSubmit={() => undefined} />
    </div>
  ),
};

export const Resultados: Story = {
  render: () => (
    <div className="max-w-xl">
      <SearchResults
        termo="go"
        itens={[
          { slug: "pacotes-em-go", titulo: "Pacotes em Go", similarity: 0.9 },
          { slug: "structs-e-metodos", titulo: "Structs e métodos", similarity: 0.6 },
        ]}
      />
    </div>
  ),
};

export const ResultadosEscuro: Story = { ...Resultados, globals: { theme: "dark" } };

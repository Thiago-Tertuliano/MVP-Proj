import type { Meta, StoryObj } from "@storybook/react";

import { ArticlePager } from "@/components/relp/ArticlePager";
import { ContinueCard } from "@/components/relp/ContinueCard";
import { ProgressSummary } from "@/components/relp/ProgressSummary";
import { RoadmapMap } from "@/components/relp/RoadmapMap";
import { TrilhaCard } from "@/components/relp/TrilhaCard";
import { roadmapLogado, roadmapVisitante } from "../../stories/fixtures";

const meta = {
  title: "Relp/Trilha",
  parameters: { layout: "padded" },
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

const base = {
  slug: "go-basico",
  titulo: "Go Básico",
  descricao: "Fundamentos de Go para ler o backend da plataforma.",
  totalModulos: 2,
};

export const CardVisitante: Story = {
  name: "TrilhaCard / visitante",
  render: () => (
    <div className="max-w-sm">
      <TrilhaCard {...base} />
    </div>
  ),
};

export const CardComProgresso: Story = {
  name: "TrilhaCard / com progresso",
  render: () => (
    <div className="max-w-sm">
      <TrilhaCard {...base} progresso={{ status: "ready", concluidos: 1, total: 3, percentual: 33.33 }} />
    </div>
  ),
};

export const CardHover: Story = {
  name: "TrilhaCard / hover",
  parameters: { pseudo: { hover: true } },
  render: CardComProgresso.render,
};

export const CardCarregandoProgresso: Story = {
  name: "TrilhaCard / carregando progresso",
  render: () => (
    <div className="max-w-sm">
      <TrilhaCard {...base} progresso={{ status: "loading" }} />
    </div>
  ),
};

export const CardEscuro: Story = { ...CardComProgresso, name: "TrilhaCard / escuro", globals: { theme: "dark" } };

export const Progresso: Story = {
  name: "ProgressSummary",
  render: () => (
    <div className="max-w-sm space-y-4">
      <ProgressSummary concluidos={0} total={3} />
      <ProgressSummary concluidos={1} total={3} />
      <ProgressSummary concluidos={3} total={3} />
    </div>
  ),
};

export const MapaLogado: Story = {
  name: "RoadmapMap / logado (lido, atual, futuro)",
  render: () => (
    <div className="max-w-2xl">
      <RoadmapMap modulos={roadmapLogado} progresso={{ concluidos: 1, total: 3 }} />
    </div>
  ),
};

export const MapaVisitante: Story = {
  name: "RoadmapMap / visitante",
  render: () => (
    <div className="max-w-2xl">
      <RoadmapMap modulos={roadmapVisitante} mostrarEstado={false} />
    </div>
  ),
};

export const MapaEscuro: Story = { ...MapaLogado, name: "RoadmapMap / escuro", globals: { theme: "dark" } };

export const Continuar: Story = {
  name: "ContinueCard",
  render: () => (
    <div className="max-w-2xl space-y-4">
      <ContinueCard artigoSlug="structs-e-metodos" artigoTitulo="Structs e métodos" trilhaTitulo="Go Básico" />
      <ContinueCard artigoSlug="pacotes-em-go" artigoTitulo="Pacotes em Go" trilhaTitulo="Go Básico" concluido />
    </div>
  ),
};

export const Paginacao: Story = {
  name: "ArticlePager",
  render: () => (
    <div className="max-w-2xl space-y-4">
      <ArticlePager
        anterior={{ slug: "a", titulo: "Pacotes em Go" }}
        proximo={{ slug: "c", titulo: "Interfaces implícitas" }}
      />
      <ArticlePager anterior={null} proximo={{ slug: "b", titulo: "Structs e métodos" }} />
      <ArticlePager anterior={{ slug: "b", titulo: "Structs e métodos" }} proximo={null} />
    </div>
  ),
};

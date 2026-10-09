import * as React from "react";
import {
  BookOpen,
  Brain,
  Code,
  Compass,
  Crown,
  Database,
  Flag,
  Flame,
  Footprints,
  Globe,
  Layers,
  Map,
  Medal,
  Network,
  Rocket,
  Server,
  Shield,
  Star,
  Swords,
  Terminal,
  Trophy,
  Zap,
  type LucideIcon,
} from "lucide-react";

/**
 * Ícones permitidos para roadmaps e conquistas. O servidor guarda só o nome (string livre do autor),
 * então tudo que não estiver aqui cai em `Map` — nunca quebra a tela.
 */
export const ICONES_ROADMAP = {
  map: Map,
  rocket: Rocket,
  code: Code,
  database: Database,
  globe: Globe,
  shield: Shield,
  server: Server,
  brain: Brain,
  terminal: Terminal,
  book: BookOpen,
  layers: Layers,
  zap: Zap,
  trophy: Trophy,
} satisfies Record<string, LucideIcon>;

const ICONES_CONQUISTA = {
  footprints: Footprints,
  compass: Compass,
  network: Network,
  swords: Swords,
  crown: Crown,
  flame: Flame,
  star: Star,
  map: Map,
  medal: Medal,
  flag: Flag,
  trophy: Trophy,
} satisfies Record<string, LucideIcon>;

const TODOS: Record<string, LucideIcon> = { ...ICONES_CONQUISTA, ...ICONES_ROADMAP };

export const NOMES_ICONES_ROADMAP = Object.keys(ICONES_ROADMAP);

export function iconePorNome(nome: string | undefined | null): LucideIcon {
  return (nome && TODOS[nome.toLowerCase()]) || Map;
}

/** Renderiza o ícone pelo nome (decorativo: o texto vizinho já diz o que é). */
export function Icone({ nome, className }: { nome?: string | null; className?: string }) {
  const Cmp = iconePorNome(nome);
  return <Cmp className={className} aria-hidden="true" />;
}

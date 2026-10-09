"use client";

import * as React from "react";
import {
  Background,
  Controls,
  Handle,
  MarkerType,
  Position,
  ReactFlow,
  type Edge,
  type Node,
  type NodeProps,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";

import { RoadmapNodeCard } from "@/components/roadmap/RoadmapNodeCard";
import { arestaPercorrida, autoLayout, LAYOUT, posicaoDoNo, semPosicoes } from "@/lib/roadmaps/graph";
import type { EstadoNo, RoadmapAresta, RoadmapNo } from "@/lib/roadmaps/types";
import { cn } from "@/lib/utils";

type FlowData = {
  no: RoadmapNo;
  estado: EstadoNo | null;
  selecionado: boolean;
  onSelecionar: (id: string) => void;
};
type FlowNode = Node<FlowData, "no">;

/** Handles invisíveis: o aluno não conecta nada, só precisamos deles para ancorar as arestas. */
const HANDLE_CLASS = "!pointer-events-none !h-1 !w-1 !border-0 !bg-transparent !opacity-0";

function NoDoRoadmap({ data }: NodeProps<FlowNode>) {
  return (
    <>
      <Handle type="target" position={Position.Top} isConnectable={false} className={HANDLE_CLASS} />
      <RoadmapNodeCard
        titulo={data.no.titulo}
        tipo={data.no.tipo}
        estado={data.estado}
        xp={data.no.xp}
        selecionado={data.selecionado}
        onSelecionar={() => data.onSelecionar(data.no.id)}
      />
      <Handle type="source" position={Position.Bottom} isConnectable={false} className={HANDLE_CLASS} />
    </>
  );
}

// Fora do componente: referência estável (o React Flow recria os nós se o objeto mudar a cada render).
const nodeTypes = { no: NoDoRoadmap };

export type RoadmapCanvasProps = {
  nos: RoadmapNo[];
  arestas: RoadmapAresta[];
  /** `null` para visitante: tudo neutro, nenhuma aresta "acesa". */
  estados: Record<string, EstadoNo> | null;
  selecionadoId?: string | null;
  onSelecionar: (id: string | null) => void;
  className?: string;
};

/**
 * Mapa interativo (React Flow). Somente leitura: arrastar o fundo/zoom navega, mas nós não se movem.
 * Há sempre a visão em lista como alternativa acessível (RoadmapLista).
 */
export function RoadmapCanvas({ nos, arestas, estados, selecionadoId, onSelecionar, className }: RoadmapCanvasProps) {
  const auto = React.useMemo(() => (semPosicoes(nos) ? autoLayout(nos, arestas) : null), [nos, arestas]);

  const flowNodes = React.useMemo<FlowNode[]>(
    () =>
      nos.map((no) => ({
        id: no.id,
        type: "no",
        position: posicaoDoNo(no, auto),
        width: LAYOUT.larguraNo,
        data: {
          no,
          estado: estados ? (estados[no.id] ?? "bloqueado") : null,
          selecionado: selecionadoId === no.id,
          onSelecionar,
        },
        draggable: false,
        connectable: false,
        selectable: false,
      })),
    [nos, auto, estados, selecionadoId, onSelecionar],
  );

  const flowEdges = React.useMemo<Edge[]>(
    () =>
      arestas.map((a) => {
        const acesa = arestaPercorrida(a, estados);
        const cor = acesa ? "hsl(var(--relp-done-fg))" : "hsl(var(--muted-foreground) / 0.55)";
        return {
          id: `${a.origem}->${a.destino}:${a.tipo}`,
          source: a.origem,
          target: a.destino,
          type: "smoothstep",
          focusable: false,
          style: { stroke: cor, strokeWidth: acesa ? 3 : 2, strokeDasharray: a.tipo === "opcional" ? "6 6" : undefined },
          markerEnd: { type: MarkerType.ArrowClosed, color: cor, width: 18, height: 18 },
        };
      }),
    [arestas, estados],
  );

  return (
    <div
      role="region"
      aria-label="Mapa do roadmap. Se preferir, use a visão em lista."
      className={cn("h-[min(72vh,680px)] min-h-[420px] w-full overflow-hidden rounded-xl border border-border bg-card", className)}
    >
      <ReactFlow
        nodes={flowNodes}
        edges={flowEdges}
        nodeTypes={nodeTypes}
        fitView
        fitViewOptions={{ padding: 0.2, maxZoom: 1 }}
        minZoom={0.25}
        maxZoom={1.5}
        nodesDraggable={false}
        nodesConnectable={false}
        nodesFocusable={false}
        edgesFocusable={false}
        elementsSelectable={false}
        zoomOnDoubleClick={false}
        // Sem um handler de clique no nó o React Flow aplica `pointer-events: none` nele (não é
        // arrastável nem selecionável), e o botão interno ficaria impossível de clicar.
        onNodeClick={(_, node) => onSelecionar(node.id)}
        onPaneClick={() => onSelecionar(null)}
        aria-label="Mapa interativo do roadmap"
      >
        <Background gap={24} size={1.5} color="hsl(var(--border))" />
        <Controls showInteractive={false} position="bottom-right" />
      </ReactFlow>
    </div>
  );
}

import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";

import { Icone } from "@/components/roadmap/icones";
import { RoadmapPlayer } from "@/components/roadmap/RoadmapPlayer";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import { EmptyState } from "@/components/ui/empty-state";
import { obterRoadmap } from "@/lib/roadmaps/client";

export const dynamic = "force-dynamic";

type Props = { params: { slug: string } };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const roadmap = await obterRoadmap(params.slug).catch(() => null);
  return { title: roadmap?.titulo ?? "Roadmap" };
}

export default async function RoadmapPage({ params }: Props) {
  const roadmap = await obterRoadmap(params.slug);
  if (!roadmap) notFound();

  const chefes = roadmap.nos.filter((n) => n.tipo === "chefe").length;

  return (
    <div className="space-y-6">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/roadmaps">Roadmaps</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{roadmap.titulo}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <header className="flex gap-4">
        <span
          aria-hidden="true"
          className="hidden h-14 w-14 shrink-0 items-center justify-center rounded-xl bg-primary-muted text-primary sm:flex"
        >
          <Icone nome={roadmap.icone} className="h-8 w-8" />
        </span>
        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight text-foreground">{roadmap.titulo}</h1>
          {roadmap.descricao && <p className="max-w-2xl text-muted-foreground">{roadmap.descricao}</p>}
          <p className="text-xs text-muted-foreground">
            {roadmap.nos.length} nós · {chefes} {chefes === 1 ? "chefe" : "chefes"} ·{" "}
            <span className="font-medium text-notify">{roadmap.xp_total} XP</span>
          </p>
        </div>
      </header>

      {roadmap.nos.length === 0 ? (
        <EmptyState title="Este roadmap ainda não tem nós" description="O autor ainda está montando o mapa." />
      ) : (
        <RoadmapPlayer roadmap={roadmap} />
      )}
    </div>
  );
}

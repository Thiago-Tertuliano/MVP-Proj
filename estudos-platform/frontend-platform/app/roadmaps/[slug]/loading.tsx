import { Skeleton } from "@/components/ui/skeleton";

export default function RoadmapLoading() {
  return (
    <div className="space-y-6" role="status" aria-label="Carregando roadmap">
      <Skeleton className="h-5 w-48" />
      <div className="space-y-3">
        <Skeleton className="h-9 w-72" />
        <Skeleton className="h-5 w-full max-w-xl" />
      </div>
      <Skeleton className="h-24 w-full rounded-xl" />
      <Skeleton className="h-96 w-full rounded-xl" />
      <span className="sr-only">Carregando…</span>
    </div>
  );
}

import { Skeleton } from "@/components/ui/skeleton";

export default function HomeLoading() {
  return (
    <div className="space-y-8" role="status" aria-label="Carregando trilhas">
      <div className="space-y-3">
        <Skeleton className="h-9 w-64" />
        <Skeleton className="h-5 w-full max-w-xl" />
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        {[0, 1].map((i) => (
          <Skeleton key={i} className="h-44 rounded-xl" />
        ))}
      </div>
      <span className="sr-only">Carregando…</span>
    </div>
  );
}

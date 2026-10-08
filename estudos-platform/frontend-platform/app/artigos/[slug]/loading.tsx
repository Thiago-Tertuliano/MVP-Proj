import { Skeleton } from "@/components/ui/skeleton";

export default function ArtigoLoading() {
  return (
    <div className="space-y-8" role="status" aria-label="Carregando artigo">
      <Skeleton className="h-5 w-64" />
      <div className="grid gap-8 lg:grid-cols-[minmax(0,1fr)_20rem]">
        <div className="space-y-4">
          <Skeleton className="h-10 w-3/4" />
          <Skeleton className="h-5 w-1/2" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-5/6" />
          <Skeleton className="h-40 w-full" />
        </div>
        <Skeleton className="hidden h-48 lg:block" />
      </div>
      <span className="sr-only">Carregando…</span>
    </div>
  );
}

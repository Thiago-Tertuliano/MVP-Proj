import Link from "next/link";

export function SiteHeader() {
  return (
    <header className="border-b border-surface-border bg-surface-raised/80 backdrop-blur">
      <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-4">
        <Link href="/" className="flex items-center gap-2 font-semibold tracking-tight">
          <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-accent text-accent-foreground">
            EP
          </span>
          Relp!
        </Link>
        <nav className="flex items-center gap-3 text-sm">
          <Link href="/" className="text-muted-foreground transition hover:text-foreground">
            Trilhas
          </Link>
          <Link href="/login" className="btn-ghost py-1.5">
            Entrar
          </Link>
        </nav>
      </div>
    </header>
  );
}

import Link from "next/link";

export function SiteHeader() {
  return (
    <header className="border-b border-border bg-surface-raised/80 backdrop-blur">
      <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-4">
        <Link href="/" className="flex items-center gap-2 font-semibold tracking-tight text-fg">
          <span className="flex h-8 w-8 items-center justify-center rounded-relp-md bg-primary-muted text-primary">
            EP
          </span>
          Estudos Platform
        </Link>
        <nav className="flex items-center gap-3 text-sm">
          <Link href="/" className="text-muted transition hover:text-fg">
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

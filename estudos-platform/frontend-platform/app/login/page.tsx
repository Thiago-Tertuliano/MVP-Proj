export default function LoginPage() {
  return (
    <div className="mx-auto max-w-md space-y-6">
      <header className="space-y-2 text-center">
        <h1 className="text-2xl font-bold text-fg">Entrar</h1>
        <p className="text-sm text-muted-foreground">
          Mock de login — depois liga em POST /api/v1/auth/login
        </p>
      </header>

      <form className="card space-y-4">
        <label className="block space-y-1.5 text-sm">
          <span className="text-muted-foreground">E-mail</span>
          <input
            type="email"
            defaultValue="autor.seed@estudos.local"
            className="w-full rounded-relp-md border border-border bg-surface px-3 py-2 text-fg outline-none focus:border-primary"
          />
        </label>
        <label className="block space-y-1.5 text-sm">
          <span className="text-muted-foreground">Senha</span>
          <input
            type="password"
            defaultValue="senha1234"
            className="w-full rounded-relp-md border border-border bg-surface px-3 py-2 text-fg outline-none focus:border-primary"
          />
        </label>
        <button type="button" className="btn-primary w-full">
          Entrar
        </button>
      </form>
    </div>
  );
}

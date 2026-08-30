export default function LoginPage() {
  return (
    <div className="mx-auto max-w-md space-y-6">
      <header className="space-y-2 text-center">
        <h1 className="text-2xl font-bold">Entrar</h1>
        <p className="text-sm text-zinc-400">
          Mock de login — depois liga em POST /api/v1/auth/login
        </p>
      </header>

      <form className="card space-y-4">
        <label className="block space-y-1.5 text-sm">
          <span className="text-zinc-400">E-mail</span>
          <input
            type="email"
            defaultValue="autor.seed@estudos.local"
            className="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 outline-none focus:border-accent"
          />
        </label>
        <label className="block space-y-1.5 text-sm">
          <span className="text-zinc-400">Senha</span>
          <input
            type="password"
            defaultValue="senha1234"
            className="w-full rounded-lg border border-surface-border bg-surface px-3 py-2 outline-none focus:border-accent"
          />
        </label>
        <button type="button" className="btn-primary w-full">
          Entrar
        </button>
      </form>
    </div>
  );
}

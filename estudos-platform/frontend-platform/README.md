# Frontend — Estudos Platform

Next.js 14 (App Router) + TypeScript. Consome a API em `backend-platform`.

## Telas (mock — fidelidade visual)

| Rota | Tela |
|------|------|
| `/` | Home — trilhas publicadas + rascunho |
| `/trilhas/go-basico` | Mapa roadmap (Sintaxe + Interfaces) |
| `/artigos/pacotes-em-go` | Aula Go com blocos + quiz |
| `/artigos/spark` | Stub de catálogo (link + checkpoint) |
| `/login` | Login mock |

Dados em `lib/mock-data.ts` — alinhados ao seed e ao content-job.

## Desenvolvimento

```powershell
cd estudos-platform\frontend-platform
copy .env.example .env.local
npm install
npm run dev
```

Abre [http://localhost:3000](http://localhost:3000). API padrão: `http://localhost:8080`.

## Storybook

```powershell
npm run storybook
```

Abre [http://localhost:6006](http://localhost:6006). Addons: essentials, a11y, themes (light/dark). Preview carrega `app/globals.css` (Tailwind + tokens).

```powershell
npm run build-storybook
```

## Scripts (CI)

| Comando | Uso |
|---------|-----|
| `npm run lint` | ESLint (`next lint`) |
| `npm run typecheck` | TypeScript sem emitir |
| `npm run build` | Build de produção |
| `npm run storybook` | DS visual (porta 6006) |
| `npm run build-storybook` | Build estático do Storybook |

A esteira **Frontend CI** (`.github/workflows/frontend.yml`) roda lint/typecheck/build em PRs/push que tocam `frontend-platform/`.

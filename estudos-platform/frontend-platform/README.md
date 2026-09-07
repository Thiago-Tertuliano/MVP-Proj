# Frontend — Relp!

Next.js 14 (App Router) + TypeScript. Consome a API em `backend-platform`.

## Arquitetura de pastas

Três camadas. Não misturar:

| Pasta | O que entra | O que **não** entra |
|-------|-------------|---------------------|
| `app/` | Rotas, layouts, pages (App Router). Só orquestra dados + monta UI. | Lógica visual reutilizável, primitivos, client HTTP genérico |
| `components/ui/` | Primitivos de mercado (shadcn/ui): Button, Input, Dialog… | Regras de negócio Relp (trilha, progresso, quiz) |
| `components/relp/` | Componentes de domínio: header, card de trilha, roadmap, renderer de aula | Primitivos genéricos (button “cru”), rotas |

```
frontend-platform/
├── app/                      # Rotas Next.js
│   ├── layout.tsx
│   ├── page.tsx              # Home (lista de trilhas)
│   ├── login/page.tsx
│   ├── trilhas/[slug]/page.tsx
│   └── artigos/[slug]/page.tsx
├── components/
│   ├── ui/                   # shadcn — Button, Input, …
│   └── relp/                 # Domínio Relp — SiteHeader, TrilhaCard, …
├── lib/                      # utils (cn), types, mock-data, futuro api.ts
├── stories/                  # Foundations do Storybook (cores, tipografia)
├── .storybook/
└── components.json           # Config shadcn (new-york)
```

**Regras rápidas**

1. Página em `app/` importa de `@/components/relp/*` e `@/components/ui/*` — nunca o contrário (`ui` não importa `relp`).
2. `relp` pode usar `ui` (ex.: botão shadcn dentro do card de trilha).
3. Stories: `components/ui/*.stories.tsx` e `components/relp/*.stories.tsx` (quando existirem); foundations em `stories/`.
4. Novo primitivo → `npx shadcn@latest add …` cai em `components/ui/`.
5. Novo bloco de produto (busca, anotação, progresso) → `components/relp/`.

## Telas (mock — fidelidade visual)

| Rota | Tela |
|------|------|
| `/` | Home — trilhas publicadas + rascunho |
| `/trilhas/go-basico` | Mapa roadmap (Sintaxe + Interfaces) |
| `/artigos/pacotes-em-go` | Aula Go com blocos + quiz |
| `/artigos/spark` | Stub de catálogo (link + checkpoint) |
| `/login` | Login mock |

Dados em `lib/mock-data.ts` — alinhados ao seed e ao content-job.

## Design system

- Primitivos shadcn/ui em `components/ui/` (Button, Input)
- Domínio Relp em `components/relp/`
- Ícones: `lucide-react`
- Utilitário `cn()` em `lib/utils.ts`
- Config: `components.json` (estilo new-york)

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

Abre [http://localhost:6006](http://localhost:6006). Addons: essentials, a11y, themes (light/dark). Preview carrega `app/globals.css` (Tailwind + tokens). Stories: `UI/Button`, `UI/Input`, `Foundations/Colors`.

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

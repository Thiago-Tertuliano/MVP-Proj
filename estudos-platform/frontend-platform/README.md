# Frontend — Relp!

Next.js 14 (App Router) + TypeScript + Tailwind + shadcn/ui. Consome a API em `backend-platform`
(`/api/v1`, cookies HttpOnly para a sessão).

## Arquitetura de pastas

Três camadas. Não misturar:

| Pasta | O que entra | O que **não** entra |
|-------|-------------|---------------------|
| `app/` | Rotas, layouts, pages (App Router). Só orquestra dados + monta UI. | Lógica visual reutilizável, primitivos, client HTTP genérico |
| `components/ui/` | Primitivos de mercado (shadcn/ui): Button, Input, Dialog… | Regras de negócio Relp (trilha, progresso, quiz) |
| `components/relp/` | Componentes de domínio: header, card de trilha, roadmap, renderer de aula, formulários de auth | Primitivos genéricos (button “cru”), rotas |

```
frontend-platform/
├── app/                      # Rotas Next.js
│   ├── layout.tsx            # SessionProvider + header + Toaster
│   ├── page.tsx              # Home: trilhas + "Continue de onde parou"
│   ├── login/ registro/      # Autenticação
│   ├── trilhas/[slug]/       # Mapa (roadmap) da trilha
│   ├── artigos/[slug]/       # Leitura + progresso + anotações + quiz
│   ├── busca/                # Resultados da busca (?q=)
│   └── error.tsx / not-found.tsx / loading.tsx
├── components/
│   ├── ui/                   # shadcn — Button, Input, Select, Dialog, Alert, …
│   └── relp/                 # Domínio Relp — AppHeader, TrilhaCard, RoadmapMap, ArticleRenderer, AuthForm, …
├── lib/
│   ├── api.ts                # Client HTTP único (credentials, erros PT-BR, refresh silencioso)
│   ├── session.tsx           # SessionProvider / useSession (login, registro, logout)
│   ├── content.ts            # Leitura de conteúdo no servidor (RSC)
│   ├── progress-store.ts     # Cache de progresso com optimistic UI
│   ├── roadmap.ts            # Modelo do mapa (lido / atual / futuro), anterior/próximo
│   └── *.test.ts             # Vitest
├── stories/                  # Foundations (paleta) + fixtures SÓ para Storybook
├── .storybook/
└── components.json           # Config shadcn (new-york)
```

**Regras rápidas**

1. Página em `app/` importa de `@/components/relp/*` e `@/components/ui/*` — nunca o contrário (`ui` não importa `relp`).
2. `relp` pode usar `ui` (ex.: botão shadcn dentro do card de trilha).
3. **Nada de dado fake no app.** `stories/fixtures` existe só para Storybook e o ESLint (`no-restricted-imports`) bloqueia o import fora de `*.stories.tsx`/testes.
4. Componente `relp` é **puro** (props → UI); a ligação com sessão/API fica em containers (`SessionHeader`, `TrilhaCardConnected`, `ArticleProgress`, …). Isso mantém tudo testável no Storybook.
5. Cores só por token (`bg-primary`, `text-done-fg`…); sem hex em componente. Verde = concluído/progresso, azul = marca/ação.
6. Novo primitivo → `npx shadcn@latest add …` cai em `components/ui/`.

## Telas

| Rota | Tela | Dados |
|------|------|-------|
| `/` | Home: cards de trilha (+ barra de progresso quando logado) e card **Continuar** | `GET /trilhas`, `/progresso/trilhas/{id}`, `/progresso/continuar` |
| `/trilhas/[slug]` | Mapa da trilha, módulos e nós lido/atual/futuro | `GET /trilhas/{slug}`, `/trilhas/{slug}/artigos` |
| `/artigos/[slug]` | Leitura, marcar como lido, anotações (autosave), quiz, anterior/próximo | `GET /artigos/{slug}`, `PUT /progresso/artigos/{id}`, `/artigos/{id}/anotacoes` |
| `/busca?q=` | Resultados (mínimo 2 caracteres) | `GET /busca?q=` |
| `/login`, `/registro` | Autenticação (`?next=` interno e validado) | `POST /auth/login`, `/auth/registrar` |

Estados cobertos em todas as telas de dados: **loading** (skeleton), **vazio** (`EmptyState`), **erro** (`error.tsx` / alertas com "Tentar de novo") e **404** (`not-found.tsx`).

### Sessão

- JWT só em cookie **HttpOnly** (nunca em `localStorage`). O navegador manda `credentials: "include"`.
- `localStorage["relp:sessao"]` guarda apenas uma *dica* ("já entrei aqui") para decidir se vale tentar o refresh silencioso — visitantes não gastam o rate limit (10/min) de `/auth/*`.
- Resposta 401 → um único `POST /auth/refresh` compartilhado entre chamadas simultâneas → repete a chamada original. Se o refresh falhar, volta a visitante e avisa por toast.
- **Use `localhost`, não `127.0.0.1`**: o CORS da API libera `http://localhost:3000` e os cookies não atravessam hosts diferentes.

### Tema

O app roda no tema **claro** (padrão). Os tokens `.dark` existem e são validados no Storybook (toolbar de tema), mas o app não expõe alternância ainda — decisão de produto pendente.

## Desenvolvimento

```powershell
cd estudos-platform\frontend-platform
copy .env.example .env.local
npm install
npm run dev
```

Abre [http://localhost:3000](http://localhost:3000). A API precisa estar no ar em `NEXT_PUBLIC_API_URL` (padrão `http://localhost:8080`) — veja `backend-platform/` (`docker compose up -d` + `go run ./cmd/api`). Usuário de demo (seed): `autor.seed@estudos.local` / `senha1234`, ou crie uma conta em `/registro`.

## Storybook

```powershell
npm run storybook
```

Abre [http://localhost:6006](http://localhost:6006). Addons: essentials, a11y, themes (light/dark), interactions e **pseudo-states** (hover/focus fixos via `parameters.pseudo`).

- `Foundations/Paleta de cores` — todos os tokens, com o valor HSL lido do CSS em tempo real (claro e escuro).
- `UI/*` — Button, formulário (Input, Textarea, Select, Checkbox, Radio, Switch), feedback (Badge, Alert, Toast, Progress, Skeleton, EmptyState) e estrutura (Card, Breadcrumb, Tabs, Dialog, Tooltip, Avatar).
- `Relp/*` — AppHeader (visitante/logado/carregando), Trilha (card, progresso, roadmap, continuar, paginação), Artigo (renderer, código, destaques, quiz, anotações), AuthForm (com testes de interação), Busca, Marca/Status.

```powershell
npm run build-storybook
```

## Scripts (CI)

| Comando | Uso |
|---------|-----|
| `npm run lint` | ESLint (`next lint`) |
| `npm run typecheck` | TypeScript sem emitir |
| `npm test` | Vitest (client HTTP, roadmap, store de progresso, validação, realce de código, redirect seguro) |
| `npm run build` | Build de produção |
| `npm run storybook` | DS visual (porta 6006) |
| `npm run build-storybook` | Build estático do Storybook |

A esteira **Frontend CI** (`.github/workflows/frontend.yml`) roda lint, typecheck, testes, build e build do Storybook em PRs/push que tocam `frontend-platform/`.

## QA / homologação

Checklist manual de aceite do MVP em [`../documentacao/QA-FRONT-MVP.md`](../documentacao/QA-FRONT-MVP.md).

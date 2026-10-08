# QA / Homologação — Front MVP Relp (F6-1)

Checklist de aceite do MVP. Itens marcados **[auto]** são cobertos por CI (lint, typecheck, vitest, build do Next e do Storybook).
Itens marcados **[smoke]** foram executados manualmente em ambiente local (API Go + Postgres + front) pelo autor do PR.
Os demais dependem de **Staging** e de aprovação humana — ficam em aberto até serem executados.

## 1. Ambiente

| Item | Como |
|------|------|
| API | `docker compose up -d` + `go run ./cmd/api` em `backend-platform` (porta 8080) |
| Front | `npm run dev` em `frontend-platform` (porta 3000) |
| CORS | `CORS_ALLOWED_ORIGINS` deve conter `http://localhost:3000`; acessar o front por `localhost` (não `127.0.0.1`) |
| Conteúdo | Trilha `go-basico` publicada (seed + content-job) |

## 2. Jornada do aluno

| # | Passo | Resultado esperado | Status |
|---|-------|--------------------|--------|
| 1 | Abrir `/` deslogado | Cards de trilha, header com **Entrar / Criar conta**, sem barras de progresso | [smoke] ok |
| 2 | `/registro` com campos vazios | Mensagens por campo, sem chamar a API | [smoke] ok |
| 3 | `/registro` com dados válidos (`?next=/trilhas/go-basico`) | Conta criada, toast, redireciona para `next`, header mostra avatar + nome | [smoke] ok |
| 4 | Registrar o mesmo e-mail de novo | Alerta "E-mail já cadastrado" (409) | [auto] mapeamento em `auth-errors`; validar em staging |
| 5 | Login com senha errada | Alerta "E-mail ou senha incorretos." (401), sem refresh em laço | [auto] `api.test.ts`; validar em staging |
| 6 | 11 tentativas seguidas de login | Alerta "Muitas tentativas… 1 minuto" (429) | [auto] `api.test.ts`; validar em staging |
| 7 | Recarregar a página logado | Continua logado (cookie HttpOnly + `/auth/me`) | [smoke] ok |
| 8 | Abrir `/trilhas/go-basico` | Mapa por módulos; 1º nó "Continue aqui"; barra 0% | [smoke] ok |
| 9 | Abrir um artigo, **Marcar como lido** | Botão muda na hora (optimistic); toast com ação "Próximo" | [smoke] ok |
| 10 | Voltar ao mapa | Nó verde, barra 33%, próximo nó "Continue aqui" | [smoke] ok |
| 11 | Digitar uma anotação | "Salvando…" → "Salvo"; reabrir o artigo mostra o texto | [smoke] ok (salvo); persistência após reload a confirmar em staging |
| 12 | Home logado | Card **Continue de onde parou** + barra no card da trilha | [smoke] ok |
| 13 | Busca `/busca?q=go` | Lista de artigos; `q` com 1 caractere pede mais letras; sem resultado → estado vazio | [smoke] ok (resultado); demais [auto] |
| 14 | **Sair** | Header volta a visitante, progresso do usuário some, `/progresso/*` deixa de ser chamado | [smoke] ok |
| 15 | **Completar 100% da trilha** (todos os artigos, com anotações) | Barra 100%, aviso "Você concluiu esta trilha" | **Pendente (Staging)** |
| 16 | Sessão expirada no meio da leitura (access token vencido) | Refresh silencioso; se falhar, toast "Sua sessão expirou" e volta a visitante | [auto] `api.test.ts`; validar em staging (TTL real de 15 min) |

## 3. Estados de erro e vazio

| Cenário | Esperado | Status |
|---------|----------|--------|
| API fora do ar ao abrir `/` | `error.tsx` com "Tentar de novo" (sem stack) | Pendente (derrubar a API e conferir) |
| `/artigos/slug-inexistente` ou artigo em rascunho | Página "não encontrada" | [smoke] ok (conteúdo; status HTTP 200 por causa do streaming do `loading.tsx` — comportamento do Next 14) |
| Trilha sem artigos publicados | `EmptyState` | [auto] Storybook |
| Falha ao carregar progresso | Alerta com "Tentar de novo"; mapa continua visível | Pendente |
| Falha ao salvar anotação | "Não foi possível salvar" + "Tentar de novo" | [auto] Storybook (`NoteEditor / estados`) |

## 4. Acessibilidade e responsividade

| Item | Como auditar | Status |
|------|--------------|--------|
| Contraste AA nos componentes | Storybook → aba **Accessibility** em cada story (claro e escuro) | **Pendente** — rodar e registrar violações |
| Teclado: "Pular para o conteúdo", foco visível, ordem lógica, Radix (Select/Dialog/Radio) | Navegar só com Tab/Enter/Espaço/Setas | **Pendente** |
| Leitor de tela: estados do roadmap (Lido / Próximo a ler), erros de formulário (`role=alert`), status de anotação (`aria-live`) | NVDA/VoiceOver | **Pendente** |
| Lighthouse (Accessibility ≥ 90) nas rotas `/`, `/trilhas/[slug]`, `/artigos/[slug]`, `/login` | Chrome DevTools → Lighthouse (build de produção) | **Pendente** |
| Responsividade 360 / 768 / 1280 px (header com busca em linha própria no mobile; coluna lateral do artigo abaixo do conteúdo) | DevTools device toolbar | **Pendente** |
| `prefers-reduced-motion` | Ativar no SO | Regra CSS global aplicada; validar |

## 5. Aprovação

| Item | Responsável | Status |
|------|-------------|--------|
| Zero bugs críticos no fluxo principal em Staging | QA | Pendente |
| Deploy em Staging concluído | Tech Lead | Pendente |
| Aprovação do Tech Lead para merge na `main` | Thiago Tertuliano | Pendente |

> A issue F6-1 só deve ser fechada depois dos itens **Pendente** acima. O PR de implementação entrega o código e o checklist; ele não substitui a homologação em Staging.

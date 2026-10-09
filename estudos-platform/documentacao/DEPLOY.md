# Deploy gratuito — Vercel + Render + Neon

Publica o conteúdo da branch `dev` sem custo e sem domínio próprio.

| Peça | Serviço (plano free) | O que roda |
|------|----------------------|------------|
| Front | [Vercel](https://vercel.com) Hobby | `frontend-platform` (Next 14) |
| API | [Render](https://render.com) Free | `backend-platform/Dockerfile` (roda `migrate` e sobe a API) |
| Banco | [Neon](https://neon.tech) Free | Postgres com `pgvector` e `pgcrypto` |

## Como as peças conversam

```
navegador ──► https://<front>.vercel.app/api/v1/*  ──(rewrite do Next)──►  https://<api>.onrender.com/api/v1/*  ──►  Neon
```

Sem domínio próprio, `vercel.app` e `onrender.com` são sites diferentes e o navegador descartaria os
cookies de sessão (HttpOnly, `SameSite=Lax`). Por isso o browser fala só com o front
(`NEXT_PUBLIC_API_URL=/`) e o Next repassa `/api/v1/*` para `API_URL` (`next.config.mjs`).
Os cookies ficam first-party e o CORS deixa de importar para o browser.

## Limitações do plano gratuito

| Limitação | Impacto | Mitigação |
|-----------|---------|-----------|
| Render hiberna após 15 min sem tráfego | 1ª requisição leva ~50s; o proxy do Vercel pode dar timeout | Workflow `keepalive.yml` pinga `/health` a cada 10 min |
| Render free: 750 h/mês | Uma única API ligada 24/7 cabe (~744 h) | Não suba um segundo serviço free na mesma conta |
| Neon free: 0,5 GB e compute suspende quando ocioso | Primeira query após pausa demora ~1s | Nenhuma necessária para MVP |
| Rate limit em memória | Reinício da API zera os contadores | Aceitável com 1 instância |

## Passo a passo

### 1. Neon (banco)

1. Crie um projeto (região mais próxima do Render, ex.: `AWS us-east-1` se a API ficar em Virginia, ou `AWS sa-east-1`).
2. Copie a *connection string* (formato `postgres://...neon.tech/neondb?sslmode=require&channel_binding=require`).
   O `channel_binding` é removido automaticamente pela API.
3. Não é preciso criar extensões à mão: as migrations criam `vector` e `pgcrypto`.

### 2. Render (API)

Opção A — Blueprint: **New → Blueprint**, aponte para o repositório; o `render.yaml` da raiz cria o serviço `relp-api`.
Preencha as variáveis marcadas como `sync: false`.

Opção B — manual: **New → Web Service**, repositório, branch `dev`, *Language* `Docker`,
*Root Directory* `estudos-platform/backend-platform`, *Instance type* `Free`, *Health Check Path* `/health`.

Variáveis de ambiente:

| Variável | Valor |
|----------|-------|
| `APP_ENV` | `production` (já é o padrão da imagem) |
| `DATABASE_URL` | connection string do Neon |
| `JWT_SECRET` | 32+ caracteres aleatórios (a API não sobe com menos em produção) |
| `TRUST_PROXY` | `true` |
| `CORS_ALLOWED_ORIGINS` | `https://<front>.vercel.app` |
| `EDITOR_EMAILS` | seu e-mail real (CSV). **Não** use `autor.seed@estudos.local` |

O Render injeta `PORT`; a API usa ela quando `APP_PORT` não existe. A cada deploy o container roda
`./migrate` antes da API. Em produção o `migrate` invalida a senha conhecida do usuário de seed
(`autor.seed@estudos.local`), então ninguém consegue entrar com `senha1234`.

Teste: `https://<api>.onrender.com/health` deve responder `ok`.

### 3. Conteúdo (trilhas e artigos)

As migrations criam só o esqueleto. Para importar o conteúdo de `fontes/` e `content/` no Neon, rode o
content-job da sua máquina apontando para o banco de produção:

```powershell
cd estudos-platform/backend-platform
$env:DATABASE_URL = "<connection string do Neon>"
go run ./cmd/content-job --dry-run   # confere antes
go run ./cmd/content-job
Remove-Item Env:DATABASE_URL
```

### 4. Vercel (front)

1. **Add New → Project**, importe o repositório.
2. *Root Directory*: `estudos-platform/frontend-platform` (framework detectado: Next.js).
3. *Production Branch* (em Settings → Git): `dev`.
4. Variáveis de ambiente:

| Variável | Valor |
|----------|-------|
| `NEXT_PUBLIC_API_URL` | `/` |
| `API_URL` | `https://<api>.onrender.com` (sem `/api/v1`, sem barra no fim) |

5. Faça o deploy e volte ao Render para ajustar `CORS_ALLOWED_ORIGINS` com a URL final do Vercel.

### 5. Keepalive (GitHub)

Em **Settings → Secrets and variables → Actions → Variables**, crie `API_HEALTH_URL` =
`https://<api>.onrender.com/health`. O workflow `Keepalive API` passa a rodar a cada 10 min
(agendamentos só rodam na branch padrão do repositório, hoje `dev`).

## Checklist pós-deploy

- [ ] `/health` da API responde `ok`
- [ ] Home do front carrega trilhas
- [ ] Criar conta em `/registro`, sair e entrar de novo (cookies funcionando)
- [ ] Recarregar a página logado mantém a sessão (refresh via `/api/v1/auth/refresh`)
- [ ] Login com `autor.seed@estudos.local` / `senha1234` **falha**
- [ ] Marcar artigo como lido persiste após recarregar

## Atualizações

Push/merge em `dev` dispara deploy automático no Render e no Vercel. Migrations novas rodam sozinhas
no boot da API.

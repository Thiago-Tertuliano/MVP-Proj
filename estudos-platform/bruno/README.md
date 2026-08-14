# Collection Bruno — Estudos Platform API

Collection para testar a API no [Bruno](https://www.usebruno.com/) (alternativa open-source ao Postman).

## Como abrir

1. Instale o Bruno: https://www.usebruno.com/downloads  
2. **Open Collection** → selecione a pasta:
   ```
   estudos-platform/bruno
   ```
3. Em cima, escolha o environment **Local**
4. Suba a API (`docker compose up -d` + `go run ./cmd/api` em `backend-platform`)

## Ordem sugerida de execução

| # | Pasta | Request | Observação |
|---|--------|---------|------------|
| 1 | `00-health` | Health | Tem que retornar `ok` |
| 2 | `01-auth` | Registrar **ou** Login | Salva tokens no env |
| 3 | `01-auth` | Me | Precisa de Bearer |
| 4 | `01-auth` | Refresh | Rotaciona tokens |
| 5 | `02-trilhas` | Criar → Módulo → Publicar → Listar → Obter | Branch com trilhas |
| 6 | `03-artigos` | Criar → Atualizar → Publicar → Listar → Obter | Precisa do PR #8 |
| 7 | `01-auth` | Logout | Precisa do PR #7 |

## Variáveis (environment Local)

| Var | Uso |
|-----|-----|
| `baseUrl` | `http://localhost:8080` |
| `email` / `senha` / `nome` | Credenciais de teste |
| `accessToken` / `refreshToken` | Preenchidos após Login/Registrar |
| `trilhaId` / `trilhaSlug` | Preenchidos após Criar Trilha |
| `artigoId` / `artigoSlug` | Preenchidos após Criar Artigo |

Se **Registrar** der `409` (e-mail já cadastrado), rode **Login** ou mude o `email` no environment.

## Compatibilidade por branch/PR

| Feature | Status |
|---------|--------|
| Health + Auth (registrar/login/refresh/me) | Base `dev` / sprint 1 |
| Logout | PR #7 |
| Trilhas | PR #9 / branch `feat/trilhas-modulos` |
| Artigos | PR #8 / branch `feat/artigos-crud` |

Requests de feature ainda não mergeada vão retornar **404** — normal até o merge.

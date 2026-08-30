# Fontes do content-job

Markdown de **entrada** para o binário `backend-platform/cmd/content-job`. Não edite `content/generated/` — é saída do job.

| Arquivo | Quem edita | O que vira |
|---------|------------|------------|
| `Courses.md` | Mariana | trilhas + artigos rascunho (links externos) |
| `AULA-CODE-REVIEW-GO-SPRINT.md` | Thiago | enriquece a trilha seed `go-basico` |

## Fluxo

1. PR alterando `fontes/**` dispara a **Content CI** (`content.yml`).
2. CI roda `go test` nos parsers e `content-job --dry-run`.
3. Após merge, Thiago/Bruno roda `make content-job` (Postgres) e Bruno publica via API.

Spec completa: [`documentacao/JOB-CONTEUDO.md`](../documentacao/JOB-CONTEUDO.md).

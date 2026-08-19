# Estudos Platform — Backend (Go + DDD)

[![CI](https://github.com/Thiago-Tertuliano/MVP-Proj/actions/workflows/ci.yml/badge.svg?branch=dev)](https://github.com/Thiago-Tertuliano/MVP-Proj/actions/workflows/ci.yml)

API monolítica modular em Go para a plataforma de estudos de TI. Segue **Domain-Driven Design (DDD)** estrito, **Vertical Slicing** por funcionalidade e **PostgreSQL + pgvector** para conteúdo dinâmico e busca semântica.

PRs contra `dev`: branch `tipo/descricao` (`feat/`, `fix/`, `bug/`, `ref/`, `ci/`, `docs/`, `chore/`) e título conventional (`feat: …`, `fix(auth): …`). Sem isso o job `branch-name` falha e o build/lint não rodam.

---

## 📌 Stack

| Camada | Tecnologia |
|---|---|
| Linguagem | Go 1.26 |
| HTTP | `chi` router |
| Banco | PostgreSQL 16 + `pgvector` |
| Acesso a dados | `pgx` (pool) + migrations SQL |
| Validação | `go-playground/validator` |
| Config | env vars (dotenv) |
| Build | Docker + Makefile |
| Testes | `testing` padrão + `testcontainers-go` |

---

## 🗂️ Estrutura de Pastas

```
backend-platform/
├── cmd/
│   └── api/
│       └── main.go                     # Entry point: config, DI, sobe HTTP server
│
└── internal/
    ├── applications/                   # 🎯 CAMADA DE APLICAÇÃO
    │   └── estudos/                    # Bounded Context: Estudos
    │       ├── usecase/                #   Casos de uso (orquestração)
    │       ├── dto/                    #   Contratos de entrada/saída
    │       └── port/                   #   Ports: interfaces que a infra implementa
    │
    ├── domain/                         # 🧠 CAMADA DE DOMÍNIO (pura)
    │   ├── shared/                     #   Código compartilhado entre contextos
    │   │   ├── kernel/                 #   Base Entity, ValueObject, DomainEvent
    │   │   └── errors/                 #   DomainError padronizado
    │   └── estudos/
    │       ├── entity/                 #   Entidades/Aggregates (Artigo, Trilha)
    │       ├── valueobject/            #   Slug, Embedding, ConteudoJSON
    │       ├── repository/             #   Interfaces (portas) de repositório
    │       └── event/                  #   Eventos de domínio
    │
    ├── infrastructure/                 # 🔧 CAMADA DE INFRAESTRUTURA
    │   ├── persistence/
    │   │   └── postgres/
    │   │       ├── repository/         #   Implementação concreta dos repositórios
    │   │       └── migration/          #   SQL migrations versionadas
    │   └── external/                   #   Integrações: OpenAI, S3, etc.
    │
    └── presentation/                   # 🌐 CAMADA DE APRESENTAÇÃO
        └── http/
            ├── handler/                #   Handlers HTTP (binding + validação de borda)
            ├── middleware/             #   Auth, CORS, RequestID, Recovery
            └── router/                 #   Composição de rotas (chi)
```

---

## 📏 Regra de Dependência (Arquitetura em Camadas)

```
┌─────────────────────────────────────────────────────────────┐
│  PRESENTATION  (handler → usecase)                          │
│  ├── importa: applications, dto                             │
│  └── NUNCA importa: infrastructure, domain*                 │
├─────────────────────────────────────────────────────────────┤
│  APPLICATION  (usecase → domain)                            │
│  ├── importa: domain, dto, port                             │
│  └── NUNCA importa: infrastructure, presentation            │
├─────────────────────────────────────────────────────────────┤
│  DOMAIN  (puro — zero dependências externas)                │
│  └── NUNCA importa: infrastructure, application, present.   │
├─────────────────────────────────────────────────────────────┤
│  INFRASTRUCTURE  (implementa os ports do domínio)           │
│  └── importa: domain (para satisfazer as interfaces)        │
└─────────────────────────────────────────────────────────────┘
```

> **Regra de ouro:** a seta de dependência **sempre aponta para dentro**. O domínio não conhece o mundo externo — a infra conhece o domínio.

---

## 🧠 Camada de Domínio (detalhamento)

### O que vive aqui
- **Entidades** (`entity/`) — possuem identidade própria (ID UUID). Ex.: `Artigo`, `Trilha`.
- **Value Objects** (`valueobject/`) — sem identidade, igualdade por valor, **imutáveis**. Ex.: `Slug`, `Embedding`, `ConteudoJSON`.
- **Interfaces de repositório** (`repository/`) — contratos que a infra implementa. O domínio define o "quê", a infra decide o "como".
- **Eventos de domínio** (`event/`) — fatos que ocorreram e interessam a outros módulos.

### O que NÃO vive aqui
- Queries SQL, ORM, HTTP, cache, filas.
- Nenhum import de biblioteca externa que não seja da stdlib do Go.

### Exemplo de entidade — `Artigo`

```go
// internal/domain/estudos/entity/artigo.go
package entity

type Artigo struct {
	ID          string      `json:"id"`
	Slug        string      `json:"slug"`
	Titulo      string      `json:"titulo"`
	Conteudo    []byte      `json:"conteudo"`    // JSONB
	Metadados   []byte      `json:"metadados"`   // JSONB
	Embedding   []float32   `json:"embedding"`   // pgvector (1536 dims)
	Status      string      `json:"status"`      // rascunho | revisao | publicado | arquivado
	TrilhaID    *string     `json:"trilha_id"`
	ModuloID    *string     `json:"modulo_id"`
	AutorID     string      `json:"autor_id"`
	PublicadoEm *time.Time  `json:"publicado_em"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Regras de negócio — comportamento dentro da entidade
func (a *Artigo) Publicar() error {
	if a.Status != "revisao" {
		return errors.ErrInvalidState("Artigo.Publicar", "só pode publicar artigo em revisão", nil)
	}
	if len(a.Conteudo) == 0 {
		return errors.ErrInvalidArgument("Artigo.Publicar", "conteúdo vazio", nil)
	}
	a.Status = "publicado"
	a.PublicadoEm = now()
	return nil
}
```

---

## 🎯 Camada de Aplicação (detalhamento)

### O que vive aqui
- **Casos de uso** (`usecase/`) — orquestram: recebem DTO → validam regras de aplicação → chamam domínio/repositório → retornam DTO.
- **DTOs** (`dto/`) — contratos de entrada (`CreateArtigoRequest`) e saída (`ArtigoResponse`).
- **Ports** (`port/`) — interfaces que a infraestrutura implementa e injeta no use case (inversão de dependência).

### Exemplo de caso de uso — `CriarArtigo`

```go
// internal/applications/estudos/usecase/criar_artigo.go
package usecase

type CriarArtigoUseCase struct {
	repo        repository.ArtigoRepository
	slugger     *slug.Generator
	embedPort   port.EmbeddingProvider
}

type CriarArtigoRequest struct {
	Titulo    string `json:"titulo" validate:"required,min=3,max=300"`
	Conteudo  []byte `json:"conteudo" validate:"required"`
	TrilhaID  *string `json:"trilha_id"`
}

func (uc *CriarArtigoUseCase) Execute(ctx context.Context, req CriarArtigoRequest, autorID string) (*dto.ArtigoResponse, error) {
	// 1. valida DTO de borda
	// 2. monta entidade (regras de domínio)
	// 3. persiste via repositório
	// 4. dispara evento ArtigoCriado (embedding async)
	// 5. retorna DTO de saída
}
```

---

## 🔧 Camada de Infraestrutura (detalhamento)

### O que vive aqui
- **Implementações de repositório** (`postgres/repository/`) — SQL via `pgx`, satisfazem as interfaces do domínio.
- **Migrations** (`postgres/migration/`) — SQL versionado (`0001_create_artigos.up.sql`).
- **pgvector** — índice HNSW e queries de similaridade.
- **Externals** (`external/`) — OpenAI (embeddings), S3/R2 (armazenamento de imagens).

### Exemplo — repositório PostgreSQL

```go
// internal/infrastructure/persistence/postgres/repository/artigo_repo_pg.go
package repository

type ArtigoRepoPG struct {
	pool *pgxpool.Pool
}

func NewArtigoRepoPG(pool *pgxpool.Pool) *ArtigoRepoPG {
	return &ArtigoRepoPG{pool: pool}
}

func (r *ArtigoRepoPG) Save(ctx context.Context, a *entity.Artigo) error {
	// INSERT INTO artigos (...) ON CONFLICT (id) DO UPDATE
}
```

---

## 🌐 Camada de Apresentação (detalhamento)

### O que vive aqui
- **Handlers** (`http/handler/`) — recebem HTTP, fazem binding/validação de borda e chamam o use case. **Nunca** contêm regra de negócio.
- **Middlewares** (`http/middleware/`) — autenticação JWT, CORS, RequestID, Recovery (panic → 500), logging.
- **Router** (`http/router/`) — composição das rotas com chi, versionamento (`/api/v1`).

### Exemplo — handler

```go
// internal/presentation/http/handler/artigo_handler.go
package handler

type ArtigoHandler struct {
	createUC *usecase.CriarArtigoUseCase
}

func (h *ArtigoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. decode do body
	// 2. validação de borda
	// 3. chama use case
	// 4. serializa resposta (201)
}
```

---

## 🛣️ Fluxo de uma Requisição (Vertical Slice)

```
POST /api/v1/artigos
        │
        ▼
┌────────────────────────┐
│  Handler (http)        │  ← decode + validação de borda
└───────────┬────────────┘
            ▼
┌────────────────────────┐
│  UseCase (application) │  ← orquestração + regras de aplicação
└───────────┬────────────┘
            ▼
┌────────────────────────┐
│  Entity + Repository   │  ← regras de domínio + persistência (interface)
└───────────┬────────────┘
            ▼
┌────────────────────────┐
│  RepoPG (infra)        │  ← SQL real + pgvector
└───────────┬────────────┘
            ▼
┌────────────────────────┐
│  Event Bus (async)     │  ← ArtigoCriado → gera embedding
└────────────────────────┘
```

---

## 🗄️ Banco de Dados

### Tabelas núcleo

| Tabela | Tipo | Uso |
|---|---|---|
| `usuarios` | Relacional | Auth, perfil, plano |
| `trilhas` | Relacional | Trilhas de estudo (aggregate root) |
| `modulos` | Relacional | Módulos dentro de trilhas |
| `artigos` | JSONB + pgvector | Conteúdo rico + embedding p/ busca semântica |
| `progresso_estudo` | Relacional | Progresso do usuário |
| `anotacoes` | JSONB | Highlights/notas do estudante |

### Índices essenciais (pgvector)

```sql
-- Busca semântica (ANN — approximate nearest neighbor)
CREATE INDEX idx_artigos_embedding_hnsw ON artigos
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- Consultas em JSONB
CREATE INDEX idx_artigos_conteudo_gin ON artigos USING GIN (conteudo);
```

### Busca semântica

```
GET /api/v1/busca?q=como evoluir monólito para microserviços
```
1. Gera embedding da query (OpenAI, 1536 dims)
2. `SELECT ... ORDER BY embedding <=> $1 LIMIT 10`
3. Retorna metadados + trecho de contexto

---

## 🚀 Como Rodar

### Pré-requisitos
- Go 1.26+
- Docker (para Postgres local)

### 1. Subir o banco (pgvector)
```bash
docker run --name estudos-pg \
  -e POSTGRES_DB=estudos_platform \
  -e POSTGRES_USER=estudos \
  -e POSTGRES_PASSWORD=estudos_dev \
  -p 5432:5432 \
  -d pgvector/pgvector:pg16
```

### 2. Configurar ambiente
```bash
cp .env.example .env   # preencha os valores
```

### 3. Rodar migrations
```bash
make migrate-up
```

### 4. Subir a API (hot reload)
```bash
make dev
```

### 5. Testar
```bash
curl http://localhost:8080/health
```

---

## 🧪 Testes

```bash
make test        # unitários + integração (testcontainers)
make test-cov    # com cobertura
```

| Tipo | Onde | O que cobre |
|---|---|---|
| Unitário | junto a cada pacote | Regras de domínio, VOs, use cases |
| Integração | `internal/infrastructure` | Repositórios reais contra Postgres em container |
| API (e2e) | `internal/presentation` | Fluxo HTTP completo (handler → banco) |

---

## 📜 Makefile (comandos principais)

```makefile
dev:          # roda API com hot-reload (air)
build:        # compila binário otimizado
migrate-up:   # aplica migrations
migrate-down: # reverte migrations
test:         # roda todos os testes
lint:         # golangci-lint
generate:     # gera código (sqlc, mocks)
```

---

## 📐 Decisões de Arquitetura (resumo)

| Decisão | Motivo |
|---|---|
| Go + DDD estrito | Performance nativa, domínio rico e isolado, fácil de testar |
| `chi` em vez de Gin/Echo | Leve, padrão `net/http`, middlewares explícitos |
| `pgx` direto (sem ORM) | Controle total do SQL, zero magic, suporte nativo a pgvector |
| Interface no domínio | Inversão de dependência: teste use cases com mock fácil |
| JSONB para conteúdo | Flexibilidade total do editor de artigos (blocos ricos) |
| Embedding assíncrono | Criação do artigo não bloqueia na geração do vetor |
| Vertical Slicing | Cada feature ponta-a-ponta (migration → SQL → use case → HTTP) |

---

## 📚 Referências

- [The Clean Architecture — Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design — Eric Evans](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [pgvector](https://github.com/pgvector/pgvector)
- [chi router](https://github.com/go-chi/chi)
- [pgx](https://github.com/jackc/pgx)

## Creators
-Bruno
-Thiago
# Aula — Code Review do Backend Go (Estudos Platform)

Documento didático para explicar o que construímos no MVP, como ler os PRs e a **sintaxe Go** por trás de cada decisão.

| Item | Valor |
|------|--------|
| Projeto | `estudos-platform/backend-platform` |
| Linguagem | Go 1.26 |
| Arquitetura | DDD em camadas + vertical slice |
| Branch base | `dev` |
| Repositório | [Thiago-Tertuliano/MVP-Proj](https://github.com/Thiago-Tertuliano/MVP-Proj) |

---

## 1. O que é este projeto (em uma frase)

API HTTP em Go para uma plataforma de estudos de TI: autenticação JWT, artigos com conteúdo JSONB e trilhas com módulos — tudo organizado em Domain-Driven Design.

---

## 2. Mapa mental das camadas (DDD)

```
Presentation  →  Application  →  Domain  ←  Infrastructure
   (HTTP)         (use cases)    (puro)      (Postgres, JWT, bcrypt)
```

| Pasta | Responsabilidade | Pode importar |
|-------|------------------|---------------|
| `internal/domain/` | Regras de negócio, entidades, VOs, interfaces de repo | stdlib (+ uuid) |
| `internal/applications/` | Orquestra casos de uso, DTOs, ports | domain |
| `internal/infrastructure/` | SQL, JWT, bcrypt, config | domain (+ libs externas) |
| `internal/presentation/` | Handlers HTTP, middleware, router | applications (DTOs/use cases) |
| `cmd/api/` | `main` — sobe o servidor e faz DI | infrastructure + presentation |

**Regra de ouro (Dependency Rule):** dependências apontam para **dentro**. O domínio nunca importa HTTP, SQL ou JWT.

### Por que `internal/`?

Em Go, o diretório `internal/` é especial: **só o próprio módulo** pode importar pacotes dentro dele. É encapsulamento no nível do módulo — ninguém de fora do repo consegue importar `.../internal/domain/...`.

---

## 3. Sintaxe Go — guia rápido para a aula

Use esta seção como “dicionário” enquanto revisa o código.

### 3.1 Pacotes e arquivos

```go
package entity   // todo arquivo .go na pasta declara o mesmo package

import (
    "fmt"                                    // stdlib
    "github.com/google/uuid"                 // externo (go.mod)
    ".../internal/domain/shared/errors"      // interno do módulo
)
```

- Um diretório = um pacote (em geral).
- Nome exportado começa com **maiúscula** (`Usuario`, `NewEmail`).
- Nome privado começa com **minúscula** (`nome`, `normalizeSlug`).

### 3.2 Struct, campos e métodos

```go
type Usuario struct {
    nome  string              // privado
    email valueobject.Email   // VO embutido
}

// Método com receiver por valor (não altera o struct original)
func (u Usuario) Nome() string { return u.nome }

// Método com receiver por ponteiro (pode alterar)
func (u *Usuario) AlterarSenha(h valueobject.SenhaHash) {
    u.senhaHash = h
}
```

**Ponteiro `*T` vs valor `T`:**

| Receiver | Quando usar |
|----------|-------------|
| `*Usuario` | Precisa mutar o estado |
| `Usuario` | Só lê; cópia barata (structs pequenas) |

### 3.3 Embedding (composição)

Go **não tem herança de classes**. Tem *embedding*:

```go
type Usuario struct {
    kernel.BaseEntity  // embed: Usuario “ganha” ID(), CreatedAt(), Touch()
    nome string
}
```

Em aula: “é como composição automática — os métodos de `BaseEntity` sobem para `Usuario`”.

### 3.4 Interfaces (contratos)

```go
type TokenGerador interface {
    Gerar(claims Claims, accessTTL, refreshTTL time.Duration) (*TokenPar, error)
    ValidarAccessToken(token string) (*Claims, error)
}
```

- Interfaces são **satisfeitas implicitamente** (não existe `implements`).
- Qualquer tipo com esses métodos vira `TokenGerador`.
- Por isso dá para mockar nos testes sem frameworks pesados.

### 3.5 Múltiplos retornos e `error`

```go
func NewEmail(raw string) (Email, error) {
    // ...
    return Email{}, errors.ErrInvalidArgument("email inválido", "...", nil)
}
```

Idioma Go: último retorno costuma ser `error`. Caller sempre verifica:

```go
email, err := valueobject.NewEmail(req.Email)
if err != nil {
    return nil, err
}
```

### 3.6 `errors.As` / `errors.Is` (Go 1.13+)

```go
var de *errors.DomainError
if stderrors.As(err, &de) && de.Kind == errors.NotFound {
    // trata NotFound
}
```

`As` “desembrulha” a cadeia de erros até achar o tipo. Combinamos com `Unwrap()` no `DomainError`.

### 3.7 `context.Context`

Quase toda função de I/O recebe `ctx context.Context` **como primeiro parâmetro**:

```go
func (uc *LoginUsuario) Execute(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
```

Serve para cancelamento, timeout e valores de request (ex.: `usuario_id` no middleware).

### 3.8 Goroutines e channels (só no `main`)

```go
go func() {
    srv.ListenAndServe()  // roda em paralelo
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit  // bloqueia até Ctrl+C
```

Em aula: “`go` dispara concorrência; `chan` sincroniza o shutdown gracioso”.

### 3.9 Testes

Arquivo `*_test.go` no mesmo pacote:

```go
func TestLogin_Sucesso(t *testing.T) {
    // arrange → act → assert
    if err != nil {
        t.Fatalf("não deveria retornar erro: %v", err)
    }
}
```

Rodar: `go test ./internal/... -count=1`

---

## 4. Base já existente (Sprint 1 — Auth)

Antes dos PRs #5–#9, o backend já tinha o vertical slice de autenticação.

### 4.1 Entry point — `cmd/api/main.go`

**O que faz:** carrega config, abre pool Postgres, sobe `http.Server`, espera sinal e faz `Shutdown` com timeout.

**Pontos para a aula:**

1. `defer pool.Close()` — cleanup garantido ao sair da função.
2. Timeouts no server (`ReadTimeout`, `WriteTimeout`) — proteção contra conexões lentas.
3. Graceful shutdown — não corta request no meio.

### 4.2 Entidade `Usuario` — encapsulamento

```go
type Usuario struct {
    kernel.BaseEntity
    nome      string                 // minúsculo = privado ao pacote
    email     valueobject.Email
    senhaHash valueobject.SenhaHash
    status    StatusConta
}
```

**Code review (positivo):**

- Campos privados + getters públicos → ninguém de fora muda `email` sem passar por regra.
- `NovoUsuario` valida na construção (factory).
- `ReconstruirUsuario` separa “criar novo” de “hidratar do banco” — padrão clássico em DDD + Go.

**Sintaxe:** `type StatusConta string` é um *defined type* — tipagem forte sobre `string` (não mistura com qualquer string sem conversão).

### 4.3 Value Object `Email`

```go
func NewEmail(raw string) (Email, error) {
    raw = strings.TrimSpace(strings.ToLower(raw))
    // valida com net/mail
    return Email{value: raw}, nil
}
```

VO = imutável, igualdade por valor, validação na criação. Em Go isso vira struct pequena + construtor.

### 4.4 Ports (inversão de dependência)

```go
// application/port
type SenhaHasher interface {
    Hash(plain string) (string, error)
    Comparar(hash, plain string) bool
}
```

Use case depende da **interface**; `bcrypt_hasher.go` na infra implementa. Trocar bcrypt por argon2 = mudar só a infra.

### 4.5 `DomainError` → HTTP

```go
switch de.Kind {
case domainErros.InvalidArgument:
    status = http.StatusBadRequest      // 400
case domainErros.AlreadyExists:
    status = http.StatusConflict        // 409
case domainErros.Unauthorized:
    status = http.StatusUnauthorized    // 401
}
```

O domínio fala em linguagem de negócio (`AlreadyExists`); a presentation traduz para HTTP. Domínio **não** conhece status codes.

### 4.6 Endpoints da Sprint 1

| Método | Rota | Auth? | Função |
|--------|------|-------|--------|
| POST | `/api/v1/auth/registrar` | Não | Cria usuário + tokens |
| POST | `/api/v1/auth/login` | Não | Valida senha + tokens |
| POST | `/api/v1/auth/refresh` | Não | Rotaciona refresh |
| GET | `/api/v1/auth/me` | Bearer | Perfil do token |

**Detalhe de segurança (aula):** no login, usuário inexistente e senha errada retornam a **mesma** mensagem (`credenciais inválidas`) — anti-enumeração de e-mails.

**Refresh:** o token em texto puro **nunca** vai para o banco; só o hash SHA-256.

---

## 5. PRs abertos — code review didático

Ordem sugerida de merge (e de apresentação na aula): **#5 → #6 → #7 → #8 → #9**.

---

### PR #5 — `fix(auth): persiste refresh token no registro`

🔗 https://github.com/Thiago-Tertuliano/MVP-Proj/pull/5  
Branch: `feat/auth-register-persist-refresh`

#### Problema

O login salvava o refresh no Postgres. O register gerava o token na resposta, mas **não persistia**. Resultado: usuário se registra → chama `/auth/refresh` → 401.

#### O que mudou

`RegistrarUsuario` passou a receber `RefreshTokenRepository` (mesmo fluxo do login):

```go
type RegistrarUsuario struct {
    repo    repository.UsuarioRepository
    refresh repository.RefreshTokenRepository  // novo
    hasher  port.SenhaHasher
    tokens  port.TokenGerador
    cfg     RegistrarConfig
}
```

Depois de `Gerar(...)`:

```go
tokenHash := sha256.Sum256([]byte(tokens.RefreshToken))
rt := &repository.RefreshToken{
    ID:        uuid.New().String(),
    UsuarioID: usuario.ID().String(),
    TokenHash: hex.EncodeToString(tokenHash[:]),
    ExpiraEm:  tokens.RefreshExp,
}
uc.refresh.Save(ctx, rt)
```

#### Sintaxe em destaque

- `sha256.Sum256` retorna `[32]byte` (array). `hash[:]` vira slice para `hex.EncodeToString`.
- Injeção via construtor `NewRegistrarUsuario(...)` — Composition Root no `router.go`.

#### Code review

| ✅ Bom | ⚠️ Observação |
|--------|----------------|
| Paridade login/register | Register ainda não é atômico (salva user, depois token). Falha no 2º passo deixa user sem refresh — aceitável no MVP |
| Teste cobre “refresh foi salvo” | — |

**Pergunta para a turma:** por que hashear o refresh e não guardar o JWT cru?

---

### PR #6 — `feat(http): adiciona middleware CORS`

🔗 https://github.com/Thiago-Tertuliano/MVP-Proj/pull/6  
Branch: `feat/http-cors-middleware`

#### Problema

Frontend (Next/Vite em `localhost:3000`) seria bloqueado pelo browser sem CORS.

#### O que mudou

Middleware próprio (sem lib extra):

```go
type CORS struct {
    allowedOrigins map[string]struct{}  // set idiomático em Go
    allowAll       bool
}

func (c *CORS) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // seta headers se Origin permitida
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent) // preflight
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

Config via env:

```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

#### Sintaxe em destaque

**`map[string]struct{}` como set**

```go
allowedOrigins map[string]struct{}
c.allowedOrigins[o] = struct{}{}  // valor vazio ocupa 0 bytes
_, ok := c.allowedOrigins[origin]
```

Em Go não existe `Set<T>` nativo; esse padrão é o idiomático.

**Closure que captura `next`:**

```go
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // ...
    next.ServeHTTP(w, r)
})
```

`http.HandlerFunc` é um *adapter*: converte `func(ResponseWriter, *Request)` em `http.Handler`.

**Chain no Chi:**

```go
r.Use(middleware.RequestID)
r.Use(middleware.Recovery)
r.Use(middleware.NewCORS(cfg.CORSAllowedOrigins).Handler)
```

Ordem importa: Recovery por fora evita panic sem resposta; CORS cedo para preflight.

#### Code review

| ✅ Bom | ⚠️ Observação |
|--------|----------------|
| Allowlist explícita | `*` + `Allow-Credentials: true` ecoa a Origin (correto); não usar `*` em produção |
| Testes de origem permitida/negada/preflight | — |

---

### PR #7 — `feat(auth): endpoint POST /auth/logout`

🔗 https://github.com/Thiago-Tertuliano/MVP-Proj/pull/7  
Branch: `feat/auth-logout`

#### O que faz

Rota **protegida** `POST /api/v1/auth/logout`:

1. Middleware JWT coloca `usuario_id` no `context`.
2. Use case chama `RevokeAllByUser`.
3. Todos os refresh daquele usuário ficam `revogado = true`.

```go
func (uc *LogoutUsuario) Execute(ctx context.Context, usuarioID string) error {
    id, err := uuid.Parse(usuarioID)
    if err != nil {
        return errors.ErrInvalidArgument("usuario_id inválido", "...", err)
    }
    return uc.refresh.RevokeAllByUser(ctx, id)
}
```

No handler:

```go
usuarioID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
```

#### Sintaxe em destaque

**Type assertion**

```go
usuarioID, _ := r.Context().Value(middleware.CtxUsuarioID).(string)
```

`Value` retorna `any` (`interface{}`). `(string)` tenta converter. O `_` ignora o `ok` — em produção cuidadosa, checar `ok`.

**Typed context key** (em `middleware`):

```go
type ctxKey string
const CtxUsuarioID ctxKey = "usuario_id"
```

Evita colisão com outras libs que usem a mesma string como key.

#### Code review

| ✅ Bom | ⚠️ Observação |
|--------|----------------|
| Logout server-side (não só apagar token no client) | Access JWT continua válido até expirar (TTL curto mitiga; denylist seria overkill no MVP) |
| Handler e use case testados | — |

**Pergunta para a turma:** logout invalida o access token imediatamente? Por quê (não)?

---

### PR #8 — `feat(artigos): CRUD de artigos com publicação`

🔗 https://github.com/Thiago-Tertuliano/MVP-Proj/pull/8  
Branch: `feat/artigos-crud`

#### Escopo

| Camada | Peças |
|--------|--------|
| Domain | `Artigo`, `Slug`, `ArtigoStatus` |
| Application | Criar, Obter, Listar, Atualizar, Publicar |
| Infra | Migration `0002`, `ArtigoRepoPG` |
| HTTP | Handler + rotas públicas/protegidas |

#### Endpoints

| Método | Rota | Auth | Comportamento |
|--------|------|------|---------------|
| GET | `/artigos` | Não | Lista publicados |
| GET | `/artigos/{slug}` | Não | Detalhe publicado |
| POST | `/artigos` | Sim | Cria rascunho |
| PUT | `/artigos/{id}` | Sim | Atualiza (só autor) |
| POST | `/artigos/{id}/publicar` | Sim | Publica (só autor) |

#### Domínio — regras no aggregate

```go
func (a *Artigo) Publicar() error {
    if a.status != Revisao && a.status != Rascunho {
        return errors.ErrInvalidState(...)
    }
    if conteudo vazio {
        return errors.ErrInvalidArgument(...)
    }
    a.status = Publicado
    now := time.Now().UTC()
    a.publicadoEm = &now  // *time.Time — nil quando não publicado
    a.Touch()
    return nil
}
```

**Sintaxe:** `*time.Time` diferencia “sem data” (`nil`) de “zero value” (`0001-01-01`). Em JSON vira `publicado_em` omitido ou unix timestamp.

#### VO `Slug`

```go
func NewSlug(raw string) (Slug, error) {
    normalized := normalizeSlug(raw) // lower, acentos, hífens
    if !slugPattern.MatchString(normalized) {
        return Slug{}, errors.ErrInvalidArgument(...)
    }
    return Slug{value: normalized}, nil
}
```

Criar artigo sem `slug` no body → deriva do título (`Arquitetura Hexagonal` → `arquitetura-hexagonal`).

#### Conteúdo JSONB

```go
conteudo json.RawMessage  // []byte que é JSON válido
```

`encoding/json.RawMessage` adia o parse — o domínio não precisa conhecer o schema dos blocks (Notion-style). Postgres guarda como `JSONB`.

#### Autorização no use case

```go
if !artigo.EhAutor(autorUUID) {
    return nil, errors.ErrForbidden("apenas o autor pode editar o artigo", "...", nil)
}
```

Handler só extrai ID do token; a regra “só o autor” fica na aplicação/domínio.

#### Chi — path params

```go
slug := chi.URLParam(r, "slug")
id := chi.URLParam(r, "id")
```

Rotas registradas no composition root (`router.New`).

#### Code review

| ✅ Bom | ⚠️ Observação |
|--------|----------------|
| Aggregate rico (não é anêmico) | Publicar aceita rascunho direto (README original pedia só revisão — simplificação de MVP) |
| Listagem só de publicados | Sem paginação cursor; limit/offset simples |
| Upsert no Save (`ON CONFLICT`) | Embedding/pgvector ainda não — próximo sprint |

**Pergunta para a turma:** por que `GET` por slug é público e `PUT` por id exige auth?

---

### PR #9 — `feat(trilhas): CRUD de trilhas e módulos`

🔗 https://github.com/Thiago-Tertuliano/MVP-Proj/pull/9  
Branch: `feat/trilhas-modulos`

#### Escopo

Aggregate `Trilha` contendo entidades filhas `Modulo` (não são aggregate roots).

```go
type Trilha struct {
    kernel.BaseEntity
    slug      valueobject.Slug
    publicada bool
    modulos   []*Modulo
}

func (t *Trilha) AdicionarModulo(...) (*Modulo, error) { ... }
func (t *Trilha) Publicar() error {
    if len(t.modulos) == 0 {
        return errors.ErrInvalidState("trilha sem módulos não pode ser publicada", ...)
    }
    t.publicada = true
    return nil
}
```

#### Endpoints

| Método | Rota | Auth | Função |
|--------|------|------|--------|
| GET | `/trilhas` | Não | Lista publicadas |
| GET | `/trilhas/{slug}` | Não | Detalhe + módulos |
| POST | `/trilhas` | Sim | Cria |
| POST | `/trilhas/{id}/modulos` | Sim | Adiciona módulo |
| POST | `/trilhas/{id}/publicar` | Sim | Publica |

#### Persistência transacional

No `TrilhaRepoPG.Save`:

```go
tx, err := r.pool.Begin(ctx)
defer tx.Rollback(ctx)          // no-op se já Commitou

// UPSERT trilha
// DELETE FROM modulos WHERE trilha_id = ...
// INSERT de cada módulo

return tx.Commit(ctx)
```

**Por quê?** Módulos vivem dentro do aggregate. Regravar = substituir a coleção inteira na mesma transação (consistência).

#### Sintaxe em destaque

- `defer tx.Rollback(ctx)` é padrão seguro com pgx.
- Slice `[]*Modulo` — ponteiros para evitar cópia grande e permitir mutação.

#### Code review

| ✅ Bom | ⚠️ Observação |
|--------|----------------|
| Invariante “sem módulo → não publica” | Ainda **não** liga artigo ↔ trilha/módulo (FK fica para depois) |
| VO `Slug` incluído nesta PR | Mesmo arquivo existe no PR #8 → conflito trivial no merge |
| — | `router.go` e `Makefile` conflitam com #8 — unir rotas e migrations 0002+0003 |

**Pergunta para a turma:** módulo é aggregate root? Por que não?

---

## 6. Fluxo ponta a ponta (desenhar no quadro)

```
Cliente
  │  POST /api/v1/artigos  + Authorization: Bearer …
  ▼
Middleware Auth  →  valida JWT → context(usuario_id)
  ▼
ArtigoHandler.Criar  →  decode JSON + validator
  ▼
CriarArtigo.Execute  →  NewSlug, NovoArtigo, repo.Save
  ▼
ArtigoRepoPG  →  INSERT JSONB no Postgres
  ▼
201 + ArtigoResponse (DTO)
```

Mesmo desenho vale para auth, trilhas, etc. **Vertical slice** = cada feature atravessa todas as camadas.

---

## 7. Tabela-resumo: conceitos Go ↔ lugar no código

| Conceito Go | Onde aparece no projeto |
|-------------|-------------------------|
| Package + `internal/` | Toda a árvore `backend-platform` |
| Struct + métodos | `Usuario`, `Artigo`, `Trilha` |
| Embedding | `kernel.BaseEntity` |
| Interface implícita | `TokenGerador`, `ArtigoRepository`, fakes nos testes |
| Ponteiro vs valor | Receivers `*Artigo`, `*time.Time` opcional |
| `error` + `errors.As` | `DomainError` nos handlers |
| `context.Context` | Use cases e repos |
| Middleware chain | Chi `r.Use(...)` |
| Goroutine + channel | Graceful shutdown no `main` |
| Table-driven / mocks | `*_test.go` + `mocks_test.go` |
| `map[K]struct{}` | Allowlist CORS |
| `json.RawMessage` | Conteúdo do artigo |
| Transação (`Begin/Commit`) | Save de trilha + módulos |

---

## 8. Roteiro sugerido da aula (45–60 min)

1. **5 min** — Problema de negócio e stack (Go + Postgres + Chi).
2. **10 min** — Camadas DDD no filesystem; Dependency Rule.
3. **10 min** — Auth Sprint 1: VO Email, ports, JWT, anti-enumeração.
4. **5 min** — PR #5 e #7: consistência de sessão (persistir + logout).
5. **5 min** — PR #6: middleware e CORS (browser).
6. **10 min** — PR #8: aggregate Artigo, slug, JSONB, Forbidden.
7. **10 min** — PR #9: aggregate com filhos, transação, invariantes.
8. **5 min** — Exercício: “onde colocariam rate limit? e vínculo artigo–trilha?”

---

## 9. Exercícios para os alunos

1. Onde a senha em texto puro **não** pode aparecer? (Resposta: domínio só vê `SenhaHash`; plain só no use case → port hasher.)
2. Por que `ReconstruirUsuario` não valida de novo? (Dado já persistido; validação na escrita.)
3. Implementar mentalmente: `DELETE /artigos/{id}` — em qual camada a regra “só autor ou admin”?
4. No merge de #8+#9, o que quebra se esquecer de registrar rotas no `router`? (Compila, mas 404.)

---

## 10. O que ainda não está feito (honestidade na aula)

- Vínculo `artigo.trilha_id` / `modulo_id`
- Embeddings + `pgvector` (busca semântica)
- Cookies HttpOnly (hoje tokens no JSON body)
- Frontend Next.js
- Rate limit em `/auth/*`
- sqlc (README menciona; código usa `pgx` direto)

Isso reforça: **documentação de arquitetura ≠ código entregue** — e está tudo bem no MVP, desde que o gap seja consciente.

---

## 11. Links rápidos dos PRs

| PR | Título | URL |
|----|--------|-----|
| #5 | Persiste refresh no registro | https://github.com/Thiago-Tertuliano/MVP-Proj/pull/5 |
| #6 | Middleware CORS | https://github.com/Thiago-Tertuliano/MVP-Proj/pull/6 |
| #7 | Logout | https://github.com/Thiago-Tertuliano/MVP-Proj/pull/7 |
| #8 | CRUD Artigos | https://github.com/Thiago-Tertuliano/MVP-Proj/pull/8 |
| #9 | CRUD Trilhas/Módulos | https://github.com/Thiago-Tertuliano/MVP-Proj/pull/9 |

---

*Documento gerado para apoio à aula sobre o backend Go do Estudos Platform — Axellion / MVP Bruno & Thiago.*

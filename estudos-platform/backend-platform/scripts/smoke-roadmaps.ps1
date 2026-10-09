# Smoke E2E de roadmaps gamificados contra uma API local já em execução.
# Pré-requisitos: migrations aplicadas, seed (autor.seed@estudos.local) e
# EDITOR_EMAILS=autor.seed@estudos.local no ambiente da API.
# Uso: .\scripts\smoke-roadmaps.ps1 [-Base http://localhost:8080/api/v1]
param([string]$Base = 'http://localhost:8080/api/v1')
$ErrorActionPreference = 'Stop'
$base = $Base

function Chamar($sessao, $metodo, $caminho, $corpo) {
  $params = @{ Uri = "$base$caminho"; Method = $metodo; WebSession = $sessao; ContentType = 'application/json; charset=utf-8' }
  if ($null -ne $corpo) { $params.Body = [System.Text.Encoding]::UTF8.GetBytes(($corpo | ConvertTo-Json -Depth 12)) }
  try {
    $r = Invoke-WebRequest @params -UseBasicParsing
    return @{ status = [int]$r.StatusCode; body = ($r.Content | ConvertFrom-Json) }
  } catch {
    $resp = $_.Exception.Response
    $txt = ''
    if ($resp) { $sr = New-Object IO.StreamReader($resp.GetResponseStream()); $txt = $sr.ReadToEnd() }
    $b = $null; try { $b = $txt | ConvertFrom-Json } catch {}
    return @{ status = [int]$resp.StatusCode; body = $b }
  }
}
function Ok($cond, $msg) { if ($cond) { Write-Host "OK   $msg" } else { Write-Host "FAIL $msg"; $script:falhas++ } }
$script:falhas = 0

# --- sessões
$ed = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$al = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$r = Chamar $ed 'POST' '/auth/login' @{ email = 'autor.seed@estudos.local'; senha = 'senha1234' }
Ok ($r.status -eq 200) "login editor ($($r.status))"
$email = "aluno$([guid]::NewGuid().ToString('N').Substring(0,8))@teste.com"
$r = Chamar $al 'POST' '/auth/registrar' @{ nome = 'Aluno Teste'; email = $email; senha = 'senha1234' }
Ok ($r.status -eq 201 -or $r.status -eq 200) "registrar aluno ($($r.status))"

$r = Chamar $ed 'GET' '/auth/me' $null
Ok ($r.body.papel -eq 'editor') "me do editor traz papel=editor ($($r.body.papel))"
$r = Chamar $al 'GET' '/auth/me' $null
Ok ($r.body.papel -eq 'aluno') "me do aluno traz papel=aluno ($($r.body.papel))"

# --- aluno não edita
$r = Chamar $al 'POST' '/roadmaps' @{ titulo = 'Invasor' }
Ok ($r.status -eq 403) "aluno recebe 403 ao criar roadmap ($($r.status))"
$r = Chamar $al 'GET' '/roadmaps/meus' $null
Ok ($r.status -eq 403) "aluno recebe 403 em /roadmaps/meus ($($r.status))"
$r = Chamar $null 'GET' '/roadmaps/meus' $null
Ok ($r.status -eq 401) "anônimo recebe 401 ($($r.status))"

# --- artigo publicado para vincular
$r = Chamar $ed 'GET' '/artigos' $null
$artigo = $r.body.itens | Select-Object -First 1
if (-not $artigo) { $artigo = $r.body.artigos | Select-Object -First 1 }
Ok ($null -ne $artigo) "existe artigo publicado para vincular ($($artigo.slug))"

# --- editor cria e salva
$titulo = "Go E2E $([guid]::NewGuid().ToString('N').Substring(0,6))"
$r = Chamar $ed 'POST' '/roadmaps' @{ titulo = $titulo; descricao = 'Mapa de teste'; icone = 'map' }
Ok ($r.status -eq 201) "editor cria rascunho ($($r.status))"
$slug = $r.body.slug
$a = [guid]::NewGuid().ToString(); $b = [guid]::NewGuid().ToString(); $c = [guid]::NewGuid().ToString()
$grafo = @{
  titulo = $titulo; descricao = 'Mapa de teste'; icone = 'map'
  nos = @(
    @{ id = $a; tipo = 'topico'; titulo = 'Variáveis'; pos_x = 0; pos_y = 0; xp = 10; conclusao = 'manual' },
    @{ id = $b; tipo = 'topico'; titulo = 'Artigo vinculado'; pos_x = 0; pos_y = 150; xp = 20; conclusao = 'artigo'; artigo_id = $artigo.id },
    @{ id = $c; tipo = 'chefe'; titulo = 'Chefe'; pos_x = 0; pos_y = 300; xp = 100; conclusao = 'quiz'; quiz = @{ questoes = @(
      @{ id = 'q1'; enunciado = '1+1?'; correta = 'b'; explicacao = 'Dois.'; opcoes = @(@{ id = 'a'; texto = '1' }, @{ id = 'b'; texto = '2' }) },
      @{ id = 'q2'; enunciado = '2+2?'; correta = 'a'; opcoes = @(@{ id = 'a'; texto = '4' }, @{ id = 'b'; texto = '5' }) }
    ) } }
  )
  arestas = @(@{ origem = $a; destino = $b; tipo = 'requer' }, @{ origem = $b; destino = $c; tipo = 'requer' })
}
$r = Chamar $ed 'PUT' "/roadmaps/$slug" $grafo
Ok ($r.status -eq 200 -and $r.body.nos.Count -eq 3) "editor salva grafo ($($r.status))"

# rascunho invisível ao público
$r = Chamar $null 'GET' "/roadmaps/$slug" $null
Ok ($r.status -eq 404) "rascunho é 404 para o público ($($r.status))"
$r = Chamar $al 'GET' "/roadmaps/$slug/progresso" $null
Ok ($r.status -eq 404) "rascunho é 404 no progresso ($($r.status))"

# ciclo rejeitado
$ciclo = $grafo.Clone(); $ciclo.arestas = @(@{ origem = $a; destino = $b }, @{ origem = $b; destino = $a })
$r = Chamar $ed 'PUT' "/roadmaps/$slug" $ciclo
Ok ($r.status -eq 400) "ciclo é rejeitado com 400 ($($r.status)) $($r.body.erro)"

# aluno não edita roadmap alheio nem como editor diferente: aluno -> 403
$r = Chamar $al 'PUT' "/roadmaps/$slug" $grafo
Ok ($r.status -eq 403) "aluno não salva ($($r.status))"

$r = Chamar $ed 'POST' "/roadmaps/$slug/publicar" $null
Ok ($r.status -eq 200 -and $r.body.publicado) "editor publica ($($r.status))"

# --- público
$r = Chamar $null 'GET' '/roadmaps' $null
Ok (($r.body.itens | Where-Object { $_.slug -eq $slug }).total_nos -eq 3) "catálogo lista com 3 nós"
$r = Chamar $null 'GET' "/roadmaps/$slug" $null
$chefe = $r.body.nos | Where-Object { $_.tipo -eq 'chefe' }
Ok ($r.status -eq 200) "público lê roadmap publicado"
$json = ($r.body | ConvertTo-Json -Depth 12)
Ok (-not ($json -match '"correta"') -and -not ($json -match 'Dois\.')) "gabarito NÃO vaza para o público"
Ok ($chefe.quiz.questoes.Count -eq 2) "perguntas do chefe presentes"
$noArtigo = $r.body.nos | Where-Object { $_.tipo -eq 'topico' -and $_.conclusao -eq 'artigo' }
Ok ($noArtigo.artigo_slug -eq $artigo.slug) "nó traz artigo_slug ($($noArtigo.artigo_slug))"

# --- aluno joga
$r = Chamar $al 'GET' "/roadmaps/$slug/progresso" $null
Ok ($r.body.estados.$a -eq 'disponivel' -and $r.body.estados.$b -eq 'bloqueado' -and $r.body.estados.$c -eq 'bloqueado') "estados iniciais"
$r = Chamar $al 'POST' "/roadmaps/$slug/nos/$b/concluir" $null
Ok ($r.status -eq 400 -or $r.status -eq 409) "nó bloqueado e por artigo não conclui manualmente ($($r.status))"
$r = Chamar $al 'POST' "/roadmaps/$slug/nos/$c/quiz" @{ respostas = @{ q1 = 'b'; q2 = 'a' } }
Ok ($r.status -eq 409) "chefe bloqueado => 409 ao responder quiz ($($r.status))"
$r = Chamar $al 'POST' "/roadmaps/$slug/nos/$a/concluir" $null
Ok ($r.status -eq 200 -and $r.body.resultado.xp_ganho -eq 10 -and $r.body.estados.$b -eq 'disponivel') "concluir A: +10xp e libera B"
Ok ($r.body.resultado.streak -eq 1) "streak = 1"
Ok (($r.body.resultado.conquistas_novas | Where-Object { $_.codigo -eq 'primeiro-passo' }) -ne $null) "conquista primeiro-passo"
$r = Chamar $al 'POST' "/roadmaps/$slug/nos/$a/concluir" $null
Ok ($r.body.ja_concluido -and $r.body.resultado.xp_ganho -eq 0) "repetir não dá XP"
$r = Chamar $al 'POST' "/roadmaps/$slug/nos/$b/concluir" $null
Ok ($r.status -eq 400) "nó por artigo não aceita 'dominei' ($($r.status))"

# lê o artigo -> conclui nó B
$r = Chamar $al 'PUT' "/progresso/artigos/$($artigo.id)" @{ concluido = $true }
Ok ($r.status -eq 200 -and $r.body.gamificacao.xp_ganho -eq 20) "ler artigo vinculado: +20xp via hook ($($r.body.gamificacao.xp_ganho))"
$r = Chamar $al 'PUT' "/progresso/artigos/$($artigo.id)" @{ concluido = $true }
Ok ($r.body.gamificacao.xp_ganho -eq 0) "reler não duplica XP"

# quiz
$r = Chamar $al 'POST' "/roadmaps/$slug/nos/$c/quiz" @{ respostas = @{ q1 = 'a'; q2 = 'b' } }
Ok ($r.status -eq 200 -and -not $r.body.aprovado -and $null -eq $r.body.correcao -and $null -eq $r.body.conclusao) "quiz reprovado não conclui nem vaza correção"
$r = Chamar $al 'POST' "/roadmaps/$slug/nos/$c/quiz" @{ respostas = @{ q1 = 'b'; q2 = 'a' } }
Ok ($r.body.aprovado -and $r.body.conclusao.resultado.xp_ganho -ge 100) "quiz aprovado: +100xp (+bônus)"
Ok ($r.body.conclusao.resultado.subiu_de_nivel -eq $true) "subiu de nível"
Ok (($r.body.conclusao.resultado.conquistas_novas | Where-Object { $_.codigo -eq 'cacador-de-chefes' }) -ne $null) "conquista caçador de chefes"
Ok (($r.body.conclusao.resultado.conquistas_novas | Where-Object { $_.codigo -eq 'cartografo' }) -ne $null) "conquista cartógrafo (roadmap completo)"
Ok ($r.body.correcao.Count -eq 2) "correção enviada após aprovação"

$r = Chamar $al 'GET' '/gamificacao/me' $null
Ok ($r.body.xp_total -eq (10 + 20 + 100 + 50) -and $r.body.nivel -ge 2) "gamificacao/me: xp=$($r.body.xp_total) nivel=$($r.body.nivel)"
$r = Chamar $al 'GET' '/gamificacao/roadmaps' $null
$it = $r.body.itens | Where-Object { $_.slug -eq $slug }
Ok ($it.concluidos -eq 3 -and $it.total -eq 3) "progresso no catálogo 3/3"

# despublicar
$r = Chamar $ed 'POST' "/roadmaps/$slug/despublicar" $null
Ok ($r.status -eq 200 -and -not $r.body.publicado) "despublicar"
$r = Chamar $null 'GET' "/roadmaps/$slug" $null
Ok ($r.status -eq 404) "volta a 404 para o público"

Write-Host "`nFalhas: $script:falhas"

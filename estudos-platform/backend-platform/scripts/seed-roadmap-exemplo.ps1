# Cria (ou atualiza) e PUBLICA um roadmap de exemplo "Go do zero ao deploy" para validar o visual.
# Pré-requisitos: API local em execução, seed (autor.seed@estudos.local) e
# EDITOR_EMAILS=autor.seed@estudos.local no ambiente da API. Idempotente: roda quantas vezes quiser.
# Uso: .\scripts\seed-roadmap-exemplo.ps1 [-Base http://localhost:8080/api/v1]
param([string]$Base = 'http://localhost:8080/api/v1')
$ErrorActionPreference = 'Stop'

function Chamar($sessao, $metodo, $caminho, $corpo) {
  $params = @{ Uri = "$Base$caminho"; Method = $metodo; WebSession = $sessao; ContentType = 'application/json; charset=utf-8'; UseBasicParsing = $true }
  if ($null -ne $corpo) { $params.Body = [System.Text.Encoding]::UTF8.GetBytes(($corpo | ConvertTo-Json -Depth 12)) }
  $r = Invoke-WebRequest @params
  if ($r.Content) { return $r.Content | ConvertFrom-Json }
}

$ed = New-Object Microsoft.PowerShell.Commands.WebRequestSession
Chamar $ed 'POST' '/auth/login' @{ email = 'autor.seed@estudos.local'; senha = 'senha1234' } | Out-Null

# IDs estáveis: reexecutar o seed não cria nós novos.
function Id($n) { return '00000000-0000-4000-8000-{0:D12}' -f $n }

# Vincula até 2 artigos publicados (se houver) aos nós "artigo".
$lista = Chamar $ed 'GET' '/artigos' $null
$artigos = @(); if ($lista.itens) { $artigos = @($lista.itens) } elseif ($lista.artigos) { $artigos = @($lista.artigos) }
function NoConteudo($id, $titulo, $descricao, $xp, $x, $y, $indiceArtigo) {
  if ($artigos.Count -gt $indiceArtigo) {
    return @{ id = (Id $id); tipo = 'topico'; titulo = $titulo; descricao = $descricao; xp = $xp; pos_x = $x; pos_y = $y; conclusao = 'artigo'; artigo_id = $artigos[$indiceArtigo].id }
  }
  return @{ id = (Id $id); tipo = 'topico'; titulo = $titulo; descricao = $descricao; xp = $xp; pos_x = $x; pos_y = $y; conclusao = 'manual' }
}
function NoManual($id, $titulo, $descricao, $xp, $x, $y, $tipo = 'topico') {
  return @{ id = (Id $id); tipo = $tipo; titulo = $titulo; descricao = $descricao; xp = $xp; pos_x = $x; pos_y = $y; conclusao = 'manual' }
}
function Opcoes($a, $b, $c) { return @(@{ id = 'a'; texto = $a }, @{ id = 'b'; texto = $b }, @{ id = 'c'; texto = $c }) }

$nos = @(
  (NoManual 1 'Instalar Go e rodar o primeiro programa' 'Instale o toolchain, crie um módulo e imprima "olá, mundo".' 10 260 0),
  (NoConteudo 2 'Variáveis e tipos' 'Declaração com var e :=, tipos básicos e valores zero.' 15 0 150 0),
  (NoManual 3 'Funções' 'Múltiplos retornos, funções variádicas e defer.' 15 260 150),
  (NoManual 4 'Controle de fluxo' 'if, for e switch: a linguagem só tem um laço.' 15 520 150),
  @{ id = (Id 5); tipo = 'chefe'; titulo = 'Chefe: Fundamentos'; descricao = 'Prove que domina o básico para seguir adiante.'; xp = 120; pos_x = 260; pos_y = 300; conclusao = 'quiz'
     quiz = @{ questoes = @(
       @{ id = 'q1'; enunciado = 'Qual forma declara uma variável dentro de uma função?'; correta = 'b'; explicacao = 'O operador := declara e infere o tipo.'; opcoes = (Opcoes 'let x = 1' 'x := 1' 'int x = 1') },
       @{ id = 'q2'; enunciado = 'Uma função em Go pode retornar...'; correta = 'c'; explicacao = 'Múltiplos retornos são idiomáticos (valor, erro).'; opcoes = (Opcoes 'Apenas um valor' 'Nenhum valor' 'Vários valores') },
       @{ id = 'q3'; enunciado = 'Qual laço existe em Go?'; correta = 'a'; explicacao = 'Só existe for; ele cobre while e do-while.'; opcoes = (Opcoes 'for' 'while' 'foreach') }
     ) } },
  (NoConteudo 6 'Structs e métodos' 'Modelagem de dados sem classes.' 20 0 450 1),
  (NoManual 7 'Pacotes e módulos' 'Organização do código e go.mod.' 20 520 450 'marco'),
  (NoManual 8 'Interfaces' 'Polimorfismo implícito: satisfaça pelo comportamento.' 25 260 600),
  (NoManual 9 'Testes automatizados' 'testing, tabelas de casos e cobertura (opcional, mas recomendado).' 15 0 750),
  (NoManual 10 'Concorrência: goroutines e channels' 'O diferencial da linguagem.' 30 520 750),
  @{ id = (Id 11); tipo = 'chefe'; titulo = 'Chefe final: Go em produção'; descricao = 'Concorrência, interfaces e boas práticas.'; xp = 200; pos_x = 260; pos_y = 900; conclusao = 'quiz'
     quiz = @{ questoes = @(
       @{ id = 'q1'; enunciado = 'Como duas goroutines se comunicam de forma idiomática?'; correta = 'a'; explicacao = '"Não comunique compartilhando memória; compartilhe memória comunicando."'; opcoes = (Opcoes 'Channels' 'Variáveis globais' 'Arquivos temporários') },
       @{ id = 'q2'; enunciado = 'Uma interface em Go é satisfeita...'; correta = 'b'; explicacao = 'Não há palavra "implements": basta ter os métodos.'; opcoes = (Opcoes 'Declarando implements' 'Implicitamente, pelos métodos' 'Só com herança') }
     ) } }
)
$nos[6].tipo = 'marco'

function Aresta($o, $d, $t = 'requer') { return @{ origem = (Id $o); destino = (Id $d); tipo = $t } }
$arestas = @(
  (Aresta 1 2), (Aresta 1 3), (Aresta 1 4),
  (Aresta 2 5), (Aresta 3 5), (Aresta 4 5),
  (Aresta 5 6), (Aresta 5 7),
  (Aresta 6 8), (Aresta 6 9 'opcional'),
  (Aresta 8 10), (Aresta 7 10),
  (Aresta 10 11), (Aresta 9 11 'opcional')
)

$titulo = 'Go do zero ao deploy'
$meus = Chamar $ed 'GET' '/roadmaps/meus' $null
$itens = @(); if ($meus.itens) { $itens = @($meus.itens) }
$existente = $itens | Where-Object { $_.titulo -eq $titulo } | Select-Object -First 1
if ($existente) { $slug = $existente.slug; Write-Host "Atualizando roadmap existente: $slug" }
else {
  $criado = Chamar $ed 'POST' '/roadmaps' @{ titulo = $titulo; descricao = 'Do primeiro programa à concorrência: derrote dois chefes e domine Go.'; icone = 'rocket' }
  $slug = $criado.slug; Write-Host "Criado: $slug"
}

Chamar $ed 'PUT' "/roadmaps/$slug" @{ titulo = $titulo; descricao = 'Do primeiro programa à concorrência: derrote dois chefes e domine Go.'; icone = 'rocket'; nos = $nos; arestas = $arestas } | Out-Null
Chamar $ed 'POST' "/roadmaps/$slug/publicar" $null | Out-Null
Write-Host "Publicado: http://localhost:3100/roadmaps/$slug"

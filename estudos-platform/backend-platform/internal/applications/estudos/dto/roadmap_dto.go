package dto

// ---- Quiz ----

type OpcaoQuizDTO struct {
	ID    string `json:"id"`
	Texto string `json:"texto"`
}

// QuestaoQuizDTO: Correta/Explicacao só trafegam para o editor (omitidos para o aluno).
type QuestaoQuizDTO struct {
	ID         string         `json:"id"`
	Enunciado  string         `json:"enunciado"`
	Opcoes     []OpcaoQuizDTO `json:"opcoes"`
	Correta    string         `json:"correta,omitempty"`
	Explicacao string         `json:"explicacao,omitempty"`
}

type QuizDTO struct {
	Questoes []QuestaoQuizDTO `json:"questoes"`
}

// ---- Grafo ----

type RoadmapNoDTO struct {
	ID           string   `json:"id"`
	Tipo         string   `json:"tipo"`
	Titulo       string   `json:"titulo"`
	Descricao    string   `json:"descricao,omitempty"`
	ArtigoID     string   `json:"artigo_id,omitempty"`
	ArtigoSlug   string   `json:"artigo_slug,omitempty"`
	ArtigoTitulo string   `json:"artigo_titulo,omitempty"`
	PosX         float64  `json:"pos_x"`
	PosY         float64  `json:"pos_y"`
	XP           int      `json:"xp"`
	Conclusao    string   `json:"conclusao"`
	Quiz         *QuizDTO `json:"quiz,omitempty"`
}

type RoadmapArestaDTO struct {
	Origem  string `json:"origem"`
	Destino string `json:"destino"`
	Tipo    string `json:"tipo"`
}

type RoadmapResponse struct {
	ID           string             `json:"id"`
	Slug         string             `json:"slug"`
	Titulo       string             `json:"titulo"`
	Descricao    string             `json:"descricao"`
	Icone        string             `json:"icone"`
	Publicado    bool               `json:"publicado"`
	Nos          []RoadmapNoDTO     `json:"nos"`
	Arestas      []RoadmapArestaDTO `json:"arestas"`
	XPTotal      int                `json:"xp_total"`
	AtualizadoEm int64              `json:"atualizado_em"`
}

type RoadmapResumoResponse struct {
	ID           string `json:"id"`
	Slug         string `json:"slug"`
	Titulo       string `json:"titulo"`
	Descricao    string `json:"descricao"`
	Icone        string `json:"icone"`
	Publicado    bool   `json:"publicado"`
	TotalNos     int    `json:"total_nos"`
	TotalChefes  int    `json:"total_chefes"`
	XPTotal      int    `json:"xp_total"`
	AtualizadoEm int64  `json:"atualizado_em"`
}

type ListarRoadmapsResponse struct {
	Itens []RoadmapResumoResponse `json:"itens"`
}

// ---- Edição ----

type CriarRoadmapRequest struct {
	Titulo    string `json:"titulo"`
	Descricao string `json:"descricao"`
	Icone     string `json:"icone"`
}

type SalvarRoadmapRequest struct {
	Titulo    string             `json:"titulo"`
	Descricao string             `json:"descricao"`
	Icone     string             `json:"icone"`
	Nos       []RoadmapNoDTO     `json:"nos"`
	Arestas   []RoadmapArestaDTO `json:"arestas"`
}

// ---- Progresso do aluno ----

type RoadmapProgressoResponse struct {
	Slug       string            `json:"slug"`
	Estados    map[string]string `json:"estados"`
	Concluidos int               `json:"concluidos"`
	Total      int               `json:"total"`
	Percentual float64           `json:"percentual"`
	XPGanho    int               `json:"xp_ganho"`
	XPTotal    int               `json:"xp_total"`
	Completo   bool              `json:"completo"`
}

type RoadmapProgressoItem struct {
	Slug       string  `json:"slug"`
	Concluidos int     `json:"concluidos"`
	Total      int     `json:"total"`
	Percentual float64 `json:"percentual"`
}

type ListarProgressoRoadmapsResponse struct {
	Itens []RoadmapProgressoItem `json:"itens"`
}

// ---- Gamificação ----

type ConquistaDTO struct {
	Codigo      string `json:"codigo"`
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	Icone       string `json:"icone"`
	Conquistada bool   `json:"conquistada"`
	Em          *int64 `json:"em,omitempty"`
}

// ResultadoGamificacao resume o que mudou após uma ação do aluno (para toasts/animações).
type ResultadoGamificacao struct {
	XPGanho           int            `json:"xp_ganho"`
	XPTotal           int            `json:"xp_total"`
	Nivel             int            `json:"nivel"`
	NivelAnterior     int            `json:"nivel_anterior"`
	SubiuDeNivel      bool           `json:"subiu_de_nivel"`
	Streak            int            `json:"streak"`
	NosConcluidos     []string       `json:"nos_concluidos"`
	NosDesbloqueados  []string       `json:"nos_desbloqueados"`
	RoadmapsCompletos []string       `json:"roadmaps_completos"`
	ConquistasNovas   []ConquistaDTO `json:"conquistas_novas"`
}

type ConclusaoNoResponse struct {
	NoID        string               `json:"no_id"`
	JaConcluido bool                 `json:"ja_concluido"`
	Resultado   ResultadoGamificacao `json:"resultado"`
	Estados     map[string]string    `json:"estados"`
	Concluidos  int                  `json:"concluidos"`
	Total       int                  `json:"total"`
}

type ResponderQuizRequest struct {
	// Respostas mapeia id da questão → id da opção escolhida.
	Respostas map[string]string `json:"respostas"`
}

type CorrecaoQuestaoDTO struct {
	QuestaoID  string `json:"questao_id"`
	Correta    bool   `json:"correta"`
	Explicacao string `json:"explicacao,omitempty"`
}

type QuizResultadoResponse struct {
	Aprovado      bool `json:"aprovado"`
	Acertos       int  `json:"acertos"`
	Total         int  `json:"total"`
	NotaMinimaPct int  `json:"nota_minima_pct"`
	// Correcao só é enviada quando o aluno é aprovado (evita descobrir o gabarito por tentativa e erro).
	Correcao  []CorrecaoQuestaoDTO `json:"correcao,omitempty"`
	Conclusao *ConclusaoNoResponse `json:"conclusao,omitempty"`
}

type GamificacaoResponse struct {
	XPTotal       int            `json:"xp_total"`
	Nivel         int            `json:"nivel"`
	XPNoNivel     int            `json:"xp_no_nivel"`
	XPParaProximo int            `json:"xp_para_proximo"`
	StreakAtual   int            `json:"streak_atual"`
	StreakMax     int            `json:"streak_max"`
	AtivoHoje     bool           `json:"ativo_hoje"`
	Conquistas    []ConquistaDTO `json:"conquistas"`
}

package valueobject

// SenhaHash guarda APENAS o hash. O algoritmo de hashing (bcrypt/argon2)
// fica na infraestrutura, injetado via port — assim o domínio permanece puro.
type SenhaHash struct {
	value string
}

func NovoHashSenha(value string) SenhaHash { return SenhaHash{value: value} }

func (s SenhaHash) Value() string  { return s.value }
func (s SenhaHash) IsZero() bool   { return s.value == "" }
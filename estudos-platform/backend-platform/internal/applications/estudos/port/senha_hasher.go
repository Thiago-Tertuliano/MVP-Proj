package port

// SenhaHasher abstrai o algoritmo de hash.
// Permite trocar bcrypt → argon2 na infra sem tocar em nada da aplicação.
type SenhaHasher interface {
	Hash(plain string) (string, error)
	Comparar(hash, plain string) bool
}

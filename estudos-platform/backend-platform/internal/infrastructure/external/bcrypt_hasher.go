package external

import "golang.org/x/crypto/bcrypt"

// BcryptHasher implementa port.SenhaHasher usando bcrypt.
type BcryptHasher struct {
	custo int
}

func NewBcryptHasher(custo int) *BcryptHasher {
	if custo == 0 {
		custo = bcrypt.DefaultCost // 10
	}
	return &BcryptHasher{custo: custo}
}

func (h *BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.custo)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h *BcryptHasher) Comparar(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

package errors

import "fmt"

// Kind classifica o erro para o handler mapear para o HTTP status correto.
type Kind string

const (
	InvalidArgument Kind = "INVALID_ARGUMENT"
	NotFound        Kind = "NOT_FOUND"
	AlreadyExists   Kind = "ALREADY_EXISTS"
	InvalidState    Kind = "INVALID_STATE"
	Unauthorized    Kind = "UNAUTHORIZED"
	Forbidden       Kind = "FORBIDDEN"
	Internal        Kind = "INTERNAL"
)

// DomainError é o erro estruturado da camada de domínio.
type DomainError struct {
	Kind    Kind
	Message string
	Op      string // operação que falhou, para logs/diagnóstico
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Kind, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Kind, e.Message)
}

func (e *DomainError) Unwrap() error { return e.Err }

func New(kind Kind, message, op string, err error) *DomainError {
	return &DomainError{Kind: kind, Message: message, Op: op, Err: err}
}

func ErrInvalidArgument(message, op string, err error) *DomainError {
	return New(InvalidArgument, message, op, err)
}

func ErrNotFound(message, op string, err error) *DomainError {
	return New(NotFound, message, op, err)
}

func ErrAlreadyExists(message, op string, err error) *DomainError {
	return New(AlreadyExists, message, op, err)
}

func ErrInvalidState(message, op string, err error) *DomainError {
	return New(InvalidState, message, op, err)
}

func ErrUnauthorized(message, op string, err error) *DomainError {
	return New(Unauthorized, message, op, err)
}

func ErrInternal(message, op string, err error) *DomainError {
	return New(Internal, message, op, err)
}
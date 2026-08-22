package usecase

import (
	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func parseOptionalUUID(raw *string, op string) (*uuid.UUID, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*raw)
	if err != nil {
		return nil, errors.ErrInvalidArgument("id inválido", op, err)
	}
	return &id, nil
}

func uuidPtrString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}

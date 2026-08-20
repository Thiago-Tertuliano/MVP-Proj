package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	domainErrors "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

const (
	pgUndefinedTable  = "42P01"
	pgUniqueViolation = "23505"
	pgForeignKey      = "23503"
)

func MapPG(err error, op string) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case pgUndefinedTable:
		return domainErrors.ErrInternal(
			"banco sem tabelas — rode .\\scripts\\migrate-up.ps1",
			op,
			err,
		)
	case pgUniqueViolation, pgForeignKey:
		// 23505/23503 já são tratados no repo (Save, Upsert). MapPG não reescreve.
		return err
	default:
		return err
	}
}

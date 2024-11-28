//go:build public

package handlers

import (
	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	"github.com/jackc/pgx/v5"
)

func PublicDb() *db.Queries {
	return db.New(dbConn)
}

func withTransaction(fn func(tx pgx.Tx) error) error {
	tx, err := dbConn.Begin(*dbCtx)
	if err != nil {
		return customErrors.Database.TransactionStartError
	}
	defer func() {
		_ = tx.Rollback(*dbCtx)
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(*dbCtx); err != nil {
		return customErrors.Database.TransactionCommitError
	}

	return nil
}
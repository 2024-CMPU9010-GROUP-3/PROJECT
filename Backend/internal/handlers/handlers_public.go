//go:build public

package handlers

import (
	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
)

func PublicDb() *db.Queries {
	return db.New(dbConn)
}

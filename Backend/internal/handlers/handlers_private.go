//go:build private

package handlers

import (
	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/private"
)

func PrivateDb() *db.Queries {
	return db.New(dbConn)
}

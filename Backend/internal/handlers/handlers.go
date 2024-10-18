package handlers

import (
	"context"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db"
)

type PointsHandler struct{}
type AuthHandler struct{}

var dbCtx *context.Context
var dbConn db.DBConn

// this needs to be refactored into proper dependency injection in the future
func RegisterDatabaseConnection(ctx *context.Context, conn db.DBConn) {
	dbCtx = ctx
	dbConn = conn
}

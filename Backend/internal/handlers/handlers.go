package handlers

import (
	"context"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/db"
)

type PointsHandler struct{}
type AuthHandler struct{}
type LocationHistoryHandler struct {}

var dbCtx *context.Context
var dbConn db.DBConn

// this needs to be refactored into proper dependency injection in the future
func RegisterDatabaseConnection(ctx *context.Context, conn db.DBConn) {
	dbCtx = ctx
	dbConn = conn
}

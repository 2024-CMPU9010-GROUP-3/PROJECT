package handlers

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PointsHandler struct{}
type AuthHandler struct{}

var dbCtx *context.Context
var dbConn *pgx.Conn

// this needs to be refactored into proper dependency injection in the future
func RegisterDatabaseConnection(ctx *context.Context, conn *pgx.Conn) {
	dbCtx = ctx
	dbConn = conn
}

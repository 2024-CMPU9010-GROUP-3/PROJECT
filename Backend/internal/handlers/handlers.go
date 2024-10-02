package handlers

import (
	"context"

	publicDb "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/public"
	"github.com/jackc/pgx/v5"
)

const contentType = "Content-Type"
const applicationJson = "application/json"

type PointsHandler struct{}
type AuthHandler struct{}

var dbCtx *context.Context
var dbConn *pgx.Conn
var dbQueries *publicDb.Queries

// this needs to be refactored into proper dependency injection in the future
func RegisterDatabaseConnection(ctx *context.Context, conn *pgx.Conn) {
	dbCtx = ctx
	dbConn = conn
	dbQueries = publicDb.New(dbConn)
}

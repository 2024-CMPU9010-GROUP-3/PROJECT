package util

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenKey string

type placeholder struct {
	IsPlaceholder bool
	Endpoint      string
}

func Placeholder(endpoint string) *placeholder {
	return &placeholder{true, endpoint}
}

func CheckResponseError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, fmt.Sprintf("response error, %v", err), http.StatusInternalServerError)
	}
}

func GetLatestMigrationVersion(migrationsPath string) (uint, error) {
	var latestVersion uint // use uint to match migrate package

	err := filepath.Walk(migrationsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
					return err
			}
			// only look at up-migration files
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".up.sql") {
					parts := strings.Split(info.Name(), "_")
					if len(parts) > 0 {
							version, err := strconv.Atoi(parts[0]) // migration version is the first component
							if err == nil && uint(version) > latestVersion {
									latestVersion = uint(version)
							}
					}
			}
			return nil // no error occurred
	})

	return latestVersion, err
}

// The following is a workaround for a problem with sqlc, pgx and custom types. 
// Original Author: https://github.com/louisrli Source: https://github.com/sqlc-dev/sqlc/issues/2116
//
// Any custom DB types made with CREATE TYPE need to be registered with pgx.
// https://github.com/kyleconroy/sqlc/issues/2116
// https://stackoverflow.com/questions/75658429/need-to-update-psql-row-of-a-composite-type-in-golang-with-jack-pgx
// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype
func GetCustomDataTypes(ctx context.Context, pool *pgxpool.Pool) ([]*pgtype.Type, error) {
	// Get a single connection just to load type information.
	conn, err := pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
					return nil, err
	}

	dataTypeNames := []string{
					"point_type",
					// An underscore prefix is an array type in pgtypes.
					"_point_type",
	}

	var typesToRegister []*pgtype.Type
	for _, typeName := range dataTypeNames {
					dataType, err := conn.Conn().LoadType(ctx, typeName)
					if err != nil {
									return nil, errors.Database.UnknownDatabaseError.WithCause(err)
					}
					// You need to register only for this connection too, otherwise the array type will look for the register element type.
					conn.Conn().TypeMap().RegisterType(dataType)
					typesToRegister = append(typesToRegister, dataType)
	}
	return typesToRegister, nil
}

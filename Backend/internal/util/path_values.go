package util

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
)

func GetUserIdFromRequest(r *http.Request) (pgtype.UUID, error) {
	var userId pgtype.UUID
	userIdPathParam := r.PathValue("id")
	if err := userId.Scan(userIdPathParam); err != nil {
		return userId, customErrors.Parameter.InvalidUUIDError
	}
	return userId, nil
}

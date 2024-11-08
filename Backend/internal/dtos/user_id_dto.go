//go:build public

package dtos

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserIdDto struct {
	UserId pgtype.UUID `json:"userid"`
}

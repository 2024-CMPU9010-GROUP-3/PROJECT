package dtos

import "github.com/jackc/pgx/v5/pgtype"

type UserLoginResponseDto struct {
	UserId pgtype.UUID `json:"userid"`
	Token  string      `json:"token"`
}

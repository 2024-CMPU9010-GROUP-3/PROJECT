package dtos

import (
	"io"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserIdDto struct {
	UserId pgtype.UUID `json:"userid"`
}

func (self UserIdDto) Decode(r io.Reader) error {
	panic("not implemented") // TODO: Implement
}

func (self UserIdDto) Encode() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (self UserIdDto) Validate() error {
	panic("not implemented") // TODO: Implement
}


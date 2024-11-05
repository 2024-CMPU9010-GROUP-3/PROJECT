package dtos

import (
	"encoding/json"
	"io"

	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateUserDto struct {
	Username       string      `json:"username"`
	Email          string      `json:"email"`
	Password       string      `json:"password"`
	FirstName      string      `json:"firstname"`
	LastName       string      `json:"lastname"`
	ProfilePicture pgtype.Text `json:"profilepicture"`
}

func (self *UpdateUserDto) Decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(&self)
	if err != nil {
		return customErrors.Payload.InvalidPayloadUserError
	}

	return self.Validate()
}

func (self *UpdateUserDto) Validate() error {
	return nil
}

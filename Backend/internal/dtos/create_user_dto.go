package dtos

import (
	"encoding/json"
	"fmt"
	"io"

	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserDto struct {
	Username       string      `json:"username"`
	Email          string      `json:"email"`
	Password       string      `json:"password"`
	FirstName      string      `json:"firstname"`
	LastName       string      `json:"lastname"`
	ProfilePicture pgtype.Text `json:"profilepicture"`
}

func (self *CreateUserDto) Decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(&self)
	if err != nil {
		return customErrors.Payload.InvalidPayloadUserError
	}

	return self.Validate()
}

func (self *CreateUserDto) Validate() error {
	if len(self.Username) == 0 {
		return customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Username is required"))
	}

	if len(self.Email) == 0 {
		return customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Email is required"))
	}

	if len(self.Password) == 0 {
		return customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Password is required"))
	}
	return nil
}

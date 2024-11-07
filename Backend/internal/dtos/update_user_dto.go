//go:build public

package dtos

import (
	"encoding/json"
	"fmt"
	"io"

	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateUserDto struct {
	Username       string      `json:"username"`
	Email          string      `json:"email"`
	Password       *string     `json:"password"` // can be nil
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
	if self.Password != nil {
		if len(*self.Password) == 0 {
			err := customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Password cannot be empty"))
			return err
		}
		if len(*self.Password) > 72 {
			err := customErrors.Payload.PasswordTooLongError
			return err
		}
	}
	if len(self.Username) == 0 {
		err := customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Username is required"))
		return err
	}
	if len(self.Email) == 0 {
		err := customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Email is required"))
		return err
	}
	return nil
}

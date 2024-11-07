//go:build public

package dtos

import (
	"encoding/json"
	"fmt"
	"io"

	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
)

type UserLoginDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (self *UserLoginDto) Decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(&self)
	if err != nil {
		return customErrors.Payload.InvalidPayloadUserError
	}

	return self.Validate()
}

func (self *UserLoginDto) Validate() error {
	if len(self.Password) == 0 {
		err := customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Password is required"))
		return err
	}

	if len(self.Username) == 0 && len(self.Email) == 0 {
		err := customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Username or Email is required"))
		return err
	}

	return nil
}
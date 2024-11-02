package dtos

import "io"

type UserLoginDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (self UserLoginDto) Decode(r io.Reader) error {
	panic("not implemented") // TODO: Implement
}

func (self UserLoginDto) Encode() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (self UserLoginDto) Validate() error {
	panic("not implemented") // TODO: Implement
}


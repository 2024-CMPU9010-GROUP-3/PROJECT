package dtos

import "io"

type CreateUserDto struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	ProfilePicture string `json:"profilepicture"`
}

func (self CreateUserDto) Decode(r io.Reader) error {
	panic("not implemented") // TODO: Implement
}

func (self CreateUserDto) Encode() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (self CreateUserDto) Validate() error {
	panic("not implemented") // TODO: Implement
}


package dtos

import "io"

type DTO interface{
	Decode(r io.Reader) error
	Encode() (string, error)
	Validate() error
}


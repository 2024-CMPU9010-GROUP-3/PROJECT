package dtos

import "io"

type DTO interface{
	Decode(r io.Reader) error
	Validate() error
}


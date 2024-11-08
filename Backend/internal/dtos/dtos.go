package dtos

import (
	"io"

	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
)

type DTO interface{
	Decode(r io.Reader) *customErrors.CustomError
	Validate() *customErrors.CustomError
}
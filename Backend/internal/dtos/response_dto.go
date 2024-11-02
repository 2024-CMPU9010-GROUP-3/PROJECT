package dtos

import (
	"io"

	errs "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
)

type ResponseDto struct {
	Error    *errs.CustomError   `json:"error"`
	Response *ResponseContentDto `json:"response"`
}

func (self ResponseDto) Decode(r io.Reader) error {
	panic("not implemented") // TODO: Implement
}

func (self ResponseDto) Encode() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (self ResponseDto) Validate() error {
	panic("not implemented") // TODO: Implement
}


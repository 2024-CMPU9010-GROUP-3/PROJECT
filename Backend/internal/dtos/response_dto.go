package dtos

import (
	errs "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
)

type ResponseDto struct {
	Error    *errs.CustomError   `json:"error"`
	Response *ResponseContentDto `json:"response"`
}

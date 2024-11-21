package dtos

import (
	"encoding/json"
	"fmt"
	"io"

	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
)

type IntIdListDto struct {
	IdList []int64 `json:"idlist"`
}

func (self *IntIdListDto) Decode(r io.Reader) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(self)
	if err != nil {
		return customErrors.Parameter.InvalidIntError
	}

	return self.Validate()
}

func (self *IntIdListDto) Validate() error {
	for id := range self.IdList {
		if id < 0 {
			return customErrors.Parameter.InvalidIntError.WithCause(fmt.Errorf("Id must be greater than or equal to 0"))
		}
	}
	return nil
}

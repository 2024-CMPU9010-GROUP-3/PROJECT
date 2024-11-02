package dtos

import (
	"encoding/json"
	"fmt"
	"io"

	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type CreatePointDto struct {
	Longlat geojson.Geometry `json:"longlat"`
	Type    string           `json:"type"`
	Details any              `json:"details"` // potentially unsafe, but we need to accept any json object here
}

func (self *CreatePointDto) Decode(r io.Reader) *customErrors.CustomError {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(self)
	if err != nil {
		return &customErrors.Payload.InvalidPayloadPointError
	}

	return self.Validate()
}

func (self *CreatePointDto) Validate() *customErrors.CustomError {
	if len(self.Type) == 0 {
		err := customErrors.Parameter.RequiredParameterMissingError.WithCause(fmt.Errorf("Field \"Type\" cannot be empty"))
		return &err
	}
	return nil
}

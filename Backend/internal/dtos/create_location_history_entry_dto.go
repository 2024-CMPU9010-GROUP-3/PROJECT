//go:build public

package dtos

import (
	"encoding/json"
	"fmt"
	"io"

	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type CreateLocationHistoryEntryDto struct {
	Amenitytypes []db.PointType   `json:"amenitytypes"`
	Longlat      geojson.Geometry `json:"longlat"`
	Radius       int32            `json:"radius"`
	DisplayName  pgtype.Text           `json:"displayname"`
}

func (self *CreateLocationHistoryEntryDto) Decode(r io.Reader) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(self)
	if err != nil {
		return customErrors.Payload.InvalidPayloadLocationHistoryEntryError
	}

	return self.Validate()
}

func (self *CreateLocationHistoryEntryDto) Validate() error {
	if self.Radius < 0 {
		err := customErrors.Parameter.InvalidIntError.WithCause(fmt.Errorf("Radius must be greater than or equal to 0"))
		return err
	}
	for _, t := range self.Amenitytypes {
		if !t.IsValid() {
			err := customErrors.Parameter.InvalidPointTypeError.WithCause(fmt.Errorf("Type '%s' is not supported", t))
			return err
		}
	}
	return nil
}

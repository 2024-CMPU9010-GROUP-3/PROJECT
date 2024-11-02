package dtos

import (
	"io"

	"github.com/twpayne/go-geom/encoding/geojson"
)

type CreatePointDto struct {
	Longlat geojson.Geometry `json:"longlat"`
	Type    string           `json:"type"`
	Details any              `json:"details"` // potentially unsafe, but we need to accept any json object here
}

func (self CreatePointDto) Decode(r io.Reader) error {
	panic("not implemented") // TODO: Implement
}

func (self CreatePointDto) Encode() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (self CreatePointDto) Validate() error {
	panic("not implemented") // TODO: Implement
}


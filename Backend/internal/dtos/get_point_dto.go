//go:build public

package dtos

import (
	"io"

	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type PointDto struct {
	Id      int64
	Longlat geojson.Geometry
	Type    db.PointType
}

func (self PointDto) Decode(r io.Reader) error {
	panic("not implemented") // TODO: Implement
}

func (self PointDto) Encode() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (self PointDto) Validate() error {
	panic("not implemented") // TODO: Implement
}


//go:build public

package dtos

import (
	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type GetPointDto struct {
	Id      int64
	Longlat geojson.Geometry
	Type    db.PointType
}

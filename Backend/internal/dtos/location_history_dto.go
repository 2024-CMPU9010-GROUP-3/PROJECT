//go:build public

package dtos

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type LocationHistoryEntryDto struct {
	ID           int64            `json:"id"`
	Datecreated  pgtype.Timestamp `json:"datecreated"`
	Amenitytypes []AmenityTypeWithCount `json:"amenitytypes"`
	Longlat      geojson.Geometry `json:"longlat"`
	Radius       int32            `json:"radius"`
	DisplayName  string           `json:"displayname"`
}

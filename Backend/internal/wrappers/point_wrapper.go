//go:build public

package wrappers

import (
	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/public"
	geos "github.com/twpayne/go-geos/geometry"
)

type PointWrapper struct {
	ID      int64          `json:"id"`
	Longlat *geos.Geometry `json:"longlat"`
	Type    db.PointType   `json:"type"`
}

func FromRow(p db.GetPointsInEnvelopeRow) PointWrapper {
	g := geos.NewGeometry(p.Longlat)
	return PointWrapper{p.ID, g, p.Type}
}

func FromRowList(rows []db.GetPointsInEnvelopeRow) []PointWrapper {
	pointWrappers := []PointWrapper{}
	for _, r := range rows {
		pointWrappers = append(pointWrappers, FromRow(r))
	}
	return pointWrappers
}

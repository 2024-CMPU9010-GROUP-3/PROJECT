//go:build private

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	geos "github.com/twpayne/go-geos"
)

type PointType string

const (
	PointTypePlaceholder1 PointType = "placeholder1"
	PointTypePlaceholder2 PointType = "placeholder2"
)

func (e *PointType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PointType(s)
	case string:
		*e = PointType(s)
	default:
		return fmt.Errorf("unsupported scan type for PointType: %T", src)
	}
	return nil
}

type NullPointType struct {
	PointType PointType `json:"point_type"`
	Valid     bool      `json:"valid"` // Valid is true if PointType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPointType) Scan(value interface{}) error {
	if value == nil {
		ns.PointType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PointType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPointType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PointType), nil
}

type Login struct {
	ID           pgtype.UUID `json:"id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	Passwordhash string      `json:"passwordhash"`
}

type Point struct {
	ID      int64      `json:"id"`
	Longlat *geos.Geom `json:"longlat"`
	Type    PointType  `json:"type"`
	Details []byte     `json:"details"`
}

type UserDetail struct {
	ID             pgtype.UUID      `json:"id"`
	Registerdate   pgtype.Timestamp `json:"registerdate"`
	Firstname      string           `json:"firstname"`
	Lastname       string           `json:"lastname"`
	Profilepicture pgtype.Text      `json:"profilepicture"`
	Lastloggedin   pgtype.Timestamp `json:"lastloggedin"`
}

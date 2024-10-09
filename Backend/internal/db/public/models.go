//go:build public

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	go_geom "github.com/twpayne/go-geom"
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

// Login represents a user's login credentials.
// swagger:model Login
type Login struct {
    // ID of the login session
    // in: string
    ID pgtype.UUID `json:"id"`

    // Username for login
    // required: true
    // in: string
    Username string `json:"username"`

    // Email associated with the login
    // required: true
    // in: string
    Email string `json:"email"`

    // Hashed password for security
    // required: true
    // in: string
    Passwordhash string `json:"passwordhash"`
}

// Point represents a geographical point.
// swagger:model Point
type Point struct {
    // ID of the point
    // in: int64
    ID int64 `json:"id"`

    // Longitude and latitude coordinates
    // in: object
    Longlat *go_geom.Point `json:"longlat"`

    // Type of point (e.g., landmark, checkpoint)
    // required: true
    // in: string
    Type PointType `json:"type"`

    // Additional details about the point
    // in: byte
    Details []byte `json:"details"`
}

// UserDetail contains detailed information about a user.
// swagger:model UserDetail
type UserDetail struct {
    // Unique ID of the user
    // in: string
    ID pgtype.UUID `json:"id"`

    // The date and time the user registered
    // in: string
    Registerdate pgtype.Timestamp `json:"registerdate"`

    // First name of the user
    // required: true
    // in: string
    Firstname string `json:"firstname"`

    // Last name of the user
    // required: true
    // in: string
    Lastname string `json:"lastname"`

    // Profile picture URL or path
    // in: string
    Profilepicture pgtype.Text `json:"profilepicture"`

    // Last time the user logged in
    // in: string
    Lastloggedin pgtype.Timestamp `json:"lastloggedin"`
}
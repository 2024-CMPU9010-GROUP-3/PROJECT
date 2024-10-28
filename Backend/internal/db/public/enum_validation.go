//go:build public

package db

func (p PointType) IsValid() bool {
	switch p {
	case PointTypeParking:
		return true
	}
	return false
}

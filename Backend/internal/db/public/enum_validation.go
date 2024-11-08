//go:build public

package db

func (p PointType) IsValid() bool {
	switch p {
	case PointTypeParking, PointTypeAccessibleParking, PointTypeBikeSharingStation, PointTypeBikeStand, PointTypeCoachParking, PointTypeDrinkingWaterFountain,
		PointTypeLibrary, PointTypeMultistoreyCarParking, PointTypeParkingMeter, PointTypePublicBins, PointTypePublicToilet, PointTypePublicWifiAccessPoint,
		PointTypeUnknown:
		return true
	}
	return false
}

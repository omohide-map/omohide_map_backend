package geo

import "math"

// CalculateDistance はHaversine公式を使用して2点間の距離を計算します（km単位）
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371.0 // 地球の半径（km）

	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLng := (lng2 - lng1) * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// CalculateBoundingBox は中心点と半径から矩形の境界を計算します
func CalculateBoundingBox(lat, lng, radiusKm float64) (minLat, maxLat, minLng, maxLng float64) {
	// 1度あたりおおよそ111km
	latDelta := radiusKm / 111.0
	lngDelta := radiusKm / (111.0 * math.Cos(lat*math.Pi/180.0))

	minLat = lat - latDelta
	maxLat = lat + latDelta
	minLng = lng - lngDelta
	maxLng = lng + lngDelta

	return
}

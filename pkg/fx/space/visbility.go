package space

import (
	"math"
	"time"

	"github.com/Elbujito/2112/pkg/fx/constants"
	"github.com/joshuaferrara/go-satellite"
)

// Calculate LOS dynamically by continuing propagation until the satellite is no longer visible
func ComputeLOS(satrec satellite.Satellite, centerLat float64, centerLong float64, startTime, endTime time.Time, timeStep time.Duration) time.Time {
	for t := startTime.Add(timeStep); t.Before(endTime); t = t.Add(timeStep) {
		year, month, day := t.Date()
		hour, minute, second := t.Clock()

		// Propagate satellite position
		position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)

		// Calculate GST for ECI to LLA conversion
		gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)

		// Convert ECI to Geodetic (lat, lon, alt)
		altitude, _, geoPosition := satellite.ECIToLLA(position, gmst)
		lat, lon := geoPosition.Latitude, geoPosition.Longitude

		// Check elevation
		elevation := CalculateElevation(lat, lon, altitude, centerLat, centerLong)
		if elevation <= 0 {
			// Satellite is no longer visible
			return t
		}
	}

	// Return endTime if LOS is not found
	return endTime
}

// Calculate the elevation angle of the satellite from the tile center
func CalculateElevation(satLat, satLon, satAlt float64, centerLat float64, centerLong float64) float64 {
	// Simplified elevation calculation:
	// Use haversine distance and consider the altitude difference for angular elevation
	dist := HaversineDistance(satLat, satLon, centerLat, centerLong)
	return 90.0 - dist/10.0
}

// Calculate max elevation
func CalculateMaxElevation(satLat, satLon, satAlt float64, centerLat float64, centerLong float64) float64 {
	return 90.0 - HaversineDistance(satLat, satLon, centerLat, centerLong)/10.0
}

// Haversine formula for distance calculation
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = constants.EARTH_CIRCUMFERENCE_METER / 1000
	dLat := DegreesToRadians(lat2 - lat1)
	dLon := DegreesToRadians(lon2 - lon1)

	lat1 = DegreesToRadians(lat1)
	lat2 = DegreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// Convert degrees to radians
func DegreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

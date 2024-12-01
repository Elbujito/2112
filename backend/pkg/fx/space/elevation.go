package space

import (
	"math"

	"github.com/Elbujito/2112/pkg/fx/polygon"
)

// LatLonToCartesian converts latitude, longitude, and altitude to Cartesian coordinates
func LatLonToCartesian(latitude, longitude, altitude float64) (float64, float64, float64) {
	// Earth's radius in kilometers (mean radius)
	const earthRadiusKm = 6371.0

	// Convert latitude and longitude from degrees to radians
	latRad := DegreesToRadians(latitude)
	lonRad := DegreesToRadians(longitude)

	// Calculate the Cartesian coordinates using the spherical to Cartesian transformation
	x := (earthRadiusKm + altitude) * math.Cos(latRad) * math.Cos(lonRad)
	y := (earthRadiusKm + altitude) * math.Cos(latRad) * math.Sin(lonRad)
	z := (earthRadiusKm + altitude) * math.Sin(latRad)

	// Log the calculated Cartesian coordinates for debugging purposes
	// log.Printf("Converted to Cartesian Coordinates: X=%.2f, Y=%.2f, Z=%.2f", x, y, z)

	// Return the Cartesian coordinates (x, y, z)
	return x, y, z
}

// Normalize normalizes a vector
func Normalize(x, y, z float64) (float64, float64, float64) {
	magnitude := math.Sqrt(x*x + y*y + z*z)
	return x / magnitude, y / magnitude, z / magnitude
}

// DotProduct calculates the dot product of two vectors
func DotProduct(x1, y1, z1, x2, y2, z2 float64) float64 {
	return x1*x2 + y1*y2 + z1*z2
}

// CalculateIntegratedElevationFromPoint computes the elevation of a satellite relative to a ground point
func CalculateIntegratedElevationFromPoint(satellitePos polygon.Point, satelliteAltKm float64, groundPoint polygon.Point) float64 {
	// Convert the latitude/longitude of the ground point and satellite to 3D Cartesian coordinates
	groundX, groundY, groundZ := LatLonToCartesian(groundPoint.Latitude, groundPoint.Longitude, 0)
	satX, satY, satZ := LatLonToCartesian(satellitePos.Latitude, satellitePos.Longitude, satelliteAltKm)

	// Compute the vector from the ground to the satellite
	vecX, vecY, vecZ := satX-groundX, satY-groundY, satZ-groundZ

	// Compute the vector from the ground to Earth's center (opposite of the ground point)
	earthVecX, earthVecY, earthVecZ := Normalize(-groundX, -groundY, -groundZ)

	// Normalize the satellite vector
	vecX, vecY, vecZ = Normalize(vecX, vecY, vecZ)

	// Special case handling for direct overhead (tolerance for floating point precision)
	const epsilon = 1e-6
	if math.Abs(vecX-earthVecX) < epsilon && math.Abs(vecY-earthVecY) < epsilon && math.Abs(vecZ-earthVecZ) < epsilon {
		return 90.0
	}

	// Compute the dot product of the vectors
	dotProd := DotProduct(vecX, vecY, vecZ, earthVecX, earthVecY, earthVecZ)

	// Ensure the cosine value is between -1 and 1 for valid arc cosine computation
	if dotProd > 1.0 {
		dotProd = 1.0
	} else if dotProd < -1.0 {
		dotProd = -1.0
	}

	// Compute the elevation angle in radians
	elevationRad := math.Acos(dotProd)

	// Convert the elevation angle from radians to degrees
	elevation := RadiansToDegrees(elevationRad)

	// Elevation cannot be more than 90 degrees
	if elevation > 90 {
		elevation = 90.0
	}

	return elevation
}

// Helper function to convert degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// Helper function to convert radians to degrees
func RadiansToDegrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

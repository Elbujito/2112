package space

import (
	"log"
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
	log.Printf("Converted to Cartesian Coordinates: X=%.2f, Y=%.2f, Z=%.2f", x, y, z)

	// Return the Cartesian coordinates (x, y, z)
	return x, y, z
}

func Normalize(x, y, z float64) (float64, float64, float64) {
	magnitude := math.Sqrt(x*x + y*y + z*z)
	return x / magnitude, y / magnitude, z / magnitude
}

func DotProduct(x1, y1, z1, x2, y2, z2 float64) float64 {
	return x1*x2 + y1*y2 + z1*z2
}

func CalculateIntegratedElevationFromPoint(satellitePos polygon.Point, satelliteAltKm float64, groundPoint polygon.Point) float64 {
	// Convert the latitude/longitude of the ground point and satellite to 3D Cartesian coordinates
	groundX, groundY, groundZ := LatLonToCartesian(groundPoint.Latitude, groundPoint.Longitude, 0)
	satX, satY, satZ := LatLonToCartesian(satellitePos.Latitude, satellitePos.Longitude, satelliteAltKm)

	// Log the ground and satellite positions
	log.Printf("Ground Point (X, Y, Z): (%.2f, %.2f, %.2f)", groundX, groundY, groundZ)
	log.Printf("Satellite Position (X, Y, Z): (%.2f, %.2f, %.2f)", satX, satY, satZ)

	// Compute the vector from the ground to the satellite
	vecX, vecY, vecZ := satX-groundX, satY-groundY, satZ-groundZ

	// Compute the vector from the ground to Earth's center (opposite of the ground point)
	earthVecX, earthVecY, earthVecZ := Normalize(-groundX, -groundY, -groundZ)

	// Log vectors before normalization
	log.Printf("Vector from Ground to Satellite (X, Y, Z): (%.2f, %.2f, %.2f)", vecX, vecY, vecZ)
	log.Printf("Vector from Ground to Earth's Center (X, Y, Z): (%.2f, %.2f, %.2f)", earthVecX, earthVecY, earthVecZ)

	// Normalize the satellite vector
	vecX, vecY, vecZ = Normalize(vecX, vecY, vecZ)

	// Log normalized vectors
	log.Printf("Normalized Vector from Ground to Earth's Center: (%.2f, %.2f, %.2f)", earthVecX, earthVecY, earthVecZ)
	log.Printf("Normalized Vector from Ground to Satellite: (%.2f, %.2f, %.2f)", vecX, vecY, vecZ)

	// Special case handling for direct overhead
	if vecX == earthVecX && vecY == earthVecY && vecZ == earthVecZ {
		log.Printf("Satellite directly overhead: Elevation is 90 degrees.")
		return 90.0
	}

	// Compute the dot product of the vectors
	dotProd := DotProduct(vecX, vecY, vecZ, earthVecX, earthVecY, earthVecZ)

	// Log dot product
	log.Printf("Dot Product: %.2f", dotProd)

	// Compute the cosine of the elevation angle
	cosElevation := dotProd

	// Ensure the cosine stays within [-1, 1] range
	if cosElevation > 1.0 {
		cosElevation = 1.0
	} else if cosElevation < -1.0 {
		cosElevation = -1.0
	}

	// Compute the elevation angle in radians
	elevationRad := math.Acos(cosElevation)

	// Convert the elevation angle from radians to degrees
	elevation := RadiansToDegrees(elevationRad)

	// Log the final elevation
	log.Printf("Calculated Elevation: %.2f degrees", elevation)

	// Return the elevation
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

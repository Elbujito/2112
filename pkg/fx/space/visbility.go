package space

import (
	"log"
	"math"
	"time"

	"github.com/Elbujito/2112/pkg/fx/polygon"
	"github.com/joshuaferrara/go-satellite"
)

// ComputeVisibilityWindow computes the visibility window for a satellite over a given tile.
func ComputeVisibilityWindow(
	noradID, tleLine1, tleLine2 string,
	point polygon.Point,
	radius float64,
	startTime, endTime time.Time, timeStep time.Duration,
) (time.Time, time.Time, float64) {
	log.Printf("Starting visibility computation for NORAD ID: %s", noradID)

	maxElevation := -1.0

	tileRadiusKm := radius / 1000

	// Initialize satellite data
	satrec := satellite.TLEToSat(tleLine1, tleLine2, satellite.GravityWGS84)

	// Compute AOS
	aos := ComputeAOS(satrec, point, tileRadiusKm, startTime, endTime, timeStep, &maxElevation)
	if aos.IsZero() {
		log.Printf("AOS not detected for NORAD ID %s. Aborting computation.", noradID)
		return aos, time.Time{}, maxElevation
	}

	// // Compute LOS
	los := ComputeLOS(satrec, point, tileRadiusKm, aos, endTime, timeStep, &maxElevation)
	if los.IsZero() {
		log.Printf("LOS not detected for NORAD ID %s.", noradID)
	}

	log.Printf("Visibility computation completed for NORAD ID %s: AOS=%v, LOS=%v, MaxElevation=%.2f",
		noradID, aos, aos, maxElevation)

	return aos, aos, maxElevation
}

// ComputeAOS computes the Acquisition of Signal (AOS) time for a satellite over a given tile.
func ComputeAOS(
	satrec satellite.Satellite, point polygon.Point,
	tileRadiusKm float64, startTime, endTime time.Time,
	timeStep time.Duration, maxElevation *float64,
) time.Time {
	log.Printf("[INFO] Starting AOS computation.")

	for t := startTime; t.Before(endTime); t = t.Add(timeStep) {
		// Propagate the satellite's position
		altitude, geo, err := PropagateSatellitePosition(satrec, t)
		if err != nil {
			log.Printf("[ERROR] Propagation failed at %v: %v", t, err)
			continue
		}

		// Satellite's current position
		satellitePos := polygon.Point{Latitude: geo.Latitude, Longitude: geo.Longitude}
		log.Printf("[DEBUG] Time: %v, Satellite Position: Lat=%.6f, Lon=%.6f, Alt=%.2f km", t, geo.Latitude, geo.Longitude, altitude)

		if Intersects(point, satellitePos, tileRadiusKm, altitude) {
			log.Printf("[DEBUG] Time: %v, Satellite intersects edge within radius %.2f km", t, tileRadiusKm)

			// Calculate elevation
			elevation := CalculateIntegratedElevationFromPoint(satellitePos, altitude, point)
			log.Printf("[DEBUG] Time: %v, Calculated Elevation: %.2f degrees", t, elevation)

			// Check if AOS is valid
			if elevation > 0 {
				*maxElevation = math.Max(*maxElevation, elevation)
				log.Printf("[SUCCESS] AOS detected at %v. Max Elevation: %.2f degrees", t, *maxElevation)
				return t
			}
		}
	}

	log.Printf("[INFO] AOS not found within the specified time window.")
	return time.Time{}
}

// ComputeLOS computes the Loss of Signal (LOS) time for a satellite over a given tile.
func ComputeLOS(
	satrec satellite.Satellite, point polygon.Point,
	tileRadiusKm float64, aos time.Time, endTime time.Time,
	timeStep time.Duration, maxElevation *float64,
) time.Time {
	log.Printf("[INFO] Starting LOS computation.")

	// Start the LOS computation from the AOS time
	for t := aos; t.Before(endTime); t = t.Add(timeStep) {
		// Propagate the satellite's position at the current time
		altitude, geo, err := PropagateSatellitePosition(satrec, t)
		if err != nil {
			log.Printf("[ERROR] Propagation failed at %v: %v", t, err)
			continue
		}

		// Satellite's current position
		satellitePos := polygon.Point{Latitude: geo.Latitude, Longitude: geo.Longitude}
		log.Printf("[DEBUG] Time: %v, Satellite Position: Lat=%.6f, Lon=%.6f, Alt=%.2f km", t, geo.Latitude, geo.Longitude, altitude)

		// Check if the satellite is still within the tile's radius
		if !Intersects(point, satellitePos, tileRadiusKm, altitude) {
			log.Printf("[DEBUG] Time: %v, Satellite no longer intersects tile area, LOS detected.", t)

			// Calculate the elevation at LOS
			elevation := CalculateIntegratedElevationFromPoint(satellitePos, altitude, point)
			log.Printf("[DEBUG] Time: %v, LOS Elevation: %.2f degrees", t, elevation)

			// Return the time of LOS
			return t
		}
	}

	log.Printf("[INFO] LOS not found within the specified time window.")
	return time.Time{}
}

func Intersects(tileCenter polygon.Point, satellitePos polygon.Point, tileRadiusKm float64, altitude float64) bool {
	// Compute the distance from the satellite to the tile center point (ignoring altitude)
	centerDistance := HaversineDistance(satellitePos.Latitude, satellitePos.Longitude, tileCenter.Latitude, tileCenter.Longitude, 0, 0)

	log.Printf("[DEBUG] Distance to tile center: %.2f km", centerDistance)

	// Add a small margin of error to account for floating-point precision issues
	marginOfError := 0.01

	// Check if the satellite is within the tile radius from the center of the tile
	return centerDistance <= tileRadiusKm+marginOfError
}

// PropagateSatellitePosition calculates the satellite's geodetic position at a specific time.
func PropagateSatellitePosition(satrec satellite.Satellite, t time.Time) (float64, satellite.LatLong, error) {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)
	gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)
	altitude, _, geo := satellite.ECIToLLA(position, gmst)
	return altitude, geo, nil
}

// GeneratePolygonEdges converts polygon vertices into edges.
func GeneratePolygonEdges(vertices []polygon.Point) []polygon.Edge {
	var edges []polygon.Edge
	for i := 0; i < len(vertices)-1; i++ {
		edges = append(edges, polygon.Edge{
			Start: polygon.Point{Latitude: vertices[i].Latitude, Longitude: vertices[i].Longitude},
			End:   polygon.Point{Latitude: vertices[i+1].Latitude, Longitude: vertices[i+1].Longitude},
		})
	}
	edges = append(edges, polygon.Edge{
		Start: polygon.Point{Latitude: vertices[len(vertices)-1].Latitude, Longitude: vertices[len(vertices)-1].Longitude},
		End:   polygon.Point{Latitude: vertices[0].Latitude, Longitude: vertices[0].Longitude},
	})
	return edges
}

// CalculateIntegratedElevation computes the elevation of a satellite relative to a ground point.
func CalculateIntegratedElevation(satelliteQuadKey polygon.Quadkey, satelliteAltitude float64, pointQuadKey polygon.Point) float64 {
	return CalculateIntegratedElevationFromPoint(
		polygon.Point{Latitude: satelliteQuadKey.Latitude, Longitude: satelliteQuadKey.Longitude},
		satelliteAltitude,
		polygon.Point{Latitude: pointQuadKey.Latitude, Longitude: pointQuadKey.Longitude},
	)
}

func HaversineDistance(lat1, lon1, lat2, lon2, altitude1, altitude2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := DegreesToRadians(lat2 - lat1)
	dLon := DegreesToRadians(lon2 - lon1)
	lat1Rad := DegreesToRadians(lat1)
	lat2Rad := DegreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	surfaceDistance := earthRadiusKm * c

	// Adjust for altitude (adding z-axis distance)
	altitudeDiff := altitude2 - altitude1
	return math.Sqrt(surfaceDistance*surfaceDistance + altitudeDiff*altitudeDiff)
}

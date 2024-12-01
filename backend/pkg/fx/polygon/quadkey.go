package polygon

import (
	"fmt"
	"math"
)

type Quadkey struct {
	Latitude  float64
	Longitude float64
	Level     int
}

func NewQuadkey(lat float64, long float64, level int) Quadkey {
	return Quadkey{
		Latitude:  lat,
		Longitude: long,
		Level:     level,
	}
}

func (q *Quadkey) Key() string {
	return fmt.Sprintf("%d-%f-%f", q.Level, q.Latitude, q.Longitude)
}

// Constants
const EarthRadiusKm = 6371.0 // Earth's radius in kilometers

// DistanceTo computes the great-circle distance between two Quadkeys using the haversine formula
func (q *Quadkey) DistanceTo(other Point) float64 {
	// Convert latitudes and longitudes from degrees to radians
	lat1 := q.Latitude * math.Pi / 180
	lon1 := q.Longitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	lon2 := other.Longitude * math.Pi / 180

	// Compute differences
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	return EarthRadiusKm * c
}

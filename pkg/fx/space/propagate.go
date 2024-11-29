package space

import (
	"fmt"
	"time"

	"github.com/Elbujito/2112/pkg/fx/polygon"
	"github.com/joshuaferrara/go-satellite"
)

// PropagateSatellite propagates the satellite's position to the specified time.
// Returns a QuadKey for the satellite's position at the given time or an error.
func PropagateSatellite(tleLine1, tleLine2 string, t time.Time) (polygon.Quadkey, satellite.Satellite, error) {
	// Create satellite record from TLE lines
	satrec := satellite.TLEToSat(tleLine1, tleLine2, satellite.GravityWGS84)

	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)
	if satrec.Error != 0 {
		return polygon.Quadkey{}, satellite.Satellite{}, fmt.Errorf("propagation error code: %d", satrec.Error)
	}

	// Calculate GST for ECI to LLA conversion
	gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)

	// Convert ECI to Geodetic (lat, lon, alt)
	altitude, _, geoPosition := satellite.ECIToLLA(position, gmst)

	quadKey := polygon.NewQuadkey(geoPosition.Latitude, geoPosition.Longitude, int(altitude))
	return quadKey, satrec, nil
}

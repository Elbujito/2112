package mappers

import "time"

type RawTLE struct {
	NoradID string `json:"norad_id"`
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
}

// SatelliteMetadata represents basic information about a satellite.
type SatelliteMetadata struct {
	NoradID        string
	Name           string
	IntlDesignator string
	LaunchDate     time.Time
	DecayDate      *time.Time
	ObjectType     string
	Owner          string
	Period         *float64
	Inclination    *float64
	Apogee         *float64
	Perigee        *float64
	RCS            *float64
	Altitude       *float64 // New field for altitude
}

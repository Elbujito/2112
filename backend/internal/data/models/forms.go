package models

import (
	"time"
)

type FormBase struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SatelliteForm struct {
	FormBase
	Name           string     `json:"name" validate:"required,min=2,max=255"`       // Satellite name
	NoradID        string     `json:"norad_id" validate:"required,min=1,max=255"`   // NORAD ID (unique identifier)
	Type           string     `json:"type" validate:"required,min=2,max=255"`       // Satellite type
	LaunchDate     *time.Time `json:"launch_date" validate:"omitempty"`             // Launch date
	DecayDate      *time.Time `json:"decay_date" validate:"omitempty"`              // Decay date (optional)
	IntlDesignator string     `json:"intl_designator" validate:"omitempty,max=255"` // International designator
	Owner          string     `json:"owner" validate:"omitempty,max=255"`           // Ownership information
	ObjectType     string     `json:"object_type" validate:"omitempty,max=255"`     // Object type
	Period         *float64   `json:"period" validate:"omitempty"`                  // Orbital period in minutes
	Inclination    *float64   `json:"inclination" validate:"omitempty"`             // Orbital inclination in degrees
	Apogee         *float64   `json:"apogee" validate:"omitempty"`                  // Apogee altitude in kilometers
	Perigee        *float64   `json:"perigee" validate:"omitempty"`                 // Perigee altitude in kilometers
	RCS            *float64   `json:"rcs" validate:"omitempty"`                     // Radar cross-section in square meters
}

func (form *SatelliteForm) MapToModel() *Satellite {
	return &Satellite{
		Name:           form.Name,
		NoradID:        form.NoradID,
		Type:           form.Type,
		LaunchDate:     form.LaunchDate,
		DecayDate:      form.DecayDate,
		IntlDesignator: form.IntlDesignator,
		Owner:          form.Owner,
		ObjectType:     form.ObjectType,
		Period:         form.Period,
		Inclination:    form.Inclination,
		Apogee:         form.Apogee,
		Perigee:        form.Perigee,
		RCS:            form.RCS,
	}
}

type TLEForm struct {
	FormBase
	NoradID string `json:"norad_id" validate:"required"`                            // Reference to the Satellite via NORAD ID
	Line1   string `json:"line1" validate:"required"`                               // First line of TLE
	Line2   string `json:"line2" validate:"required"`                               // Second line of TLE
	Epoch   string `json:"epoch" validate:"required,datetime=2006-01-02T15:04:05Z"` // Epoch time (ISO8601 format)
}

// MapToModel maps TLEForm to a TLE model
func (form *TLEForm) MapToModel() *TLE {
	epochTime, _ := time.Parse(time.RFC3339, form.Epoch) // Assuming validation ensures correct parsing
	return &TLE{
		NoradID: form.NoradID,
		Line1:   form.Line1,
		Line2:   form.Line2,
		Epoch:   epochTime,
	}
}

type TileForm struct {
	FormBase
	Quadkey   string  `json:"quadkey" validate:"required"`                     // Unique Quadkey for the tile
	ZoomLevel int     `json:"zoom_level" validate:"required,min=1,max=20"`     // Zoom level (1 = world, higher = more detailed)
	CenterLat float64 `json:"center_lat" validate:"required,min=-90,max=90"`   // Center latitude of the tile
	CenterLon float64 `json:"center_lon" validate:"required,min=-180,max=180"` // Center longitude of the tile
}

func (form *TileForm) MapToModel() *Tile {
	return &Tile{
		Quadkey:   form.Quadkey,
		ZoomLevel: form.ZoomLevel,
		CenterLat: form.CenterLat,
		CenterLon: form.CenterLon,
	}
}

type TileSatelliteMappingForm struct {
	FormBase
	NoradID      string  `json:"norad_id"`   // NORAD ID for the satellite
	TileID       string  `json:"tile_id"`    // Tile ID (string to match updated schema)
	Aos          string  `json:"start_time"` // ISO8601 format
	EndTime      string  `json:"end_time"`   // ISO8601 format
	MaxElevation float64 `json:"max_elevation"`
}

// MapToModel converts a VisibilityForm to a Visibility model
func (form *TileSatelliteMappingForm) MapToModel() *TileSatelliteMapping {
	startTime, _ := time.Parse(time.RFC3339, form.Aos) // Assuming validation ensures correct parsing
	return &TileSatelliteMapping{
		NoradID:      form.NoradID,
		TileID:       form.TileID, // Now TileID is a string
		Aos:          startTime,
		MaxElevation: form.MaxElevation,
	}
}

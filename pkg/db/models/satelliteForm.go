package models

import "time"

type SatelliteForm struct {
	FormBase
	Name    string `json:"name" validate:"required,min=2,max=50"`     // Satellite name
	NoradID string `json:"norad_id" validate:"required,min=1,max=10"` // NORAD ID (unique identifier)
	Type    string `json:"type" validate:"required,min=2,max=80"`     // Satellite type
}

func (form *SatelliteForm) MapToModel() *Satellite {
	return &Satellite{
		Name:    form.Name,
		NoradID: form.NoradID,
		Type:    form.Type,
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

type VisibilityForm struct {
	FormBase
	NoradID      string  `json:"norad_id"`   // NORAD ID for the satellite
	TileID       uint    `json:"tile_id"`    // Tile ID
	StartTime    string  `json:"start_time"` // ISO8601 format
	EndTime      string  `json:"end_time"`   // ISO8601 format
	MaxElevation float64 `json:"max_elevation"`
}

func (form *VisibilityForm) MapToModel() *Visibility {
	startTime, _ := time.Parse(time.RFC3339, form.StartTime) // Assuming validation ensures correct parsing
	endTime, _ := time.Parse(time.RFC3339, form.EndTime)
	return &Visibility{
		NoradID:      form.NoradID,
		TileID:       form.TileID,
		StartTime:    startTime,
		EndTime:      endTime,
		MaxElevation: form.MaxElevation,
	}
}

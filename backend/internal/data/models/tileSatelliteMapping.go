package models

import "time"

type TileSatelliteMapping struct {
	ModelBase
	NoradID      string    `gorm:"size:255;not null;index"` // Foreign key to Satellite table via NORAD ID
	TileID       string    `gorm:"type:char(36);not null"`  // Foreign key to Tile table
	Aos          time.Time `gorm:"not null"`                // Visibility start time
	MaxElevation float64   `gorm:"not null"`                // Max elevation during visibility in degrees
}

// MapToForm converts the Visibility model to a form structure
func (model *TileSatelliteMapping) MapToForm() *TileSatelliteMappingForm {
	return &TileSatelliteMappingForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		NoradID:      model.NoradID,
		TileID:       model.TileID,
		Aos:          model.Aos.Format(time.RFC3339), // ISO8601 format
		MaxElevation: model.MaxElevation,
	}
}

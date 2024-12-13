package models

// TileSatelliteMapping defines the relationship between a satellite and a tile
type TileSatelliteMapping struct {
	ModelBase
	NoradID string `gorm:"size:255;not null;index"`                       // Foreign key to Satellite table via NORAD ID
	TileID  string `gorm:"type:char(36);not null;uniqueIndex:norad_tile"` // Foreign key to Tile table
}

// MapToForm converts the TileSatelliteMapping model to a form structure
func (model *TileSatelliteMapping) MapToForm() *TileSatelliteMappingForm {
	return &TileSatelliteMappingForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		NoradID: model.NoradID,
		TileID:  model.TileID,
	}
}

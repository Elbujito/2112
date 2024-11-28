package models

// Tile Model
type Tile struct {
	ModelBase
	Quadkey        string  `gorm:"size:256;unique;not null"` // Unique identifier for the tile (Quadkey)
	ZoomLevel      int     `gorm:"not null"`                 // Zoom level for the tile
	CenterLat      float64 `gorm:"not null"`                 // Center latitude of the tile
	CenterLon      float64 `gorm:"not null"`                 // Center longitude of the tile
	NbFaces        int     `gorm:"not null"`
	Radius         float64 `gorm:"not null"`
	BoundariesJSON string  `gorm:"type:json"` // Serialized JSON of Boundaries for persistence
}

// MapToForm converts the Tile model to a TileForm structure
func (model *Tile) MapToForm() *TileForm {
	return &TileForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		Quadkey:   model.Quadkey,
		ZoomLevel: model.ZoomLevel,
		CenterLat: model.CenterLat,
		CenterLon: model.CenterLon,
	}
}

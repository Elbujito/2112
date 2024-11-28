package models

import (
	"github.com/Elbujito/2112/pkg/fx/polygon"
)

// TileService defines the operations required for Tile management
type TileService interface {
	FindByQuadkey(quadkey string) (*Tile, error)
	FindAll() ([]*Tile, error)
	Save(tile *Tile) error
	Create() error
	Update(tile *Tile) error
	Delete(id string) error
}

var tile *Tile = &Tile{}

// Tile Model
type Tile struct {
	ModelBase
	Quadkey   string  `gorm:"size:256;unique;not null"` // Unique identifier for the tile (Quadkey)
	ZoomLevel int     `gorm:"not null"`                 // Zoom level for the tile
	CenterLat float64 `gorm:"not null"`                 // Center latitude of the tile
	CenterLon float64 `gorm:"not null"`                 // Center longitude of the tile
	// Polygon        polygon.Polygon   `gorm:"embedded;embeddedPrefix:polygon_"` // Polygon associated with the tile
	// Boundaries     []polygon.Quadkey `gorm:"-"`         // Polygon boundaries (runtime only)
	// BoundariesJSON string            `gorm:"type:json"` // Serialized JSON of Boundaries for persistence
}

// NewTile creates a new Tile instance
func NewTile(quadkey string, zoomLevel int, centerLat, centerLon float64, polygon polygon.Polygon) *Tile {
	return &Tile{
		Quadkey:   quadkey,
		ZoomLevel: zoomLevel,
		CenterLat: centerLat,
		CenterLon: centerLon,
		// Polygon:   polygon,
	}
}

// // BeforeSave prepares the Boundaries field for database persistence
// func (t *Tile) BeforeSave() (err error) {
// 	boundariesJSON, err := json.Marshal(t.Boundaries)
// 	if err != nil {
// 		return err
// 	}
// 	t.BoundariesJSON = string(boundariesJSON)
// 	return nil
// }

// // AfterFind populates the Boundaries field from the stored JSON
// func (t *Tile) AfterFind() (err error) {
// 	if t.BoundariesJSON != "" {
// 		if err := json.Unmarshal([]byte(t.BoundariesJSON), &t.Boundaries); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// TileModel returns a reference to the Tile model
func TileModel() *Tile {
	return tile
}

// FindAll retrieves all Tile records
func (model *Tile) FindAll() (models []*Tile, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

// FindByQuadkey retrieves a Tile record by its Quadkey
func (model *Tile) FindByQuadkey(quadkey string) (m *Tile, err error) {
	result := db.Model(model).Where("quadkey=?", quadkey).First(&m)
	return m, result.Error
}

// Save persists a new Tile record
func (model *Tile) Save(m *Tile) error {
	return db.Model(model).Create(&m).Error
}

// Create persists a new Tile record
func (model *Tile) Create() error {
	return db.Create(model).Error
}

// Update modifies an existing Tile record
func (model *Tile) Update(m *Tile) error {
	return db.Model(model).Save(m).Error
}

// Delete removes a Tile record by its ID
func (model *Tile) Delete(id string) error {
	return db.Model(model).Where("ID=?", id).Delete(&model).Error
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

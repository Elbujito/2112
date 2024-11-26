package models

import (
	"time"
)

var satellite *Satellite = &Satellite{}
var tile *Tile = &Tile{}
var tle *TLE = &TLE{}
var visibility *Visibility = &Visibility{}

// Satellite Model
type Satellite struct {
	ModelBase
	Name    string `gorm:"size:255;not null"`        // Satellite name
	NoradID string `gorm:"size:255;unique;not null"` // NORAD ID
	Type    string `gorm:"size:255"`                 // Satellite type (e.g., telescope, communication)
}

func SatelliteModel() *Satellite {
	return satellite
}

func (model *Satellite) FindAll() (models []*Satellite, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

func (model *Satellite) Find(id string) (m *Satellite, err error) {
	result := db.Model(model).Where("ID=?", id).First(&m)
	return m, result.Error
}

func (model *Satellite) Save() error {
	return db.Model(model).Create(&model).Error
}

func (model *Satellite) Update() error {
	return db.Model(model).Updates(&model).Error
}

func (model *Satellite) Delete(id string) error {
	return db.Model(model).Where("ID=?", id).Delete(&model).Error
}

func (model *Satellite) MapToForm() *SatelliteForm {
	return &SatelliteForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		Name:    model.Name,
		NoradID: model.NoradID,
		Type:    model.Type,
	}
}

// TLE Model
type TLE struct {
	ModelBase
	SatelliteID uint      `gorm:"not null;index"` // Foreign key to Satellite table
	Line1       string    `gorm:"size:255;not null"`
	Line2       string    `gorm:"size:255;not null"`
	Epoch       time.Time `gorm:"not null"` // Time associated with the TLE
}

func TLEModel() *TLE {
	return tle
}

func (model *TLE) FindAll() (models []*TLE, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

func (model *TLE) FindBySatelliteID(satelliteID uint) (models []*TLE, err error) {
	result := db.Model(model).Where("satellite_id=?", satelliteID).Find(&models)
	return models, result.Error
}

func (model *TLE) Save() error {
	return db.Model(model).Create(&model).Error
}

func (model *TLE) Delete(id string) error {
	return db.Model(model).Where("ID=?", id).Delete(&model).Error
}

func (model *TLE) MapToForm() *TLEForm {
	return &TLEForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		SatelliteID: model.SatelliteID,
		Line1:       model.Line1,
		Line2:       model.Line2,
		Epoch:       model.Epoch.Format(time.RFC3339), // ISO8601 format
	}
}

// Tile Model
type Tile struct {
	ModelBase
	Quadkey   string  `gorm:"size:25;unique;not null"` // Unique identifier for the tile (Quadkey)
	ZoomLevel int     `gorm:"not null"`                // Zoom level for the tile
	CenterLat float64 `gorm:"not null"`                // Center latitude of the tile
	CenterLon float64 `gorm:"not null"`                // Center longitude of the tile
}

func TileModel() *Tile {
	return tile
}

func (model *Tile) FindAll() (models []*Tile, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

func (model *Tile) FindByQuadkey(quadkey string) (m *Tile, err error) {
	result := db.Model(model).Where("quadkey=?", quadkey).First(&m)
	return m, result.Error
}

func (model *Tile) Save() error {
	return db.Model(model).Create(&model).Error
}

func (model *Tile) Delete(id string) error {
	return db.Model(model).Where("ID=?", id).Delete(&model).Error
}

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

// Visibility Model
type Visibility struct {
	ModelBase
	SatelliteID  uint      `gorm:"not null;index"` // Foreign key to Satellite table
	TileID       uint      `gorm:"not null;index"` // Foreign key to Tile table
	StartTime    time.Time `gorm:"not null"`       // Visibility start time
	EndTime      time.Time `gorm:"not null"`       // Visibility end time
	MaxElevation float64   `gorm:"not null"`       // Max elevation during visibility in degrees
}

func VisibilityModel() *Visibility {
	return visibility
}

func (model *Visibility) FindAll() (models []*Visibility, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

func (model *Visibility) FindBySatelliteAndTile(satelliteID uint, tileID uint) (models []*Visibility, err error) {
	result := db.Model(model).Where("satellite_id=? AND tile_id=?", satelliteID, tileID).Find(&models)
	return models, result.Error
}

func (model *Visibility) Save() error {
	return db.Model(model).Create(&model).Error
}

func (model *Visibility) Delete(id string) error {
	return db.Model(model).Where("ID=?", id).Delete(&model).Error
}

func (model *Visibility) MapToForm() *VisibilityForm {
	return &VisibilityForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		SatelliteID:  model.SatelliteID,
		TileID:       model.TileID,
		StartTime:    model.StartTime.Format(time.RFC3339), // ISO8601 format
		EndTime:      model.EndTime.Format(time.RFC3339),   // ISO8601 format
		MaxElevation: model.MaxElevation,
	}
}

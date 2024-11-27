package models

import (
	"time"
)

// SatelliteService defines the operations required for Satellite management
type SatelliteService interface {
	FindByNoradID(noradID string) (*Satellite, error)
	Find(id string) (*Satellite, error)
	FindAll() ([]*Satellite, error)
	Save(satellite *Satellite) error
	Update(satellite *Satellite) error
	Delete(id string) error
}

// TLEService defines the operations required for TLE management
type TLEService interface {
	FindByNoradID(noradID string) ([]*TLE, error)
	FindAll() ([]*TLE, error)
	Save(tle *TLE) error
	Update(tle *TLE) error
	Delete(id string) error
}

// TileService defines the operations required for Tile management
type TileService interface {
	FindByQuadkey(quadkey string) (*Tile, error)
	FindAll() ([]*Tile, error)
	Save(tile *Tile) error
	Create() error
	Update(tile *Tile) error
	Delete(id string) error
}

// VisibilityService defines the operations required for Visibility management
type VisibilityService interface {
	FindByNoradIDAndTile(noradID string, tileID string) ([]*Visibility, error)
	FindAll() ([]*Visibility, error)
	Save(visibility *Visibility) error
	Create() error
	Update(visibility *Visibility) error
	Delete(id string) error
}

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

// FindByNoradID retrieves all TLE records for a given NoradID
func (model *Satellite) FindByNoradID(noradID string) (models *Satellite, err error) {
	result := db.Model(model).Where("norad_id=?", noradID).Find(&models)
	return models, result.Error
}

func (model *Satellite) FindAll() (models []*Satellite, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

func (model *Satellite) Find(id string) (m *Satellite, err error) {
	result := db.Model(model).Where("ID=?", id).First(&m)
	return m, result.Error
}

func (model *Satellite) Save(satellite *Satellite) error {
	return db.Model(model).Create(&model).Error
}

// Create inserts a new Satellite record with auto-generated ID
func (model *Satellite) Create() error {
	return db.Create(model).Error
}

func (model *Satellite) Update(satellite *Satellite) error {
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
	NoradID string    `gorm:"size:255;not null;index"` // Foreign key to Satellite table via Norad ID
	Line1   string    `gorm:"size:255;not null"`
	Line2   string    `gorm:"size:255;not null"`
	Epoch   time.Time `gorm:"not null"` // Time associated with the TLE
}

// TLEModel returns a reference to the TLE model
func TLEModel() *TLE {
	return tle
}

// FindAll retrieves all TLE records
func (model *TLE) FindAll() (models []*TLE, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

// FindByNoradID retrieves all TLE records for a given NoradID
func (model *TLE) FindByNoradID(noradID string) (models []*TLE, err error) {
	result := db.Model(model).Where("norad_id=?", noradID).Find(&models)
	return models, result.Error
}

// Save persists a new TLE record
func (model *TLE) Save(m *TLE) error {
	return db.Model(model).Create(&model).Error
}

// Create persists a new TLE record
func (model *TLE) Create() error {
	return db.Create(model).Error
}

// Update modifies an existing TLE record
func (model *TLE) Update(m *TLE) error {
	return db.Model(model).Save(&model).Error
}

// Delete removes a TLE record by its ID
func (model *TLE) Delete(id string) error {
	return db.Model(model).Where("ID=?", id).Delete(&model).Error
}

// MapToForm converts the TLE model to a TLEForm structure
func (model *TLE) MapToForm() *TLEForm {
	return &TLEForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		NoradID: model.NoradID,
		Line1:   model.Line1,
		Line2:   model.Line2,
		Epoch:   model.Epoch.Format(time.RFC3339), // ISO8601 format
	}
}

// Tile Model
type Tile struct {
	ModelBase
	Quadkey   string  `gorm:"size:256;unique;not null"` // Unique identifier for the tile (Quadkey)
	ZoomLevel int     `gorm:"not null"`                 // Zoom level for the tile
	CenterLat float64 `gorm:"not null"`                 // Center latitude of the tile
	CenterLon float64 `gorm:"not null"`                 // Center longitude of the tile
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

func (model *Tile) Save(m *Tile) error {
	return db.Model(model).Create(&m).Error
}

func (model *Tile) Create() error {
	return db.Create(model).Error
}

func (model *Tile) Update(m *Tile) error {
	return db.Model(model).Save(m).Error
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

type Visibility struct {
	ModelBase
	NoradID      string    `gorm:"size:255;not null;index"` // Foreign key to Satellite table via NORAD ID
	TileID       string    `gorm:"type:char(36);not null"`  // Foreign key to Tile table
	StartTime    time.Time `gorm:"not null"`                // Visibility start time
	EndTime      time.Time `gorm:"not null"`                // Visibility end time
	MaxElevation float64   `gorm:"not null"`                // Max elevation during visibility in degrees
}

func VisibilityModel() *Visibility {
	return visibility
}

// FindAll retrieves all Visibility records
func (model *Visibility) FindAll() (models []*Visibility, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

// FindByNoradIDAndTile retrieves visibility records for a given NORAD ID and Tile ID
func (model *Visibility) FindByNoradIDAndTile(noradID string, tileID string) (models []*Visibility, err error) {
	result := db.Model(model).Where("norad_id = ? AND tile_id = ?", noradID, tileID).Find(&models)
	return models, result.Error
}

// Save persists a new or existing Visibility record
func (model *Visibility) Save() error {
	return db.Model(model).Save(model).Error
}

// Create inserts a new Visibility record
func (model *Visibility) Create() error {
	return db.Create(model).Error
}

// Update modifies an existing Visibility record
func (model *Visibility) Update(visibility *Visibility) error {
	return db.Model(model).Save(visibility).Error
}

// Delete removes a Visibility record by its ID
func (model *Visibility) Delete(id string) error {
	return db.Model(model).Where("ID = ?", id).Delete(model).Error
}

// MapToForm converts the Visibility model to a form structure
func (model *Visibility) MapToForm() *VisibilityForm {
	return &VisibilityForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		NoradID:      model.NoradID,
		TileID:       model.TileID,
		StartTime:    model.StartTime.Format(time.RFC3339), // ISO8601 format
		EndTime:      model.EndTime.Format(time.RFC3339),   // ISO8601 format
		MaxElevation: model.MaxElevation,
	}
}

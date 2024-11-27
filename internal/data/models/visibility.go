package models

import "time"

// VisibilityService defines the operations required for Visibility management
type VisibilityService interface {
	FindByNoradIDAndTile(noradID string, tileID string) ([]*Visibility, error)
	FindAll() ([]*Visibility, error)
	Save(visibility *Visibility) error
	Create() error
	Update(visibility *Visibility) error
	Delete(id string) error
}

var visibility *Visibility = &Visibility{}

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

package models

import "time"

// TLEService defines the operations required for TLE management
type TLEService interface {
	FindByNoradID(noradID string) ([]*TLE, error)
	FindAll() ([]*TLE, error)
	Save(tle *TLE) error
	Update(tle *TLE) error
	Delete(id string) error
	Upsert(tle *TLE) error
}

var tle *TLE = &TLE{}

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

// Upsert inserts a new TLE record or updates an existing one based on NoradID
func (model *TLE) Upsert(tle *TLE) error {
	existingTLEs, err := model.FindByNoradID(tle.NoradID)
	if err != nil {
		return err
	}
	if len(existingTLEs) > 0 {
		// Update the first matching record
		existingTLE := existingTLEs[0]
		existingTLE.Line1 = tle.Line1
		existingTLE.Line2 = tle.Line2
		existingTLE.Epoch = tle.Epoch
		return model.Update(existingTLE)
	}
	// Insert as a new record
	return model.Save(tle)
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

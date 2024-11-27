package models

// SatelliteService defines the operations required for Satellite management
type SatelliteService interface {
	FindByNoradID(noradID string) (*Satellite, error)
	Find(id string) (*Satellite, error)
	FindAll() ([]*Satellite, error)
	Save(satellite *Satellite) error
	Update(satellite *Satellite) error
	Delete(id string) error
}

var satellite *Satellite = &Satellite{}

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

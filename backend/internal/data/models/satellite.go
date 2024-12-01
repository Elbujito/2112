package models

// Satellite Model
type Satellite struct {
	ModelBase
	Name    string `gorm:"size:255;not null"`        // Satellite name
	NoradID string `gorm:"size:255;unique;not null"` // NORAD ID
	Type    string `gorm:"size:255"`                 // Satellite type (e.g., telescope, communication)
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

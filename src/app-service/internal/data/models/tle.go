package models

import "time"

// TLE Model
type TLE struct {
	ModelBase
	NoradID string    `gorm:"size:255;not null;index"` // Foreign key to Satellite table via Norad ID
	Line1   string    `gorm:"size:255;not null"`
	Line2   string    `gorm:"size:255;not null"`
	Epoch   time.Time `gorm:"not null"` // Time associated with the TLE
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

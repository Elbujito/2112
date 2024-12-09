package models

import (
	"time"
)

// Satellite Model
type Satellite struct {
	ModelBase
	Name           string     `gorm:"size:255;not null"`        // Satellite name
	NoradID        string     `gorm:"size:255;unique;not null"` // NORAD ID
	Type           string     `gorm:"size:255"`                 // Satellite type (e.g., telescope, communication)
	LaunchDate     *time.Time `gorm:"type:date"`                // Launch date
	DecayDate      *time.Time `gorm:"type:date"`                // Decay date (optional)
	IntlDesignator string     `gorm:"size:255"`                 // International designator
	Owner          string     `gorm:"size:255"`                 // Ownership information
	ObjectType     string     `gorm:"size:255"`                 // Object type (e.g., "PAYLOAD")
	Period         *float64   `gorm:"type:float"`               // Orbital period in minutes (optional)
	Inclination    *float64   `gorm:"type:float"`               // Orbital inclination in degrees (optional)
	Apogee         *float64   `gorm:"type:float"`               // Apogee altitude in kilometers (optional)
	Perigee        *float64   `gorm:"type:float"`               // Perigee altitude in kilometers (optional)
	RCS            *float64   `gorm:"type:float"`               // Radar cross-section in square meters (optional)
	Altitude       *float64   `gorm:"type:float"`               // Radar cross-section in square meters (optional)
}

// MapToForm maps the Satellite model to a SatelliteForm.
func (model *Satellite) MapToForm() *SatelliteForm {
	return &SatelliteForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		Name:           model.Name,
		NoradID:        model.NoradID,
		Type:           model.Type,
		LaunchDate:     model.LaunchDate,
		DecayDate:      model.DecayDate,
		IntlDesignator: model.IntlDesignator,
		Owner:          model.Owner,
		ObjectType:     model.ObjectType,
		Period:         model.Period,
		Inclination:    model.Inclination,
		Apogee:         model.Apogee,
		Perigee:        model.Perigee,
		RCS:            model.RCS,
		Altitude:       model.Altitude,
	}
}

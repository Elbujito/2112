package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// SatelliteType represents the type of a satellite.
type SatelliteType string

const (
	// Active satellite type.
	Active SatelliteType = "ACTIVE"
	// Other satellite type from SATCAT catalogue.
	Other SatelliteType = "OTHER"
)

// IsValid checks if the SatelliteType is valid.
func (t SatelliteType) IsValid() error {
	switch t {
	case Active, Other:
		return nil
	default:
		return errors.New("invalid satellite type")
	}
}

// Satellite represents the domain entity for a satellite.
type Satellite struct {
	ID        string        // Unique identifier
	CreatedAt time.Time     // Timestamp of creation
	UpdatedAt time.Time     // Timestamp of last update
	Name      string        // Satellite name
	NoradID   string        // NORAD ID
	Type      SatelliteType // Satellite type
}

// NewSatellite creates a new Satellite instance.
func NewSatellite(name, noradID string, satType SatelliteType) (*Satellite, error) {
	if err := satType.IsValid(); err != nil {
		return nil, err
	}
	return &Satellite{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		NoradID:   noradID,
		Type:      satType,
	}, nil
}

// SatelliteRepository defines the interface for Satellite operations.
type SatelliteRepository interface {
	FindByNoradID(noradID string) (Satellite, error)
	Find(id string) (Satellite, error)
	FindAll() ([]Satellite, error)
	Save(satellite Satellite) error
	Update(satellite Satellite) error
	Delete(id string) error
}

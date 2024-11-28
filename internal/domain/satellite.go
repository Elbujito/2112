package domain

import (
	"context"
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
	ID        string // Unique identifier
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	NoradID   string
	Type      SatelliteType
}

// NewSatellite creates a new Satellite instance.
func NewSatellite(name string, noradID string, satType SatelliteType) (Satellite, error) {
	if err := satType.IsValid(); err != nil {
		return Satellite{}, err
	}
	return Satellite{
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
	FindByNoradID(ctx context.Context, noradID string) (Satellite, error)
	FindAll(ctx context.Context) ([]Satellite, error)
	Save(ctx context.Context, satellite Satellite) error
	Update(ctx context.Context, satellite Satellite) error
	DeleteByNoradID(ctx context.Context, noradID string) error
}

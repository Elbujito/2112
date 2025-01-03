package domain

import (
	"context"
	"errors"
	"sort"
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

type Satellite struct {
	ID             string // Unique identifier
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	NoradID        string
	Type           SatelliteType
	LaunchDate     *time.Time // Added field for launch date
	DecayDate      *time.Time // Added field for decay date, if applicable
	IntlDesignator string     // Added field for international designator
	Owner          string     // Added field for ownership information
	ObjectType     string     // Added field for object type (e.g., PAYLOAD)
	Period         *float64   // Added field for orbital period in minutes
	Inclination    *float64   // Added field for orbital inclination in degrees
	Apogee         *float64   // Added field for apogee altitude in kilometers
	Perigee        *float64   // Added field for perigee altitude in kilometers
	RCS            *float64   // Added field for radar cross-section in square meters
	TleUpdatedAt   *time.Time `gorm:"-"`
	Altitude       *float64
}

// NewSatelliteFromStatCat creates a new Satellite instance with optional SATCAT data.
func NewSatelliteFromStatCat(
	name string,
	noradID string,
	satType SatelliteType,
	launchDate *time.Time,
	decayDate *time.Time,
	intlDesignator string,
	owner string,
	objectType string,
	period *float64,
	inclination *float64,
	apogee *float64,
	perigee *float64,
	rcs *float64,
	altitude *float64,
) (Satellite, error) {
	if err := satType.IsValid(); err != nil {
		return Satellite{}, err
	}
	return Satellite{
		ID:             uuid.NewString(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Name:           name,
		NoradID:        noradID,
		Type:           satType,
		LaunchDate:     launchDate,
		DecayDate:      decayDate,
		IntlDesignator: intlDesignator,
		Owner:          owner,
		ObjectType:     objectType,
		Period:         period,
		Inclination:    inclination,
		Apogee:         apogee,
		Perigee:        perigee,
		RCS:            rcs,
	}, nil
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
	SaveBatch(ctx context.Context, satellites []Satellite) error
	FindAllWithPagination(ctx context.Context, page int, pageSize int, searchRequest *SearchRequest) ([]Satellite, int64, error)
	FindSatelliteInfoWithPagination(ctx context.Context, page int, pageSize int, searchRequest *SearchRequest) ([]SatelliteInfo, int64, error)
}

type SatellitePosition struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Time      time.Time `json:"time"`
}

type SatelliteInfo struct {
	Satellite Satellite // The associated Satellite
	TLEs      []TLE     // List of TLEs ordered by most recent
}

// NewSatelliteInfo creates a new SatelliteInfo instance.
func NewSatelliteInfo(satellite Satellite, tles []TLE) SatelliteInfo {
	// Sort the TLEs by Epoch in descending order (most recent first)
	sort.Slice(tles, func(i, j int) bool {
		return tles[i].Epoch.After(tles[j].Epoch)
	})

	return SatelliteInfo{
		Satellite: satellite,
		TLEs:      tles,
	}
}

// GetMostRecentTLE returns the most recent TLE from the SatelliteInfo.
func (info *SatelliteInfo) GetMostRecentTLE() *TLE {
	if len(info.TLEs) == 0 {
		return nil
	}
	return &info.TLEs[0]
}

// AddTLE adds a new TLE to the SatelliteInfo and keeps the list sorted by most recent.
func (info *SatelliteInfo) AddTLE(tle TLE) {
	info.TLEs = append(info.TLEs, tle)
	// Re-sort to ensure the TLEs are still ordered by most recent
	sort.Slice(info.TLEs, func(i, j int) bool {
		return info.TLEs[i].Epoch.After(info.TLEs[j].Epoch)
	})
}

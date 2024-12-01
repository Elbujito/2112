package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	xtime "github.com/Elbujito/2112/pkg/fx/time"
	"github.com/google/uuid"
)

// TLE represents the domain entity for Two-Line Element sets.
type TLE struct {
	ID      string    // Unique identifier
	NoradID string    // NORAD ID associated with the satellite
	Line1   string    // First line of the TLE
	Line2   string    // Second line of the TLE
	Epoch   time.Time // Time associated with the TLE
}

// Validate ensures that the TLE fields are valid.
func (tle *TLE) Validate() error {
	if tle.NoradID == "" {
		return errors.New("NORAD ID cannot be empty")
	}
	if tle.Line1 == "" || tle.Line2 == "" {
		return errors.New("TLE lines cannot be empty")
	}
	if tle.Epoch.IsZero() {
		return errors.New("epoch cannot be zero")
	}
	return nil
}

// NewTLE creates a new TLE instance with the provided data.
// It validates the input and returns an error if any field is invalid.
func NewTLE(noradID string, line1 string, line2 string) (TLE, error) {

	tleEpoch, err := xtime.ParseEpoch(line1)
	if err != nil {
		return TLE{}, fmt.Errorf("failed to parse epoch from TLE line: %v", err)
	}

	tle := TLE{
		ID:      uuid.NewString(),
		NoradID: noradID,
		Line1:   line1,
		Line2:   line2,
		Epoch:   tleEpoch,
	}
	if err := tle.Validate(); err != nil {
		return TLE{}, err
	}
	return tle, nil
}

// TLERepository defines the interface for TLE repository operations.
type TLERepository interface {
	FindByNoradID(ctx context.Context, noradID string) ([]TLE, error)
	FindAll(ctx context.Context) ([]TLE, error)
	Save(ctx context.Context, tle TLE) error
	Update(ctx context.Context, tle TLE) error
	Upsert(ctx context.Context, tle TLE) error
	DeleteByNoradID(ctx context.Context, noradID string) error
}

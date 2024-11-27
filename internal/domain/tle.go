package domain

import (
	"errors"
	"time"
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
func NewTLE(noradID, line1, line2 string, epoch time.Time) (*TLE, error) {
	tle := &TLE{
		NoradID: noradID,
		Line1:   line1,
		Line2:   line2,
		Epoch:   epoch,
	}
	if err := tle.Validate(); err != nil {
		return nil, err
	}
	return tle, nil
}

// TLERepository defines the interface for TLE repository operations.
type TLERepository interface {
	FindByNoradID(noradID string) ([]TLE, error)
	FindAll() ([]TLE, error)
	Save(tle TLE) error
	Update(tle TLE) error
	Upsert(tle TLE) error
	Delete(id string) error
}

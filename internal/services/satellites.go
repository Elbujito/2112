package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/space"
)

type SatelliteService struct {
	tleRepo domain.TLERepository // Assuming you have a TLE repository to get TLEs from a database
}

// NewSatelliteService creates a new instance of SatelliteService.
func NewSatelliteService(tleRepo domain.TLERepository) SatelliteService {
	return SatelliteService{tleRepo: tleRepo}
}

// Propagate calculates the satellite's position over the next 24 hours with a specified time step.
// It returns a list of positions at each time step.
// Now, the function takes a NORAD ID and the time step as parameters.
func (s *SatelliteService) Propagate(ctx context.Context, noradID string, timeStep time.Duration) ([]domain.SatellitePosition, error) {
	// Get the TLE data for the satellite by NORAD ID
	tle, err := s.tleRepo.FindByNoradID(ctx, noradID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLE data for NORAD ID %s: %w", noradID, err)
	}

	// Set up the time range: current time and 24 hours later
	startTime := time.Now()
	endTime := startTime.Add(24 * time.Hour)

	// Create a slice to hold the satellite positions
	var positions []domain.SatellitePosition

	// Propagate satellite positions for each time step
	for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.Add(timeStep) {
		// Get the satellite position at the current time
		quadKey, _, err := space.PropagateSatellite(tle[0].Line1, tle[0].Line2, currentTime)
		if err != nil {
			return nil, fmt.Errorf("failed to propagate satellite position for NORAD ID %s at %v: %w", noradID, currentTime, err)
		}

		// Store the position along with the time
		positions = append(positions, domain.SatellitePosition{
			Latitude:  quadKey.Latitude,
			Longitude: quadKey.Longitude,
			Altitude:  float64(quadKey.Level),
			Time:      currentTime,
		})
	}

	return positions, nil
}

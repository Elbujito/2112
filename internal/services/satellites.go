package services

import (
	"context"
	"fmt"
	"time"

	propagator "github.com/Elbujito/2112/internal/clients/propagate"
	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/space"
)

type SatelliteService struct {
	tleRepo         domain.TLERepository // Assuming you have a TLE repository to get TLEs from a database
	propagateClient *propagator.PropagatorClient
}

// NewSatelliteService creates a new instance of SatelliteService.
func NewSatelliteService(tleRepo domain.TLERepository, propagateClient *propagator.PropagatorClient) SatelliteService {
	return SatelliteService{tleRepo: tleRepo, propagateClient: propagateClient}
}

func (s *SatelliteService) Propagate(ctx context.Context, noradID string, duration time.Duration, interval time.Duration) ([]space.SatellitePosition, error) {
	// Validate inputs
	if noradID == "" {
		return nil, fmt.Errorf("NORAD ID is required")
	}
	if duration <= 0 || interval <= 0 {
		return nil, fmt.Errorf("invalid duration or interval: both must be greater than zero")
	}

	// Get the TLE data for the satellite by NORAD ID
	tle, err := s.tleRepo.FindByNoradID(ctx, noradID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLE data for NORAD ID %s: %w", noradID, err)
	}
	if len(tle) == 0 {
		return nil, fmt.Errorf("no TLE data found for NORAD ID %s", noradID)
	}

	// Set up the time range
	startTime := time.Now()
	// endTime := startTime.Add(duration)

	// Use PropagatorClient to propagate satellite positions
	propagatedPositions, err := s.propagateClient.FetchPropagation(ctx, tle[0].Line1, tle[0].Line2, startTime.Format(time.RFC3339), int(duration.Minutes()), int(interval.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch propagated positions for NORAD ID %s: %w", noradID, err)
	}

	// Convert the API response to the internal SatellitePosition format
	var positions []space.SatellitePosition
	for _, pos := range propagatedPositions {
		// Parse the time from the API response
		parsedTime, err := time.Parse(time.RFC3339, pos.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time %s for NORAD ID %s: %w", pos.Time, noradID, err)
		}

		positions = append(positions, space.SatellitePosition{
			Latitude:  pos.Latitude,
			Longitude: pos.Longitude,
			Altitude:  pos.Altitude,
			Time:      parsedTime,
		})
	}

	return positions, nil
}

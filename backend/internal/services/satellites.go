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
	tleRepo          domain.TLERepository // Assuming you have a TLE repository to get TLEs from a database
	propagateClient  *propagator.PropagatorClient
	celestrackClient celestrackClient
	repo             domain.SatelliteRepository
}

// NewSatelliteService creates a new instance of SatelliteService.
func NewSatelliteService(tleRepo domain.TLERepository, propagateClient *propagator.PropagatorClient, celestrackClient celestrackClient, repo domain.SatelliteRepository) SatelliteService {
	return SatelliteService{tleRepo: tleRepo, propagateClient: propagateClient, celestrackClient: celestrackClient, repo: repo}
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

// GetSatelliteByNoradID retrieves a satellite by NORAD ID.
func (s *SatelliteService) GetSatelliteByNoradID(ctx context.Context, noradID string) (domain.Satellite, error) {
	return s.repo.FindByNoradID(ctx, noradID)
}

// ListAllSatellites retrieves all stored satellites.
func (s *SatelliteService) ListAllSatellites(ctx context.Context) ([]domain.Satellite, error) {
	return s.repo.FindAll(ctx)
}
func (s *SatelliteService) FetchAndStoreAllSatellites(ctx context.Context) ([]domain.Satellite, error) {
	// Fetch all satellite metadata from CelestrackClient
	rawSatellites, err := s.celestrackClient.FetchSatelliteMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch satellite metadata: %w", err)
	}

	if len(rawSatellites) == 0 {
		return nil, fmt.Errorf("no satellite metadata available")
	}

	var storedSatellites []domain.Satellite
	for _, rawSatellite := range rawSatellites {

		// Use the updated constructor to create a Satellite
		satellite, err := domain.NewSatelliteFromStatCat(
			rawSatellite.Name,
			rawSatellite.NoradID,
			domain.Other, // Default type; adjust based on metadata if available
			&rawSatellite.LaunchDate,
			rawSatellite.DecayDate,
			rawSatellite.IntlDesignator,
			rawSatellite.Owner,
			rawSatellite.ObjectType,
			rawSatellite.Period,
			rawSatellite.Inclination,
			rawSatellite.Apogee,
			rawSatellite.Apogee,
			rawSatellite.RCS,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create satellite for NORAD ID %s: %w", rawSatellite.NoradID, err)
		}
		// Add the satellite to the result list
		storedSatellites = append(storedSatellites, satellite)
	}

	// Save the satellite to the repository
	if err := s.repo.SaveBatch(ctx, storedSatellites); err != nil {
		return nil, fmt.Errorf("failed to save satellite to database: %w", err)
	}

	return storedSatellites, nil
}

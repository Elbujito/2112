package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/polygon"
	"github.com/Elbujito/2112/pkg/fx/space"
)

type ComputeVisibilitiesHandler struct {
	tileRepo       domain.TileRepository
	tleRepo        domain.TLERepository
	satelliteRepo  domain.SatelliteRepository
	visibilityRepo domain.TileSatelliteMappingRepository
}

func NewComputeVisibilitiesHandler(
	tileRepo domain.TileRepository,
	tleRepo domain.TLERepository,
	satelliteRepo domain.SatelliteRepository,
	visibilityRepo domain.TileSatelliteMappingRepository,
) ComputeVisibilitiesHandler {
	return ComputeVisibilitiesHandler{
		tileRepo:       tileRepo,
		tleRepo:        tleRepo,
		satelliteRepo:  satelliteRepo,
		visibilityRepo: visibilityRepo,
	}
}

func (h *ComputeVisibilitiesHandler) GetTask() Task {
	return Task{
		Name:         "execComputeVisibilitiesTask",
		Description:  "Computes satellite visibilities for all tiles",
		RequiredArgs: []string{},
	}
}

func (h *ComputeVisibilitiesHandler) Run(ctx context.Context, args map[string]string) error {
	// Fetch all satellites, TLEs, and tiles
	satellites, err := h.satelliteRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch satellites: %w", err)
	}

	tles, err := h.tleRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch TLEs: %w", err)
	}
	tleMap := make(map[string]domain.TLE)
	for _, tle := range tles {
		tleMap[tle.NoradID] = tle
	}

	tiles, err := h.tileRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch tiles: %w", err)
	}

	startTime := time.Now()
	endTime := startTime.Add(24 * time.Hour)

	for _, sat := range satellites {
		err := h.computeSatelliteVisibility(ctx, sat, tleMap, tiles, startTime, endTime)
		if err != nil {
			return err
		}
	}

	return nil
}

// Compute visibility for a single satellite
func (h *ComputeVisibilitiesHandler) computeSatelliteVisibility(
	ctx context.Context,
	sat domain.Satellite,
	tleMap map[string]domain.TLE,
	tiles []domain.Tile,
	startTime, endTime time.Time,
) error {
	tle, ok := tleMap[sat.NoradID]
	if !ok {
		return fmt.Errorf("no TLE data found for satellite %s", sat.NoradID)
	}

	const timeStep = 1 * time.Hour
	visibilityBatch := make([]domain.TileSatelliteMapping, 0, len(tiles))

	for t := startTime; t.Before(endTime); t = t.Add(timeStep) {
		for _, tile := range tiles {

			aos, maxElevation := space.ComputeVisibilityWindow(
				tle.NoradID, tle.Line1, tle.Line2,
				polygon.Point{Latitude: tile.CenterLat, Longitude: tile.CenterLon}, tile.Radius, t, endTime, timeStep,
			)

			if !aos.IsZero() {
				visibility := domain.NewVisibility(
					sat.NoradID,
					tile.ID,
					aos,
					maxElevation,
				)
				visibilityBatch = append(visibilityBatch, visibility)
			}

			// Save in batches
			if len(visibilityBatch) >= 100 {
				if err := h.visibilityRepo.SaveBatch(ctx, visibilityBatch); err != nil {
					log.Printf("Failed to save batch for satellite %s: %v\n", sat.NoradID, err)
				}
				visibilityBatch = visibilityBatch[:0] // Reset batch
			}
		}
	}

	// Save remaining visibilities in batch
	if len(visibilityBatch) > 0 {
		if err := h.visibilityRepo.SaveBatch(ctx, visibilityBatch); err != nil {
			log.Printf("Failed to save remaining batch for satellite %s: %v\n", sat.NoradID, err)
		}
	}

	return nil
}

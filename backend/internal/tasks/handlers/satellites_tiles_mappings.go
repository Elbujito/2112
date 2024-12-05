package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/polygon"
	"github.com/Elbujito/2112/pkg/fx/space"
)

type SatellitesTilesMappingsHandler struct {
	tileRepo       domain.TileRepository
	tleRepo        domain.TLERepository
	satelliteRepo  domain.SatelliteRepository
	visibilityRepo domain.MappingRepository
}

func NewSatellitesTilesMappingsHandler(
	tileRepo domain.TileRepository,
	tleRepo domain.TLERepository,
	satelliteRepo domain.SatelliteRepository,
	visibilityRepo domain.MappingRepository,
) SatellitesTilesMappingsHandler {
	return SatellitesTilesMappingsHandler{
		tileRepo:       tileRepo,
		tleRepo:        tleRepo,
		satelliteRepo:  satelliteRepo,
		visibilityRepo: visibilityRepo,
	}
}

func (h *SatellitesTilesMappingsHandler) GetTask() Task {
	return Task{
		Name:         "satellites_tiles_mappings",
		Description:  "Computes satellite visibilities for all tiles",
		RequiredArgs: []string{"timeStepInSeconds", "periodInMinutes"},
	}
}

func (h *SatellitesTilesMappingsHandler) Run(ctx context.Context, args map[string]string) error {

	argTimeStep, ok := args["timeStepInSeconds"]
	if !ok || argTimeStep == "" {
		return fmt.Errorf("missing required argument: timeStepInSeconds")
	}

	timeStepInSeconds, err := strconv.Atoi(argTimeStep)
	if err != nil {
		return err
	}
	timeStepDuration := time.Duration(timeStepInSeconds) * time.Second

	argPeriod, ok := args["periodInMinutes"]
	if !ok || argTimeStep == "" {
		return fmt.Errorf("missing required argument: periodInMinutes")
	}

	periodInMinutes, err := strconv.Atoi(argPeriod)
	if err != nil {
		return err
	}
	periodDuration := time.Duration(periodInMinutes) * time.Minute

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
	endTime := startTime.Add(periodDuration)

	for _, sat := range satellites {
		err := h.computeSatellitesTilesMappings(ctx, sat, tleMap, tiles, startTime, endTime, timeStepDuration)
		if err != nil {
			continue
		}
	}

	return nil
}

// Compute visibility for a single satellite
func (h *SatellitesTilesMappingsHandler) computeSatellitesTilesMappings(
	ctx context.Context,
	sat domain.Satellite,
	tleMap map[string]domain.TLE,
	tiles []domain.Tile,
	startTime, endTime time.Time,
	timeStepDuration time.Duration,
) error {

	tle, ok := tleMap[sat.NoradID]
	if !ok {
		return fmt.Errorf("no TLE data found for satellite %s", sat.NoradID)
	}

	visibilityBatch := make([]domain.TileSatelliteMapping, 0, len(tiles))

	for t := startTime; t.Before(endTime); t = t.Add(timeStepDuration) {
		for _, tile := range tiles {

			aos, maxElevation := space.ComputeVisibilityWindow(
				tle.NoradID, tle.Line1, tle.Line2,
				polygon.Point{Latitude: tile.CenterLat, Longitude: tile.CenterLon}, tile.Radius, t, endTime, timeStepDuration,
			)

			if !aos.IsZero() {
				visibility := domain.NewMapping(
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

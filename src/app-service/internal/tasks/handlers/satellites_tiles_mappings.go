package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xpolygon"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xspace"
)

type SatellitesTilesMappingsHandler struct {
	tileRepo       domain.TileRepository
	tleRepo        repository.TleRepository
	satelliteRepo  domain.SatelliteRepository
	visibilityRepo domain.MappingRepository
}

func NewSatellitesTilesMappingsHandler(
	tileRepo domain.TileRepository,
	tleRepo repository.TleRepository,
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
	// Parse period argument
	periodInMinutes, err := ParseIntArg(args, "periodInMinutes")
	if err != nil {
		return fmt.Errorf("missing or invalid argument 'periodInMinutes': %w", err)
	}
	periodDuration := time.Duration(periodInMinutes) * time.Minute

	// Parse timestep argument (optional)
	timeStepDuration := time.Duration(0)
	if argTimeStep, ok := args["timeStepInSeconds"]; ok && argTimeStep != "" {
		timeStepInSeconds, err := strconv.Atoi(argTimeStep)
		if err != nil {
			return fmt.Errorf("invalid 'timeStepInSeconds' argument: %w", err)
		}
		timeStepDuration = time.Duration(timeStepInSeconds) * time.Second
	}

	// Fetch satellites, TLEs, and tiles
	satellites, err := h.satelliteRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch satellites: %w", err)
	}

	// tles, err := h.tleRepo.FindAll(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch TLEs: %w", err)
	// }
	tles := []domain.TLE{}
	tleMap := mapTLEsByNoradID(tles)

	tiles, err := h.tileRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch tiles: %w", err)
	}

	// Compute mappings
	startTime := time.Now()
	endTime := startTime.Add(periodDuration)

	for _, sat := range satellites {
		err := h.computeSatellitesTilesMappings(ctx, sat, tleMap, tiles, startTime, endTime, timeStepDuration)
		if err != nil {
			log.Printf("Error computing mappings for satellite %s: %v", sat.NoradID, err)
			continue
		}
	}

	return nil
}

func mapTLEsByNoradID(tles []domain.TLE) map[string]domain.TLE {
	tleMap := make(map[string]domain.TLE, len(tles))
	for _, tle := range tles {
		tleMap[tle.NoradID] = tle
	}
	return tleMap
}

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

	// Dynamically calculate timestep if not provided
	if timeStepDuration == 0 && sat.Altitude != nil {
		timeStepDuration = xspace.CalculateOptimalTimestep(*sat.Altitude, tiles[0].Radius)
	} else {
		timeStepDuration = time.Hour * 1
	}

	visibilityBatch := make([]domain.TileSatelliteMapping, 0, len(tiles))
	for t := startTime; t.Before(endTime); t = t.Add(timeStepDuration) {
		for _, tile := range tiles {
			aos, maxElevation := xspace.ComputeVisibilityWindow(
				tle.NoradID, tle.Line1, tle.Line2,
				xpolygon.Point{Latitude: tile.CenterLat, Longitude: tile.CenterLon}, tile.Radius, t, endTime, timeStepDuration,
			)

			if !aos.IsZero() {
				visibilityBatch = append(visibilityBatch, domain.NewMapping(
					sat.NoradID, tile.ID, aos, maxElevation,
				))
			}

			// Save in batches
			if len(visibilityBatch) >= 100 {
				h.saveVisibilityBatch(ctx, visibilityBatch, sat.NoradID)
				visibilityBatch = visibilityBatch[:0] // Reset batch
			}
		}
	}

	// Save remaining visibilities
	if len(visibilityBatch) > 0 {
		h.saveVisibilityBatch(ctx, visibilityBatch, sat.NoradID)
	}

	return nil
}

func (h *SatellitesTilesMappingsHandler) saveVisibilityBatch(ctx context.Context, batch []domain.TileSatelliteMapping, noradID string) {
	if err := h.visibilityRepo.SaveBatch(ctx, batch); err != nil {
		log.Printf("Failed to save batch for satellite %s: %v", noradID, err)
	}
}

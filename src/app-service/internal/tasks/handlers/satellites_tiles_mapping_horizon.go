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

type SatellitesTilesMappingsByHorizonHandler struct {
	tileRepo      domain.TileRepository
	tleRepo       repository.TleRepository
	satelliteRepo domain.SatelliteRepository
	mappingRepo   domain.MappingRepository
}

func NewSatellitesTilesMappingsByHorizonHandler(
	tileRepo domain.TileRepository,
	tleRepo repository.TleRepository,
	satelliteRepo domain.SatelliteRepository,
	visibilityRepo domain.MappingRepository,
) SatellitesTilesMappingsByHorizonHandler {
	return SatellitesTilesMappingsByHorizonHandler{
		tileRepo:      tileRepo,
		tleRepo:       tleRepo,
		satelliteRepo: satelliteRepo,
		mappingRepo:   visibilityRepo,
	}
}

// GetTask returns the task metadata
func (h *SatellitesTilesMappingsByHorizonHandler) GetTask() Task {
	return Task{
		Name:         "satellites_tiles_mapping_horizon",
		Description:  "Computes satellite visibilities for all tiles",
		RequiredArgs: []string{"timeStepInSeconds", "periodInMinutes"},
	}
}

// Run executes the visibility computation process
func (h *SatellitesTilesMappingsByHorizonHandler) Run(ctx context.Context, args map[string]string) error {

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

	// tles, err := h.tleRepo.FindAll(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch TLEs: %w", err)
	// }
	tles := []domain.TLE{}
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

	// Group tiles by region (e.g., by latitude/longitude or zoom level)
	tileGroups := groupTilesByRegion(tiles)

	// For each satellite, compute visibility for the grouped tiles
	for _, sat := range satellites {
		err := h.computeMappings(ctx, sat, tleMap, tileGroups, startTime, endTime, timeStepDuration)
		if err != nil {
			return err
		}
	}

	return nil
}

// Compute visibility for a single satellite, optimized with tile grouping and satellite horizon.
func (h *SatellitesTilesMappingsByHorizonHandler) computeMappings(
	ctx context.Context,
	sat domain.Satellite,
	tleMap map[string]domain.TLE,
	tileGroups map[string][]domain.Tile,
	startTime, endTime time.Time,
	timeStep time.Duration,
) error {
	tle, ok := tleMap[sat.NoradID]
	if !ok {
		return fmt.Errorf("no TLE data found for satellite %s", sat.NoradID)
	}

	visibilityBatch := make([]domain.TileSatelliteMapping, 0, 100)

	// Iterate over time steps
	for t := startTime; t.Before(endTime); t = t.Add(timeStep) {
		// Compute the satellite's horizon at the current time step
		visibleRegion, err := xspace.ComputeSatelliteHorizon(t, tle.Line1, tle.Line2)
		if err != nil {
			return fmt.Errorf("failed to compute satellite horizon: %w", err)
		}

		// Process each group of tiles (based on region grouping)
		for _, regionTiles := range tileGroups {
			// For each tile, check if it falls within the satellite's visible region (horizon)
			for _, tile := range regionTiles {
				if len(tile.Vertices) == 0 {
					log.Printf("Skipping tile %s due to invalid polygon data\n", tile.ID)
					continue
				}

				// Check if the tile's center is within the satellite's visible region (horizon)
				if isTileVisibleFromRegion(tile, visibleRegion) {
					// Calculate visibility for the tile
					aos, maxElevation := xspace.ComputeVisibilityWindow(
						tle.NoradID, tle.Line1, tle.Line2,
						xpolygon.Point{Latitude: tile.CenterLat, Longitude: tile.CenterLon},
						tile.Radius, t, endTime, timeStep,
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
				}
			}
		}

		// Save in batches
		if len(visibilityBatch) >= 100 {
			if err := h.mappingRepo.SaveBatch(ctx, visibilityBatch); err != nil {
				log.Printf("Failed to save batch for satellite %s: %v\n", sat.NoradID, err)
			}
			visibilityBatch = visibilityBatch[:0] // Reset batch
		}
	}

	// Save any remaining visibilities in batch
	if len(visibilityBatch) > 0 {
		if err := h.mappingRepo.SaveBatch(ctx, visibilityBatch); err != nil {
			log.Printf("Failed to save remaining batch for satellite %s: %v\n", sat.NoradID, err)
		}
	}

	return nil
}

// Group tiles by region (e.g., zoom level, latitude range)
func groupTilesByRegion(tiles []domain.Tile) map[string][]domain.Tile {
	tileGroups := make(map[string][]domain.Tile)

	for _, tile := range tiles {
		regionKey := fmt.Sprintf("zoom%d_lat%.2f", tile.ZoomLevel, tile.CenterLat)
		tileGroups[regionKey] = append(tileGroups[regionKey], tile)
	}

	return tileGroups
}

// Check if a tile is visible from a given region
func isTileVisibleFromRegion(tile domain.Tile, visibleRegion []xpolygon.Point) bool {
	return xpolygon.IsPointInPolygon(xpolygon.Point{Latitude: tile.CenterLat, Longitude: tile.CenterLon}, visibleRegion)
}

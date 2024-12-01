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

type VisbilityBySatelliteHorizonHandler struct {
	tileRepo       domain.TileRepository
	tleRepo        domain.TLERepository
	satelliteRepo  domain.SatelliteRepository
	visibilityRepo domain.TileSatelliteMappingRepository
}

func NewVisbilityBySatelliteHorizonHandler(
	tileRepo domain.TileRepository,
	tleRepo domain.TLERepository,
	satelliteRepo domain.SatelliteRepository,
	visibilityRepo domain.TileSatelliteMappingRepository,
) VisbilityBySatelliteHorizonHandler {
	return VisbilityBySatelliteHorizonHandler{
		tileRepo:       tileRepo,
		tleRepo:        tleRepo,
		satelliteRepo:  satelliteRepo,
		visibilityRepo: visibilityRepo,
	}
}

// GetTask returns the task metadata
func (h *VisbilityBySatelliteHorizonHandler) GetTask() Task {
	return Task{
		Name:         "execVisbilityBySatelliteHorizonTask",
		Description:  "Computes satellite visibilities for all tiles",
		RequiredArgs: []string{},
	}
}

// Run executes the visibility computation process
func (h *VisbilityBySatelliteHorizonHandler) Run(ctx context.Context, args map[string]string) error {
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

	// Group tiles by region (e.g., by latitude/longitude or zoom level)
	tileGroups := groupTilesByRegion(tiles)

	// For each satellite, compute visibility for the grouped tiles
	for _, sat := range satellites {
		err := h.computeSatelliteVisibility(ctx, sat, tleMap, tileGroups, startTime, endTime)
		if err != nil {
			return err
		}
	}

	return nil
}

// Compute visibility for a single satellite, optimized with tile grouping and satellite horizon.
func (h *VisbilityBySatelliteHorizonHandler) computeSatelliteVisibility(
	ctx context.Context,
	sat domain.Satellite,
	tleMap map[string]domain.TLE,
	tileGroups map[string][]domain.Tile,
	startTime, endTime time.Time,
) error {
	tle, ok := tleMap[sat.NoradID]
	if !ok {
		return fmt.Errorf("no TLE data found for satellite %s", sat.NoradID)
	}

	const timeStep = 1 * time.Hour
	visibilityBatch := make([]domain.TileSatelliteMapping, 0, 100)

	// Iterate over time steps
	for t := startTime; t.Before(endTime); t = t.Add(timeStep) {
		// Compute the satellite's horizon at the current time step
		visibleRegion, err := space.ComputeSatelliteHorizon(t, tle)
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
					aos, maxElevation := space.ComputeVisibilityWindow(
						tle.NoradID, tle.Line1, tle.Line2,
						polygon.Point{Latitude: tile.CenterLat, Longitude: tile.CenterLon},
						tile.Radius, t, endTime, timeStep,
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
				}
			}
		}

		// Save in batches
		if len(visibilityBatch) >= 100 {
			if err := h.visibilityRepo.SaveBatch(ctx, visibilityBatch); err != nil {
				log.Printf("Failed to save batch for satellite %s: %v\n", sat.NoradID, err)
			}
			visibilityBatch = visibilityBatch[:0] // Reset batch
		}
	}

	// Save any remaining visibilities in batch
	if len(visibilityBatch) > 0 {
		if err := h.visibilityRepo.SaveBatch(ctx, visibilityBatch); err != nil {
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
func isTileVisibleFromRegion(tile domain.Tile, visibleRegion []polygon.Point) bool {
	return polygon.IsPointInPolygon(polygon.Point{Latitude: tile.CenterLat, Longitude: tile.CenterLon}, visibleRegion)
}

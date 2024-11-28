package handlers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/space"
	"github.com/joshuaferrara/go-satellite"
)

type ComputeVisibilitiesHandler struct {
	tileRepo       domain.TileRepository
	tleRepo        domain.TLERepository
	satelliteRepo  domain.SatelliteRepository
	visibilityRepo domain.VisibilityRepository
}

func NewComputeVisibilitiesHandler(
	tileRepo domain.TileRepository,
	tleRepo domain.TLERepository,
	satelliteRepo domain.SatelliteRepository,
	visibilityRepo domain.VisibilityRepository,
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

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, sat := range satellites {
		wg.Add(1)
		go func(sat domain.Satellite) {
			defer wg.Done()
			err := h.computeSatelliteVisibility(ctx, sat, tleMap, tiles, startTime, endTime, &mutex)
			if err != nil {
				fmt.Printf("Error computing visibility for satellite %s: %v\n", sat.NoradID, err)
			}
		}(sat)
	}

	wg.Wait()

	return nil
}

// Compute visibility for a single satellite
func (h *ComputeVisibilitiesHandler) computeSatelliteVisibility(
	ctx context.Context,
	sat domain.Satellite,
	tleMap map[string]domain.TLE,
	tiles []domain.Tile,
	startTime, endTime time.Time,
	mutex *sync.Mutex,
) error {
	tle, ok := tleMap[sat.NoradID]
	if !ok {
		return fmt.Errorf("no TLE data found for satellite %s", sat.NoradID)
	}

	satrec := satellite.TLEToSat(tle.Line1, tle.Line2, satellite.GravityWGS84)

	const timeStep = time.Minute
	for t := startTime; t.Before(endTime); t = t.Add(timeStep) {
		year, month, day := t.Date()
		hour, minute, second := t.Clock()

		position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)

		// Calculate GST for ECI to LLA conversion
		gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)

		// Convert ECI to Geodetic (lat, lon, alt)
		altitude, _, geoPosition := satellite.ECIToLLA(position, gmst)
		lat, lon := geoPosition.Latitude, geoPosition.Longitude

		for _, tile := range tiles {
			elevation := space.CalculateElevation(lat, lon, altitude, tile.CenterLat, tile.CenterLon)
			if elevation > 0 { // Satellite is visible
				aos := t
				los := space.ComputeLOS(satrec, tile.CenterLat, tile.CenterLon, t, endTime, timeStep)
				maxElevation := space.CalculateMaxElevation(lat, lon, altitude, tile.CenterLat, tile.CenterLon)
				visibility := domain.NewVisibility(
					sat.NoradID,
					tile.ID,
					aos,
					los,
					maxElevation,
				)

				mutex.Lock()
				err := h.visibilityRepo.Save(ctx, visibility)
				mutex.Unlock()
				if err != nil {
					fmt.Printf("Failed to save visibility record for satellite %s: %v\n", sat.NoradID, err)
				}
			}
		}
	}
	return nil
}

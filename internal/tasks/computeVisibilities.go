package tasks

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Elbujito/2112/internal/data/models"
	"github.com/joshuaferrara/go-satellite"
)

func init() {
	task := &Task{
		Name:         "execComputeVisibilitiesTask",
		Description:  "Computes satellite visibilities for all tiles",
		RequiredArgs: []string{},
		Run:          execComputeVisibilitiesTask,
	}
	Tasks.AddTask(task)
}

// Task execution function
func execComputeVisibilitiesTask(env *TaskEnv, args map[string]string) error {
	// Initialize service instances
	satelliteService := models.SatelliteModel()
	tleService := models.TLEModel()
	tileService := models.TileModel()

	// Step 1: Fetch all satellites
	satellites, err := satelliteService.FindAll()
	if err != nil {
		return fmt.Errorf("failed to fetch satellites: %w", err)
	}

	// Step 2: Fetch all TLEs and store in a map for quick access
	tles, err := tleService.FindAll()
	if err != nil {
		return fmt.Errorf("failed to fetch TLEs: %w", err)
	}
	tleMap := make(map[string]*models.TLE)
	for _, tle := range tles {
		tleMap[tle.NoradID] = tle
	}

	// Step 3: Fetch all tiles
	tiles, err := tileService.FindAll()
	if err != nil {
		return fmt.Errorf("failed to fetch tiles: %w", err)
	}

	// Define time range for visibility calculation
	startTime := time.Now()
	endTime := startTime.Add(24 * time.Hour)

	// Use WaitGroup for concurrency
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Process each satellite in parallel
	for _, sat := range satellites {
		wg.Add(1)
		go func(sat *models.Satellite) {
			defer wg.Done()
			if err := computeSatelliteVisibility(sat, tleMap, tiles, startTime, endTime, &mutex); err != nil {
				fmt.Printf("Error computing visibility for satellite %s: %v\n", sat.NoradID, err)
			}
		}(sat)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return nil
}

// Compute visibility for a single satellite
func computeSatelliteVisibility(
	satModel *models.Satellite,
	tleMap map[string]*models.TLE,
	tiles []*models.Tile,
	startTime, endTime time.Time,
	mutex *sync.Mutex,
) error {
	// Fetch TLE for the satellite from the map
	tle, ok := tleMap[satModel.NoradID]
	if !ok {
		return fmt.Errorf("no TLE data found for satellite %s", satModel.NoradID)
	}

	// Parse TLE and initialize satellite model
	satrec := satellite.TLEToSat(tle.Line1, tle.Line2, satellite.GravityWGS84)

	// Iterate over the time range
	const timeStep = time.Minute
	for t := startTime; t.Before(endTime); t = t.Add(timeStep) {
		// Extract date and time components
		year, month, day := t.Date()
		hour, minute, second := t.Clock()

		// Propagate the satellite's position and velocity
		position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)

		// Calculate GST for ECI to LLA conversion
		gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)

		// Convert ECI to Geodetic (lat, lon, alt)
		altitude, _, geoPosition := satellite.ECIToLLA(position, gmst)
		lat, lon := geoPosition.Latitude, geoPosition.Longitude

		// Check visibility for each tile
		for _, tile := range tiles {
			elevation := calculateElevation(lat, lon, altitude, tile)

			// AOS and LOS logic
			if elevation > 0 { // Satellite is visible
				aos := t
				los := computeLOS(satrec, tile, t, endTime, timeStep)

				// Calculate max elevation
				maxElevation := calculateMaxElevation(lat, lon, altitude, tile)

				// Create a new visibility record
				visibility := &models.Visibility{
					NoradID:      satModel.NoradID,
					TileID:       tile.ID,
					StartTime:    aos,
					EndTime:      los,
					MaxElevation: maxElevation,
				}

				// Protect visibility creation with a mutex
				mutex.Lock()
				func() {
					defer mutex.Unlock()
					err := visibility.Create()
					if err != nil {
						fmt.Printf("Failed to save visibility record for satellite %s: %v\n", satModel.NoradID, err)
					}
				}()
			}
		}
	}

	return nil
}

// Calculate LOS dynamically by continuing propagation until the satellite is no longer visible
func computeLOS(satrec satellite.Satellite, tile *models.Tile, startTime, endTime time.Time, timeStep time.Duration) time.Time {
	for t := startTime.Add(timeStep); t.Before(endTime); t = t.Add(timeStep) {
		year, month, day := t.Date()
		hour, minute, second := t.Clock()

		// Propagate satellite position
		position, _ := satellite.Propagate(satrec, year, int(month), day, hour, minute, second)

		// Calculate GST for ECI to LLA conversion
		gmst := satellite.GSTimeFromDate(year, int(month), day, hour, minute, second)

		// Convert ECI to Geodetic (lat, lon, alt)
		altitude, _, geoPosition := satellite.ECIToLLA(position, gmst)
		lat, lon := geoPosition.Latitude, geoPosition.Longitude

		// Check elevation
		elevation := calculateElevation(lat, lon, altitude, tile)
		if elevation <= 0 {
			// Satellite is no longer visible
			return t
		}
	}

	// Return endTime if LOS is not found
	return endTime
}

// Calculate the elevation angle of the satellite from the tile center
func calculateElevation(satLat, satLon, satAlt float64, tile *models.Tile) float64 {
	// Simplified elevation calculation:
	// Use haversine distance and consider the altitude difference for angular elevation
	dist := haversineDistance(satLat, satLon, tile.CenterLat, tile.CenterLon)
	return 90.0 - dist/10.0
}

// Calculate max elevation
func calculateMaxElevation(satLat, satLon, satAlt float64, tile *models.Tile) float64 {
	return 90.0 - haversineDistance(satLat, satLon, tile.CenterLat, tile.CenterLon)/10.0
}

// Haversine formula for distance calculation
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	lat1 = degreesToRadians(lat1)
	lat2 = degreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// Convert degrees to radians
func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

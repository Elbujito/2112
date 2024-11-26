package tasks

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Elbujito/2112/pkg/api/handlers/celestrack"
	"github.com/Elbujito/2112/pkg/api/mappers"
	"github.com/Elbujito/2112/pkg/db/models"
	xtime "github.com/Elbujito/2112/pkg/utils/time"
)

func init() {
	task := &Task{
		Name:         "fetchAndUpsertTLE",
		Description:  "Fetch TLE from CelesTrak and upsert it in the database",
		RequiredArgs: []string{"category"},
		Run:          execFetchCatalogTLETask,
	}
	Tasks.AddTask(task)
}

func execFetchCatalogTLETask(env *TaskEnv, args map[string]string) error {
	category, ok := args["category"]
	if !ok || category == "" {
		return fmt.Errorf("missing required argument: category")
	}

	// Use SatelliteService and TLEService for dependency injection
	satelliteService := models.SatelliteModel()
	tleService := models.TLEModel()

	// Fetch and upsert TLEs
	return fetchAndUpsertTLEs(category, satelliteService, tleService, celestrack.FetchCategoryTLE)
}

func fetchAndUpsertTLEs(category string, satelliteService models.SatelliteService, tleService models.TLEService, fetchTLEHandler func(string) ([]*mappers.RawTLE, error)) error {
	// Fetch TLEs from the provided category
	tles, err := fetchTLEHandler(category)
	if err != nil {
		return fmt.Errorf("failed to fetch TLE catalog for category %s: %v", category, err)
	}

	// Upsert each TLE
	for _, tle := range tles {
		if err := upsertTLE(tle, satelliteService, tleService); err != nil {
			log.Printf("Failed to upsert TLE for NORAD ID %s: %v", tle.NoradID, err)
		} else {
			log.Printf("Successfully upserted TLE for NORAD ID %s", tle.NoradID)
		}
	}

	return nil
}

func upsertTLE(rawTLE *mappers.RawTLE, satelliteService models.SatelliteService, tleService models.TLEService) error {
	// Ensure the satellite exists
	_, err := ensureSatelliteExists(rawTLE.NoradID, satelliteService)
	if err != nil {
		return fmt.Errorf("failed to ensure satellite existence: %v", err)
	}

	// Check for existing TLEs
	existingTLEs, err := tleService.FindByNoradID(rawTLE.NoradID)
	if err == nil && len(existingTLEs) > 0 {
		// Update existing TLE
		existingTLE := existingTLEs[0]
		existingTLE.Line1 = rawTLE.Line1
		existingTLE.Line2 = rawTLE.Line2
		existingTLE.Epoch, err = parseEpoch(rawTLE.Line1)
		if err != nil {
			return fmt.Errorf("failed to parse epoch from TLE line: %v", err)
		}
		return tleService.Update(existingTLE)
	}

	// Insert a new TLE
	newEpoch, err := parseEpoch(rawTLE.Line1)
	if err != nil {
		return fmt.Errorf("failed to parse epoch from TLE line: %v", err)
	}
	newTLE := &models.TLE{
		NoradID: rawTLE.NoradID,
		Line1:   rawTLE.Line1,
		Line2:   rawTLE.Line2,
		Epoch:   newEpoch,
	}
	return newTLE.Create()
}

func ensureSatelliteExists(noradID string, satelliteService models.SatelliteService) (string, error) {
	// Try to find the satellite by NORAD ID
	satellite, err := satelliteService.FindByNoradID(noradID)
	if err == nil && satellite != nil && satellite.NoradID == noradID {
		return satellite.NoradID, err
	}

	// Satellite not found; create a new one
	newSatellite := &models.Satellite{
		Name:    fmt.Sprintf("Unknown Satellite %s", noradID),
		NoradID: noradID,
	}
	if err := newSatellite.Create(); err != nil {
		log.Printf("Failed to create satellite: %v %s", noradID, err)
	} else {
		log.Printf("Successfully inserted satellite for NORAD ID %s", noradID)
	}
	return newSatellite.ID, err
}

func parseEpoch(line1 string) (time.Time, error) {
	// Extract epoch substring from the TLE line
	if len(line1) < 32 {
		return time.Time{}, fmt.Errorf("invalid TLE line: epoch data missing")
	}
	epochStr := strings.TrimSpace(line1[18:32])
	return xtime.FromRawTLE(epochStr)
}

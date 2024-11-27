package tasks

import (
	"fmt"

	"github.com/Elbujito/2112/internal/api/handlers/celestrack"
	"github.com/Elbujito/2112/internal/api/mappers"
	"github.com/Elbujito/2112/internal/data/models"
	xtime "github.com/Elbujito/2112/pkg/fx/time"
)

type TLEHandler struct {
	SatelliteService models.SatelliteService
	TLEService       models.TLEService
	FetchTLEHandler  func(string) ([]*mappers.RawTLE, error)
}

func NewTLEHandler(
	satelliteService models.SatelliteService,
	tleService models.TLEService,
	fetchTLEHandler func(string) ([]*mappers.RawTLE, error),
) *TLEHandler {
	return &TLEHandler{
		SatelliteService: satelliteService,
		TLEService:       tleService,
		FetchTLEHandler:  fetchTLEHandler,
	}
}

func (h *TLEHandler) Run(category string) error {
	tles, err := h.FetchTLEHandler(category)
	if err != nil {
		return fmt.Errorf("failed to fetch TLE catalog for category %s: %v", category, err)
	}

	for _, rawTLE := range tles {
		err := h.upsertTLE(rawTLE)
		if err != nil {
			return fmt.Errorf("failed to upsert TLE for NORAD ID %s: %v", rawTLE.NoradID, err)
		}
	}

	return nil
}

func (h *TLEHandler) upsertTLE(rawTLE *mappers.RawTLE) error {
	_, err := h.ensureSatelliteExists(rawTLE.NoradID)
	if err != nil {
		return fmt.Errorf("failed to ensure satellite existence: %v", err)
	}

	tleEpoch, err := xtime.ParseEpoch(rawTLE.Line1)
	if err != nil {
		return fmt.Errorf("failed to parse epoch from TLE line: %v", err)
	}

	tle := &models.TLE{
		NoradID: rawTLE.NoradID,
		Line1:   rawTLE.Line1,
		Line2:   rawTLE.Line2,
		Epoch:   tleEpoch,
	}
	return h.TLEService.Upsert(tle)
}

func (h *TLEHandler) ensureSatelliteExists(noradID string) (string, error) {
	satellite, err := h.SatelliteService.FindByNoradID(noradID)
	if err == nil && satellite != nil && satellite.NoradID == noradID {
		return satellite.ID, nil
	}

	newSatellite := &models.Satellite{
		Name:    fmt.Sprintf("Unknown Satellite %s", noradID),
		NoradID: noradID,
	}

	err = h.SatelliteService.Save(newSatellite)
	if err != nil {
		return "", fmt.Errorf("failed to create satellite: %v", err)
	}

	return newSatellite.ID, nil
}

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
	handler := NewTLEHandler(
		models.SatelliteModel(),
		models.TLEModel(),
		celestrack.FetchCategoryTLE,
	)
	return handler.Run(category)
}

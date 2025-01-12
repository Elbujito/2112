package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
)

type TleServiceClient interface {
	FetchTLEFromSatCatByCategory(ctx context.Context, category string, contextID string) ([]domain.TLE, error)
}

type CelestrackTleUploadHandler struct {
	satelliteRepo domain.SatelliteRepository
	tleRepo       repository.TleRepository
	tleService    TleServiceClient
}

func NewCelestrackTleUploadHandler(
	satelliteRepo domain.SatelliteRepository,
	tleRepo repository.TleRepository,
	tleService TleServiceClient) CelestrackTleUploadHandler {
	return CelestrackTleUploadHandler{
		satelliteRepo: satelliteRepo,
		tleRepo:       tleRepo,
		tleService:    tleService,
	}
}

func (h *CelestrackTleUploadHandler) GetTask() Task {
	return Task{
		Name:         "celestrack_tle_upload",
		Description:  "Fetch TLE from CelesTrak and upsert it in the database",
		RequiredArgs: []string{"category", "maxCount", "contextID"},
	}
}

func (h *CelestrackTleUploadHandler) Run(ctx context.Context, args map[string]string) error {

	category, ok := args["category"]
	if !ok || category == "" {
		return fmt.Errorf("missing required argument: category")
	}

	nbTles, ok := args["maxCount"]
	if !ok || nbTles == "" {
		return fmt.Errorf("missing required argument: maxCount")
	}

	contextID, ok := args["contextID"]
	if !ok || nbTles == "" {
		return fmt.Errorf("missing required argument: maxCount")
	}

	// Convert nbTles to an integer (assuming it's a string or a similar type in args)
	maxCount, err := strconv.Atoi(nbTles)
	if err != nil {
		return fmt.Errorf("invalid value for max: %v", err)
	}

	tles, err := h.tleService.FetchTLEFromSatCatByCategory(ctx, category, contextID)
	if err != nil {
		return fmt.Errorf("failed to fetch TLE catalog for category %s: %v", category, err)
	}

	// Retain only the maximum nbTles elements
	if len(tles) > maxCount {
		tles = tles[:maxCount] // Slice to keep only the first maxCount elements
	}

	log.Printf("Returning %d TLEs for category %s", len(tles), category)

	for _, rawTLE := range tles {
		err := h.upsertTLE(ctx, rawTLE, category)
		if err != nil {
			return fmt.Errorf("failed to upsert TLE for NORAD ID %s: %v", rawTLE.NoradID, err)
		}
	}

	return nil
}

func (h *CelestrackTleUploadHandler) upsertTLE(ctx context.Context, tle domain.TLE, category string) error {
	err := h.ensureSatelliteExists(ctx, tle.NoradID, category)
	if err != nil {
		return fmt.Errorf("failed to ensure satellite existence: %v", err)
	}
	return h.tleRepo.UpdateTle(ctx, tle)
}

func (h *CelestrackTleUploadHandler) ensureSatelliteExists(ctx context.Context, noradID string, category string) error {
	satellite, err := h.satelliteRepo.FindByNoradID(ctx, noradID)
	if err == nil && satellite.NoradID == noradID {
		return nil
	}

	nowUtc := time.Now().UTC()

	newSatellite, err := domain.NewSatellite(
		fmt.Sprintf("Unknown Satellite %s", noradID),
		noradID,
		domain.SatelliteType(strings.ToUpper(category)),
		true,  //active by default
		false, // not in favourite
		nowUtc,
	)

	if err != nil {
		return err
	}

	err = h.satelliteRepo.Save(ctx, newSatellite)
	if err != nil {
		return fmt.Errorf("failed to create satellite: %v", err)
	}

	return nil
}

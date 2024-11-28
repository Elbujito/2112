package handlers

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/internal/domain"
)

type TleServiceClient interface {
	FetchTLEFromSatCatByCategory(ctx context.Context, category string) ([]domain.TLE, error)
}

type TLEProvisionHandler struct {
	satelliteRepo domain.SatelliteRepository
	tleRepo       domain.TLERepository
	tleService    TleServiceClient
}

func NewTLEProvisionHandler(
	satelliteRepo domain.SatelliteRepository,
	tleRepo domain.TLERepository,
	tleService TleServiceClient) TLEProvisionHandler {
	return TLEProvisionHandler{
		satelliteRepo: satelliteRepo,
		tleRepo:       tleRepo,
		tleService:    tleService,
	}
}

func (h *TLEProvisionHandler) GetTask() Task {
	return Task{
		Name:         "fetchAndUpsertTLE",
		Description:  "Fetch TLE from CelesTrak and upsert it in the database",
		RequiredArgs: []string{"category"},
	}
}

func (h *TLEProvisionHandler) Run(ctx context.Context, args map[string]string) error {

	category, ok := args["category"]
	if !ok || category == "" {
		return fmt.Errorf("missing required argument: category")
	}

	tles, err := h.tleService.FetchTLEFromSatCatByCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to fetch TLE catalog for category %s: %v", category, err)
	}

	for _, rawTLE := range tles {
		err := h.upsertTLE(ctx, rawTLE, category)
		if err != nil {
			return fmt.Errorf("failed to upsert TLE for NORAD ID %s: %v", rawTLE.NoradID, err)
		}
	}

	return nil
}

func (h *TLEProvisionHandler) upsertTLE(ctx context.Context, tle domain.TLE, category string) error {
	err := h.ensureSatelliteExists(ctx, tle.NoradID, category)
	if err != nil {
		return fmt.Errorf("failed to ensure satellite existence: %v", err)
	}
	return h.tleRepo.Upsert(tle)
}

func (h *TLEProvisionHandler) ensureSatelliteExists(ctx context.Context, noradID string, category string) error {
	satellite, err := h.satelliteRepo.FindByNoradID(ctx, noradID)
	if err == nil && satellite.NoradID == noradID {
		return nil
	}

	newSatellite, err := domain.NewSatellite(
		fmt.Sprintf("Unknown Satellite %s", noradID),
		noradID,
		domain.SatelliteType(category),
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

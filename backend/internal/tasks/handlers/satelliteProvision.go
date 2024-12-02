package handlers

import (
	"context"

	"github.com/Elbujito/2112/internal/domain"
)

type SatelliteServiceClient interface {
	FetchAndStoreAllSatellites(ctx context.Context) ([]domain.Satellite, error)
}

type SatelliteProvisionHandler struct {
	satelliteRepo    domain.SatelliteRepository
	satelliteService SatelliteServiceClient
}

func NewSatelliteProvisionHandler(
	satelliteRepo domain.SatelliteRepository,
	satelliteService SatelliteServiceClient) SatelliteProvisionHandler {
	return SatelliteProvisionHandler{
		satelliteRepo:    satelliteRepo,
		satelliteService: satelliteService,
	}
}

func (h *SatelliteProvisionHandler) GetTask() Task {
	return Task{
		Name:         "fetchAndUpsertSatellite",
		Description:  "Fetch Satellite Metadata from CelesTrak and upsert it in the database",
		RequiredArgs: []string{""},
	}
}

func (h *SatelliteProvisionHandler) Run(ctx context.Context, args map[string]string) error {

	_, err := h.satelliteService.FetchAndStoreAllSatellites(ctx)
	if err != nil {
		return err
	}

	return nil
}

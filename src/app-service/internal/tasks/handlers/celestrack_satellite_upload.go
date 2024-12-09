package handlers

import (
	"context"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

type SatelliteServiceClient interface {
	FetchAndStoreAllSatellites(ctx context.Context) ([]domain.Satellite, error)
}

type CelesTrackSatelliteUploadHandler struct {
	satelliteRepo    domain.SatelliteRepository
	satelliteService SatelliteServiceClient
}

func NewCelesTrackSatelliteUploadHandler(
	satelliteRepo domain.SatelliteRepository,
	satelliteService SatelliteServiceClient) CelesTrackSatelliteUploadHandler {
	return CelesTrackSatelliteUploadHandler{
		satelliteRepo:    satelliteRepo,
		satelliteService: satelliteService,
	}
}

func (h *CelesTrackSatelliteUploadHandler) GetTask() Task {
	return Task{
		Name:         "celestrack_satellite_upload",
		Description:  "Fetch Satellite Metadata from CelesTrak and upsert it in the database",
		RequiredArgs: []string{""},
	}
}

func (h *CelesTrackSatelliteUploadHandler) Run(ctx context.Context, args map[string]string) error {

	_, err := h.satelliteService.FetchAndStoreAllSatellites(ctx)
	if err != nil {
		return err
	}

	return nil
}

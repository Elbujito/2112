package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/clients/redis"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
)

type SatellitesTilesMappingsHandler struct {
	tileRepo      domain.TileRepository
	tleRepo       repository.TleRepository
	satelliteRepo domain.SatelliteRepository
	mappingRepo   domain.MappingRepository
	redisClient   *redis.RedisClient
}

// NewSatellitesTilesMappingsHandler creates a new instance of the handler.
func NewSatellitesTilesMappingsHandler(
	tileRepo domain.TileRepository,
	tleRepo repository.TleRepository,
	satelliteRepo domain.SatelliteRepository,
	mappingRepo domain.MappingRepository,
	redisClient *redis.RedisClient,
) SatellitesTilesMappingsHandler {
	return SatellitesTilesMappingsHandler{
		tileRepo:      tileRepo,
		tleRepo:       tleRepo,
		satelliteRepo: satelliteRepo,
		mappingRepo:   mappingRepo,
		redisClient:   redisClient,
	}
}

func (h *SatellitesTilesMappingsHandler) GetTask() Task {
	return Task{
		Name:         "satellites_tiles_mappings",
		Description:  "Computes satellite visibilities for all tiles by satellite horizon",
		RequiredArgs: []string{"visibilityRadiusKm"},
	}
}

// Run executes the visibility computation process.
func (h *SatellitesTilesMappingsHandler) Run(ctx context.Context, args map[string]string) error {
	radiusInKm, err := strconv.ParseFloat(args["visibilityRadiusKm"], 64)
	if err != nil {
		return fmt.Errorf("invalid radius: %w", err)
	}

	// Call the Subscribe method
	h.Subscribe(ctx, "event_satellite_positions_updated", radiusInKm)

	return nil
}

// Run executes the visibility computation process.
func (h *SatellitesTilesMappingsHandler) Exec(ctx context.Context, id string, startTime time.Time, endTime time.Time, radiusInKm float64) error {

	sat, err := h.satelliteRepo.FindByNoradID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to fetch satellites: %w", err)
	}

	positions, err := h.tleRepo.QuerySatellitePositions(ctx, sat.NoradID, startTime, endTime)
	if err != nil {
		log.Printf("Error query satellite positions for satellite %s: %v\n", sat.NoradID, err)
		return nil
	}
	for _, pos := range positions {
		err = h.computeTileMappings(ctx, sat, pos, radiusInKm)
		if err != nil {
			log.Printf("Error computing mappings for satellite %s: %v\n", sat.NoradID, err)
		}
	}
	return nil
}

// computeTileMappings computes visibility for a satellite position.
func (h *SatellitesTilesMappingsHandler) computeTileMappings(
	ctx context.Context,
	sat domain.Satellite,
	position domain.SatellitePosition,
	radiusInKm float64,
) error {
	visibleTiles, err := h.tileRepo.FindTilesVisibleFromPoint(ctx, position.Latitude, position.Longitude, radiusInKm)
	if err != nil {
		return fmt.Errorf("failed to find visible tiles: %w", err)
	}

	mappings := make([]domain.TileSatelliteMapping, len(visibleTiles))
	for i, tile := range visibleTiles {
		mappings[i] = domain.NewMapping(
			sat.NoradID,
			tile.ID,
			position.Time,
			position.Altitude,
		)
	}

	if len(mappings) > 0 {
		if err := h.mappingRepo.SaveBatch(ctx, mappings); err != nil {
			return fmt.Errorf("failed to save mappings: %w", err)
		}
	}

	return nil
}

// Subscribe listens for satellite position updates and computes visibility.
func (h *SatellitesTilesMappingsHandler) Subscribe(ctx context.Context, channel string, radiusInKm float64) error {
	err := h.redisClient.Subscribe(ctx, channel, func(message string) error {
		// Parse the incoming message
		var update struct {
			SatelliteID string `json:"satellite_id"`
			StartTime   string `json:"start_time"`
			EndTime     string `json:"end_time"`
		}
		if err := json.Unmarshal([]byte(message), &update); err != nil {
			log.Printf("Failed to parse update message: %v\n", err)
			return err
		}

		// Convert StartTime and EndTime to time.Time
		startTime, err := time.Parse(time.RFC3339, update.StartTime)
		if err != nil {
			log.Printf("Failed to parse start time: %v\n", err)
			return err
		}
		endTime, err := time.Parse(time.RFC3339, update.EndTime)
		if err != nil {
			log.Printf("Failed to parse end time: %v\n", err)
			return err
		}

		return h.Exec(ctx, update.SatelliteID, startTime, endTime, radiusInKm)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel %s: %w", channel, err)
	}

	log.Printf("Subscribed to Redis channel: %s\n", channel)
	return nil
}

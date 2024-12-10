package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

// TleServiceClient defines the interface for fetching TLE data.
type FetchTleServiceClient interface {
	FetchTLEFromSatCatByCategory(ctx context.Context, category string) ([]domain.TLE, error)
}

// RedisPublisher defines the interface for publishing messages to Redis.
type RedisPublisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}

// TleFetchAndPublishHandler handles fetching TLE data and publishing it to Redis.
type TleFetchAndPublishHandler struct {
	tleService     FetchTleServiceClient
	redisPublisher RedisPublisher
}

// NewTleFetchAndPublishHandler creates a new instance of TleFetchAndPublishHandler.
func NewTleFetchAndPublishHandler(tleService FetchTleServiceClient, redisPublisher RedisPublisher) TleFetchAndPublishHandler {
	return TleFetchAndPublishHandler{
		tleService:     tleService,
		redisPublisher: redisPublisher,
	}
}

// GetTask returns the task metadata.
func (h *TleFetchAndPublishHandler) GetTask() Task {
	return Task{
		Name:         "tle_fetch_and_publish",
		Description:  "Fetch TLE data and publish each entry to Redis",
		RequiredArgs: []string{"category", "redis_channel"},
	}
}

// Run executes the task.
func (h *TleFetchAndPublishHandler) Run(ctx context.Context, args map[string]string) error {
	category, ok := args["category"]
	if !ok || category == "" {
		return fmt.Errorf("missing required argument: category")
	}

	redisChannel, ok := args["redis_channel"]
	if !ok || redisChannel == "" {
		return fmt.Errorf("missing required argument: redis_channel")
	}

	tles, err := h.tleService.FetchTLEFromSatCatByCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to fetch TLE data for category %s: %v", category, err)
	}

	for _, tle := range tles {
		message, err := json.Marshal(tle)
		if err != nil {
			return fmt.Errorf("failed to marshal TLE for NORAD ID %s: %v", tle.NoradID, err)
		}

		err = h.redisPublisher.Publish(ctx, redisChannel, message)
		if err != nil {
			return fmt.Errorf("failed to publish TLE for NORAD ID %s: %v", tle.NoradID, err)
		}
	}

	return nil
}

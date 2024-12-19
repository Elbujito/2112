package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/clients/redis"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

// ComputeVisibilitiessHandler handles visibility computation for satellites based on user locations.
type ComputeVisibilitiessHandler struct {
	tileRepo    domain.TileRepository
	mappingRepo domain.MappingRepository
	redisClient *redis.RedisClient
}

// NewComputeVisibilitiessHandler initializes a new handler instance.
func NewComputeVisibilitiessHandler(
	tileRepo domain.TileRepository,
	mappingRepo domain.MappingRepository,
	redisClient *redis.RedisClient,
) ComputeVisibilitiessHandler {
	return ComputeVisibilitiessHandler{
		tileRepo:    tileRepo,
		mappingRepo: mappingRepo,
		redisClient: redisClient,
	}
}

func (h *ComputeVisibilitiessHandler) GetTask() Task {
	return Task{
		Name:         "compute_visibilities",
		Description:  "Computes satellite visibilities for all tiles by satellite path",
		RequiredArgs: []string{},
	}
}

// Run executes the visibility computation process.
func (h *ComputeVisibilitiessHandler) Run(ctx context.Context, args map[string]string) error {
	log.Println("Starting Run method")
	log.Println("Subscribing to visibility_requests channel")
	return h.Subscribe(ctx, "visibility_requests")
}

// Subscribe listens for user location updates and computes visibilities.
func (h *ComputeVisibilitiessHandler) Subscribe(ctx context.Context, channel string) error {
	log.Printf("Subscribing to Redis channel: %s\n", channel)

	err := h.redisClient.Subscribe(ctx, channel, func(message string) error {
		// Parse the incoming message
		var request struct {
			UID       string  `json:"uid"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Radius    float64 `json:"radius"`
			Horizon   float64 `json:"horizon"`
			StartTime string  `json:"startTime"`
			EndTime   string  `json:"endTime"`
		}

		if err := json.Unmarshal([]byte(message), &request); err != nil {
			return fmt.Errorf("failed to parse update message: %w", err)
		}

		// Convert StartTime and EndTime to time.Time
		startTime, err := time.Parse(time.RFC3339, request.StartTime)
		if err != nil {
			return fmt.Errorf("failed to parse start time: %w", err)
		}
		endTime, err := time.Parse(time.RFC3339, request.EndTime)
		if err != nil {
			return fmt.Errorf("failed to parse end time: %w", err)
		}

		log.Printf("Received visibility request for UID: %s at location (%.6f, %.6f) with radius %.2f, horizon %.2f, from %s to %s\n",
			request.UID, request.Latitude, request.Longitude, request.Radius, request.Horizon, startTime, endTime)

		// Compute visibility for the received request
		return h.computeVisibility(ctx, request.UID, request.Latitude, request.Longitude, request.Radius, startTime, endTime)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to channel %s: %w", channel, err)
	}

	log.Printf("Successfully subscribed to Redis channel: %s\n", channel)
	return nil
}

// computeVisibility computes visibilities for a given user location and time range.
func (h *ComputeVisibilitiessHandler) computeVisibility(ctx context.Context, uid string, latitude, longitude, radius float64, startTime, endTime time.Time) error {
	// Validate location bounds
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude out of bounds: %f", latitude)
	}
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude out of bounds: %f", longitude)
	}
	if radius <= 0 {
		return fmt.Errorf("radius must be greater than 0")
	}

	// Step 1: Find tiles intersecting the user location
	tiles, err := h.tileRepo.FindTilesIntersectingLocation(ctx, latitude, longitude, radius)
	if err != nil {
		return fmt.Errorf("failed to find tiles intersecting location: %w", err)
	}
	if len(tiles) == 0 {
		log.Printf("No tiles found for location: (%f, %f) with radius: %f\n", latitude, longitude, radius)
		return nil
	}

	// Step 2: Extract tile IDs
	var tileIDs []string
	for _, tile := range tiles {
		tileIDs = append(tileIDs, tile.ID)
	}

	// Step 3: Find satellites associated with the identified tiles
	satellites, err := h.mappingRepo.FindSatellitesForTiles(ctx, tileIDs)
	if err != nil {
		return fmt.Errorf("failed to find satellites for tiles: %w", err)
	}
	if len(satellites) == 0 {
		log.Printf("No satellites found for the identified tiles.\n")
		return nil
	}

	// Prepare Redis key for cached visibilities
	key := fmt.Sprintf("satellite_visibilities:%s", uid)

	// Step 4: Publish visibility results
	var visibilities []map[string]interface{}
	for _, satellite := range satellites {
		visibility := map[string]interface{}{
			"satelliteId":   satellite.NoradID,
			"satelliteName": satellite.Name,
			"aos":           startTime.Format(time.RFC3339), // Acquisition of Signal time
			"los":           endTime.Format(time.RFC3339),   // Loss of Signal time
			"userLocation": map[string]interface{}{
				"latitude":  latitude,
				"longitude": longitude,
				"radius":    radius,
				"horizon":   30, // Example horizon value
			},
			"uid": uid,
		}

		visibilities = append(visibilities, visibility)

		// Publish the visibility message to the broker
		if err := h.redisClient.Publish(ctx, key, visibility); err != nil {
			log.Printf("Failed to publish visibility for satellite %s: %v\n", satellite.NoradID, err)
		} else {
			log.Printf("Published visibility for satellite %s (UID: %s)\n", satellite.Name, uid)
		}
	}

	// Cache all visibilities for the user in Redis
	cachedData, err := json.Marshal(visibilities)
	if err != nil {
		return fmt.Errorf("failed to serialize visibilities: %w", err)
	}

	if err := h.redisClient.Set(ctx, key, cachedData); err != nil {
		return fmt.Errorf("failed to cache visibilities: %w", err)
	}

	log.Printf("Cached visibilities for UID: %s\n", uid)
	return nil
}

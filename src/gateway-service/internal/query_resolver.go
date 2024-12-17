package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
	"github.com/go-redis/redis/v8"
)

type queryResolver struct {
	*Resolver
}

func (q *queryResolver) SatellitePosition(ctx context.Context, id string) (*model.SatellitePosition, error) {
	data, err := rdb.Get(ctx, "satellite_positions:"+id).Result()
	if err != nil {
		log.Printf("Error fetching SatellitePosition from Redis: %v", err)
		return nil, err
	}

	var position model.SatellitePosition
	if err := json.Unmarshal([]byte(data), &position); err != nil {
		log.Printf("Error unmarshalling SatellitePosition: %v", err)
		return nil, err
	}

	return &position, nil
}

func (q *queryResolver) SatelliteTle(ctx context.Context, id string) (*model.SatelliteTle, error) {
	data, err := rdb.Get(ctx, "satellite_tle:"+id).Result()
	if err != nil {
		log.Printf("Error fetching SatelliteTle from Redis: %v", err)
		return nil, err
	}

	var tle model.SatelliteTle
	if err := json.Unmarshal([]byte(data), &tle); err != nil {
		log.Printf("Error unmarshalling SatelliteTle: %v", err)
		return nil, err
	}

	return &tle, nil
}

func (q *queryResolver) SatellitePositionsInRange(ctx context.Context, id string, startTime string, endTime string) ([]*model.SatellitePosition, error) {
	log.Printf("SatellitePositionsInRange query called with id: %s, startTime: %s, endTime: %s", id, startTime, endTime)

	// Parse startTime and endTime
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		log.Printf("Error parsing startTime: %v", err)
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		log.Printf("Error parsing endTime: %v", err)
		return nil, err
	}

	// Prepare Redis key and query sorted set
	// Define the Redis key pattern for satellite positions
	key := fmt.Sprintf("satellite_positions:%s", id)

	// Convert time to UNIX timestamps for range query
	startTimestamp := strconv.FormatInt(start.Unix(), 10)
	endTimestamp := strconv.FormatInt(end.Unix(), 10)

	log.Printf("Fetching positions from sorted set: %s", key)

	// Query Redis for members within the specified score range
	results, err := rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: startTimestamp,
		Max: endTimestamp,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to query Redis for key %s with score range [%s, %s]: %w", key, startTimestamp, endTimestamp, err)
	}

	// Parse results into SatellitePosition objects
	positions := []*model.SatellitePosition{}
	for _, raw := range results {
		var position model.SatellitePosition
		if err := json.Unmarshal([]byte(raw), &position); err != nil {
			log.Printf("Error unmarshalling SatellitePosition: %v", err)
			continue
		}
		positions = append(positions, &position)
	}

	log.Printf("SatellitePositionsInRange query completed. Found %d positions within range.", len(positions))
	return positions, nil
}

func (q *queryResolver) SatelliteVisibilities(ctx context.Context, latitude float64, longitude float64) ([]*model.TileVisibility, error) {
	log.Printf("SatelliteVisibilities query called with latitude: %f, longitude: %f", latitude, longitude)

	// Define the Redis key for visibilities
	key := "satellite_visibilities"

	// Fetch all visibility data from Redis (or filter by geo hash/quadkey if applicable)
	results, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("Error fetching visibility data from Redis: %v", err)
		return nil, fmt.Errorf("failed to fetch visibility data: %w", err)
	}

	// Parse visibility data
	visibilities := []*model.TileVisibility{}
	for _, raw := range results {
		var visibility model.TileVisibility
		if err := json.Unmarshal([]byte(raw), &visibility); err != nil {
			log.Printf("Error unmarshalling TileVisibility: %v", err)
			continue
		}

		// Filter by latitude and longitude if necessary
		// For now, assuming all visibility data is returned
		visibilities = append(visibilities, &visibility)
	}

	log.Printf("SatelliteVisibilities query completed. Found %d visibilities.", len(visibilities))
	return visibilities, nil
}

func (q *queryResolver) RequestSatelliteVisibilitiesInRange(ctx context.Context, latitude float64, longitude float64, startTime string, endTime string) (bool, error) {
	visibilityRequestChannel := "visibility_requests" // Redis channel for publishing visibility requests
	log.Printf("RequestSatelliteVisibilitiesInRange request received: latitude=%f, longitude=%f, startTime=%s, endTime=%s", latitude, longitude, startTime, endTime)

	// Generate a unique channel ID for the response based on request parameters
	channelID := generateVisibilityChannelID(latitude, longitude, startTime, endTime)

	// Prepare the request payload
	requestPayload := map[string]string{
		"latitude":  fmt.Sprintf("%.6f", latitude),
		"longitude": fmt.Sprintf("%.6f", longitude),
		"startTime": startTime,
		"endTime":   endTime,
		"channel":   channelID, // Clients will listen to this channel
	}

	// Serialize the payload into JSON
	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshalling request payload: %v", err)
		return false, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Publish the request to Redis
	err = rdb.Publish(ctx, visibilityRequestChannel, requestJSON).Err()
	if err != nil {
		log.Printf("Error publishing visibility request to Redis channel '%s': %v", visibilityRequestChannel, err)
		return false, fmt.Errorf("failed to publish request to channel '%s': %w", visibilityRequestChannel, err)
	}

	log.Printf("Successfully published visibility request to Redis channel '%s'. Listening channel: %s", visibilityRequestChannel, channelID)
	return true, nil
}

// Helper function to generate a deterministic Redis channel ID for visibility results
func generateVisibilityChannelID(latitude float64, longitude float64, startTime string, endTime string) string {
	return fmt.Sprintf("visibilities:lat%.6f:lon%.6f:start%s:end%s", latitude, longitude, startTime, endTime)
}

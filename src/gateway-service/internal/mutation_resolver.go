package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
)

type mutationResolver struct {
	*CustomResolver
}

func (r *mutationResolver) RequestSatelliteVisibilities(
	ctx context.Context,
	uid string,
	userLocation model.UserLocationInput,
	startTime string,
	endTime string,
) (bool, error) {
	visibilityRequestChannel := "visibility_requests" // Redis channel for visibility requests

	log.Printf(
		"RequestSatelliteVisibilities called: uid=%s, latitude=%.6f, longitude=%.6f, radius=%.2f, horizon=%.2f, startTime=%s, endTime=%s",
		uid, userLocation.Latitude, userLocation.Longitude, userLocation.Radius, userLocation.Horizon, startTime, endTime,
	)

	// Prepare the request payload
	requestPayload := map[string]interface{}{
		"uid":       uid,
		"latitude":  userLocation.Latitude,
		"longitude": userLocation.Longitude,
		"radius":    userLocation.Radius,
		"horizon":   userLocation.Horizon,
		"startTime": startTime,
		"endTime":   endTime,
	}

	// Serialize the payload into JSON
	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshalling request payload: %v", err)
		return false, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Publish the request to Redis
	err = r.rdb.Publish(ctx, visibilityRequestChannel, requestJSON).Err()
	if err != nil {
		log.Printf("Error publishing visibility request to Redis channel '%s': %v", visibilityRequestChannel, err)
		return false, fmt.Errorf("failed to publish request to channel '%s': %w", visibilityRequestChannel, err)
	}

	log.Printf("Successfully published visibility request to Redis channel '%s': %s", visibilityRequestChannel, string(requestJSON))
	return true, nil
}

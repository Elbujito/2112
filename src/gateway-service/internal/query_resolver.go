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
	*CustomResolver
}

func (q *queryResolver) SatellitePosition(ctx context.Context, id string) (*model.SatellitePosition, error) {
	// Fetch satellite position data from Redis
	data, err := q.CustomResolver.rdb.Get(ctx, "satellite_positions:"+id).Result()
	if err != nil {
		log.Printf("Error fetching SatellitePosition from Redis: %v", err)
		return nil, fmt.Errorf("satellite position not found for ID %s: %w", id, err)
	}

	// Unmarshal JSON data into SatellitePosition struct
	var position model.SatellitePosition
	if err := json.Unmarshal([]byte(data), &position); err != nil {
		log.Printf("Error unmarshalling SatellitePosition: %v", err)
		return nil, fmt.Errorf("failed to parse satellite position data: %w", err)
	}

	return &position, nil
}

func (q *queryResolver) SatelliteTle(ctx context.Context, id string) (*model.SatelliteTle, error) {
	// Fetch TLE data from Redis
	data, err := q.CustomResolver.rdb.Get(ctx, "satellite_tle:"+id).Result()
	if err != nil {
		log.Printf("Error fetching SatelliteTle from Redis: %v", err)
		return nil, fmt.Errorf("TLE data not found for satellite ID %s: %w", id, err)
	}

	// Unmarshal JSON data into SatelliteTle struct
	var tle model.SatelliteTle
	if err := json.Unmarshal([]byte(data), &tle); err != nil {
		log.Printf("Error unmarshalling SatelliteTle: %v", err)
		return nil, fmt.Errorf("failed to parse TLE data: %w", err)
	}

	return &tle, nil
}

func (q *queryResolver) SatellitePositionsInRange(ctx context.Context, id string, startTime string, endTime string) ([]*model.SatellitePosition, error) {
	// Parse startTime and endTime into time.Time objects
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse startTime: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endTime: %w", err)
	}

	// Prepare Redis query keys and timestamps
	key := fmt.Sprintf("satellite_positions:%s", id)
	startTimestamp := strconv.FormatInt(start.Unix(), 10)
	endTimestamp := strconv.FormatInt(end.Unix(), 10)

	// Query Redis for position data within the specified range
	results, err := q.CustomResolver.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: startTimestamp,
		Max: endTimestamp,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to query Redis for satellite positions: %w", err)
	}

	// Parse results into SatellitePosition objects
	var positions []*model.SatellitePosition
	for _, raw := range results {
		var position model.SatellitePosition
		if err := json.Unmarshal([]byte(raw), &position); err != nil {
			log.Printf("Error unmarshalling SatellitePosition: %v", err)
			continue
		}
		positions = append(positions, &position)
	}

	return positions, nil
}

func (q *queryResolver) CachedSatelliteVisibilities(ctx context.Context, uid string, userLocation model.UserLocationInput, startTime string, endTime string) ([]*model.SatelliteVisibility, error) {
	// Prepare Redis key for cached visibilities
	key := fmt.Sprintf("satellite_visibilities:%s", uid)

	// Fetch cached visibility data from Redis
	data, err := q.CustomResolver.rdb.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Error fetching cached satellite visibilities for UID %s: %v", uid, err)
		return nil, fmt.Errorf("cached visibilities not found for UID %s: %w", uid, err)
	}

	// Parse cached visibility data into SatelliteVisibility objects
	var visibilities []*model.SatelliteVisibility
	if err := json.Unmarshal([]byte(data), &visibilities); err != nil {
		log.Printf("Error unmarshalling CachedSatelliteVisibilities: %v", err)
		return nil, fmt.Errorf("failed to parse cached satellite visibilities: %w", err)
	}

	return visibilities, nil
}

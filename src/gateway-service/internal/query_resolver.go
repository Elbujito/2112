package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Elbujito/2112/graphql-api/graph/model"
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

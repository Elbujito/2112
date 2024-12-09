package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Elbujito/2112/graphql-api/graph/model"
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

	positions := []*model.SatellitePosition{}

	keys, err := rdb.Keys(ctx, "satellite_positions:"+id+":*").Result()
	if err != nil {
		log.Printf("Error fetching keys for SatellitePosition in Redis: %v", err)
		return nil, err
	}

	for _, key := range keys {
		data, err := rdb.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Error fetching SatellitePosition from Redis: %v", err)
			continue
		}

		var position model.SatellitePosition
		if err := json.Unmarshal([]byte(data), &position); err != nil {
			log.Printf("Error unmarshalling SatellitePosition: %v", err)
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, position.Timestamp)
		if err != nil {
			log.Printf("Error parsing timestamp: %v", err)
			continue
		}

		if timestamp.After(start) && timestamp.Before(end) {
			positions = append(positions, &position)
		}
	}

	return positions, nil
}

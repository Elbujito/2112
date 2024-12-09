package main

import (
	"context"
	"encoding/json"
	"log"

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

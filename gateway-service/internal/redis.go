package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Elbujito/2112/graphql-api/graph/model"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis(host, port string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func subscribeToRedisForTleUpdates(resolver *Resolver) {
	pubsub := rdb.Subscribe(ctx, "satellite_tle")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var tle model.SatelliteTle
		if err := json.Unmarshal([]byte(msg.Payload), &tle); err != nil {
			log.Printf("Error unmarshalling SatelliteTle: %v", err)
			continue
		}

		// Simulate position update with timestamp
		position := &model.SatellitePosition{
			ID:        tle.ID,
			Name:      tle.Name,
			Latitude:  0.0, // Placeholder value
			Longitude: 0.0, // Placeholder value
			Altitude:  0.0, // Placeholder value
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		// Publish to Redis and notify subscribers
		publishSatellitePosition(position)
		resolver.NotifySubscribers(position)
	}
}

func publishSatellitePosition(position *model.SatellitePosition) {
	payload, err := json.Marshal(position)
	if err != nil {
		log.Printf("Error marshalling SatellitePosition: %v", err)
		return
	}

	channel := "satellite_positions:" + position.ID
	if err := rdb.Publish(ctx, channel, payload).Err(); err != nil {
		log.Printf("Error publishing SatellitePosition to Redis: %v", err)
	}
}

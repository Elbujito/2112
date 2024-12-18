package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
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

func subscribeToRedisForPositionUpdates(resolver *Resolver) {
	pubsub := rdb.Subscribe(ctx, "satellite_positions")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var position model.SatellitePosition
		if err := json.Unmarshal([]byte(msg.Payload), &position); err != nil {
			log.Printf("Error unmarshalling SatellitePosition: %v", err)
			continue
		}

		resolver.NotifySubscribers(&position)
	}
}

func subscribeToRedisForVisibilityUpdates(resolver *Resolver) {
	pubsub := rdb.Subscribe(ctx, "satellite_visibilities")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var position model.SatellitePosition
		if err := json.Unmarshal([]byte(msg.Payload), &position); err != nil {
			log.Printf("Error unmarshalling SatellitePosition: %v", err)
			continue
		}

		resolver.NotifySubscribers(&position)
	}
}

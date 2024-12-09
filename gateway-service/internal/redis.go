package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
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

func subscribeToTleAndGeneratePositions(resolver *Resolver) {
	pubsub := rdb.Subscribe(ctx, "satellite_tle")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var tle model.SatelliteTle
		if err := json.Unmarshal([]byte(msg.Payload), &tle); err != nil {
			log.Printf("Error unmarshalling SatelliteTle: %v", err)
			continue
		}

		// Simulate position generation
		position := generateSatellitePosition(&tle)

		// Publish position to Redis
		publishSatellitePosition(position)

		// Notify WebSocket subscribers
		resolver.Mutex.Lock()
		for id, ch := range resolver.PositionSubscribers {
			if id == position.ID {
				ch <- position
			}
		}
		resolver.Mutex.Unlock()
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

func generateSatellitePosition(tle *model.SatelliteTle) *model.SatellitePosition {
	rand.Seed(time.Now().UnixNano())
	return &model.SatellitePosition{
		ID:        tle.ID,
		Name:      tle.Name,
		Latitude:  rand.Float64()*180 - 90,  // Random latitude
		Longitude: rand.Float64()*360 - 180, // Random longitude
		Altitude:  rand.Float64()*500 + 200, // Random altitude
	}
}

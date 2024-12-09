package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	model "github.com/Elbujito/2112/graphql-api/graph/model"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

// In-memory storage
var satelliteData = make(map[string]model.SatellitePosition)
var satelliteTleData = make(map[string]model.SatelliteTle)
var messageHistory []string
var mutex = &sync.Mutex{} // Mutex to handle concurrent access to data

// Initialize Redis client
func initRedis(host, port string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	// Test Redis connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Start the Redis listener in a background goroutine
	go subscribeToRedis()
}

// Subscribe to Redis channels and update in-memory data
func subscribeToRedis() {
	pubsub := rdb.Subscribe(ctx, "satellite_tle")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var tle model.SatelliteTle
		if err := json.Unmarshal([]byte(msg.Payload), &tle); err != nil {
			log.Printf("Error unmarshalling SatelliteTle data: %v", err)
			continue
		}

		mutex.Lock()
		satelliteTleData[tle.ID] = tle
		messageHistory = append(messageHistory, msg.Payload)
		mutex.Unlock()

		log.Printf("Updated SatelliteTle data: %v", tle)
	}
}

// Publish satellite position data to Redis
func publishSatellitePosition(position model.SatellitePosition) {
	payload, err := json.Marshal(position)
	if err != nil {
		log.Printf("Error marshalling SatellitePosition data: %v", err)
		return
	}

	err = rdb.Publish(ctx, "satellite_positions", payload).Err()
	if err != nil {
		log.Printf("Error publishing SatellitePosition data to Redis: %v", err)
		return
	}

	log.Printf("Published SatellitePosition data to Redis: %v", position)
}

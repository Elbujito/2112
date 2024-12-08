package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

// Satellite data model with TLE (Two-Line Element)
type Satellite struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TLELine1 string `json:"tle_line1"`
	TLELine2 string `json:"tle_line2"`
}

var ctx = context.Background()

// Mocked satellite TLE data
var satellites = []Satellite{
	{
		ID:       "1",
		Name:     "Satellite 1",
		TLELine1: "1 25544U 98067A   21349.25024444  .00001234  00000-0  12345-4 0  9990",
		TLELine2: "2 25544  51.6404 102.7490 0007637  35.6700 324.1674 15.48906415594989",
	},
	{
		ID:       "2",
		Name:     "Satellite 2",
		TLELine1: "1 20580U 90037B   21349.25024444  .00001023  00000-0  80000-4 0  9875",
		TLELine2: "2 20580  42.4104 315.8404 0012005  22.7400 182.1140 15.54578220107342",
	},
	{
		ID:       "3",
		Name:     "Satellite 3",
		TLELine1: "1 18624U 05037A   21349.25024444  .00004567  00000-0  30252-3 0  8961",
		TLELine2: "2 18624  35.2542  22.4517 0000565  19.2850 340.8676 15.10439051098236",
	},
}

func publishSatelliteData(rdb *redis.Client) {
	// Loop through the mocked data and publish each satellite with TLE to the "satellite_tle_data" channel
	for _, satellite := range satellites {
		satelliteData, err := json.Marshal(satellite)
		if err != nil {
			log.Fatalf("Failed to marshal satellite data: %v", err)
		}
		err = rdb.Publish(ctx, "satellite_tle_data", satelliteData).Err()
		if err != nil {
			log.Fatalf("Failed to publish satellite data: %v", err)
		}
		fmt.Printf("Published satellite with TLE: %v\n", satellite)
	}
}

func main() {
	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-service:6379", // Redis service from Docker Compose
	})

	// Test Redis connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Publish mocked satellite data with TLE to Redis on startup
	publishSatelliteData(rdb)

	// Start the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "App Service is running")
	})

	fmt.Println("App Service started on port 8081")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

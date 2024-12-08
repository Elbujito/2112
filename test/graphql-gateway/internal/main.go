package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// Redis context and client
var ctx = context.Background()
var rdb *redis.Client

// Satellite data storage (in-memory storage for simplicity)
var satelliteData = make(map[string]SatellitePositionResolver)
var satelliteTleData = make(map[string]SatelliteTleResolver)
var messageHistory []string
var mutex = &sync.Mutex{} // Mutex to handle concurrent access to satelliteData

// Define your schema inline or load from a file
var schemaString = `
schema {
	query: Query
	subscription: Subscription
}

type Query {
	satellitePosition(id: ID!): SatellitePosition
	satelliteTle(id: ID!): SatelliteTle
	messageHistory: [String!]!
}

type SatellitePosition {
	id: ID!
	name: String!
	latitude: Float!
	longitude: Float!
	altitude: Float!
}

type SatelliteTle {
	id: ID!
	name: String!
	tleLine1: String!
	tleLine2: String!
}

type Subscription {
	messageReceived: String!
}
`

// Satellite model with TLE data
type Satellite struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	TLELine1  string  `json:"tle_line1"`
	TLELine2  string  `json:"tle_line2"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

// Define your resolvers
type Resolver struct{}

// SatellitePosition resolver
func (r *Resolver) SatellitePosition(ctx context.Context, args struct{ ID graphql.ID }) *SatellitePositionResolver {
	mutex.Lock()
	defer mutex.Unlock()

	// Return satellite position data if exists
	if satellite, exists := satelliteData[string(args.ID)]; exists {
		return &satellite
	}
	return nil // Return nil if not found
}

// SatelliteTle resolver
func (r *Resolver) SatelliteTle(ctx context.Context, args struct{ ID graphql.ID }) *SatelliteTleResolver {
	mutex.Lock()
	defer mutex.Unlock()

	// Return satellite TLE data if exists
	if satellite, exists := satelliteTleData[string(args.ID)]; exists {
		return &satellite
	}
	return nil // Return nil if not found
}

// messageHistory resolver
func (r *Resolver) MessageHistory(ctx context.Context) []string {
	mutex.Lock()
	defer mutex.Unlock()
	return messageHistory
}

// Subscription for incoming messages
func (r *Resolver) MessageReceived(ctx context.Context) <-chan string {
	ch := make(chan string)

	// Start the Redis subscription in a separate goroutine
	go subscribeToRedis(ch)

	return ch
}

// SatellitePositionResolver for the data model
type SatellitePositionResolver struct {
	IDField        graphql.ID
	NameField      string
	LatitudeField  float64
	LongitudeField float64
	AltitudeField  float64
}

// SatelliteTleResolver for the data model
type SatelliteTleResolver struct {
	IDField       graphql.ID
	NameField     string
	TleLine1Field string
	TleLine2Field string
}

// Implementing methods to resolve fields for SatellitePosition
func (r *SatellitePositionResolver) ID() graphql.ID {
	return r.IDField
}

func (r *SatellitePositionResolver) Name() string {
	return r.NameField
}

func (r *SatellitePositionResolver) Latitude() float64 {
	return r.LatitudeField
}

func (r *SatellitePositionResolver) Longitude() float64 {
	return r.LongitudeField
}

func (r *SatellitePositionResolver) Altitude() float64 {
	return r.AltitudeField
}

// Implementing methods to resolve fields for SatelliteTle
func (r *SatelliteTleResolver) ID() graphql.ID {
	return r.IDField
}

func (r *SatelliteTleResolver) Name() string {
	return r.NameField
}

func (r *SatelliteTleResolver) TleLine1() string {
	return r.TleLine1Field
}

func (r *SatelliteTleResolver) TleLine2() string {
	return r.TleLine2Field
}

func main() {
	// Get environment variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	appPort := os.Getenv("GATEWAY_PORT")

	if redisHost == "" {
		redisHost = "localhost" // default if not set
	}
	if redisPort == "" {
		redisPort = "6379" // default if not set
	}
	if appPort == "" {
		appPort = "4000" // default if not set
	}

	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort, // Use the environment variables
	})

	// Test Redis connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Parse the GraphQL schema with resolvers
	schema := graphql.MustParseSchema(schemaString, &Resolver{})

	// Start the Redis listener in a background goroutine
	go subscribeToRedis(nil)

	// Set up the HTTP server
	http.Handle("/query", &relay.Handler{Schema: schema})
	log.Printf("GraphQL Gateway is running on http://localhost:%s/query", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}

// Function to handle Redis subscription for both satellite_positions and propagated_positions
func subscribeToRedis(ch chan string) {
	pubsub := rdb.Subscribe(ctx, "satellite_positions", "propagated_positions") // Subscribe to both channels
	defer pubsub.Close()

	// Continuously listen for messages on the "satellite_positions" and "propagated_positions" channels
	for msg := range pubsub.Channel() {
		// Parse the satellite data and update the in-memory store
		var satellite Satellite
		if err := json.Unmarshal([]byte(msg.Payload), &satellite); err != nil {
			log.Printf("Error unmarshalling satellite data: %v", err)
			continue
		}

		// If the message is TLE data, store it in satelliteTleData
		if satellite.TLELine1 != "" && satellite.TLELine2 != "" {
			mutex.Lock()
			satelliteTleData[satellite.ID] = SatelliteTleResolver{
				IDField:       graphql.ID(satellite.ID),
				NameField:     satellite.Name,
				TleLine1Field: satellite.TLELine1,
				TleLine2Field: satellite.TLELine2,
			}
			mutex.Unlock()
		}

		// If the message is a satellite position, store it in satelliteData
		if satellite.Latitude != 0 && satellite.Longitude != 0 {
			mutex.Lock()
			satelliteData[satellite.ID] = SatellitePositionResolver{
				IDField:        graphql.ID(satellite.ID),
				NameField:      satellite.Name,
				LatitudeField:  satellite.Latitude,
				LongitudeField: satellite.Longitude,
				AltitudeField:  satellite.Altitude, // Assuming the altitude is passed with the position data
			}
			mutex.Unlock()
		}

		// Add message to history
		mutex.Lock()
		messageHistory = append(messageHistory, msg.Payload)
		mutex.Unlock()

		// Send the message to the WebSocket channel (if ch is not nil)
		if ch != nil {
			ch <- msg.Payload
		}

		// Log the update
		log.Printf("Updated satellite data: %v", satellite)
	}
}

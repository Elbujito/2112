package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// Redis context and client
var ctx = context.Background()
var rdb *redis.Client

// Define your schema inline or load from a file
var schemaString = `
schema {
	query: Query
}

type Query {
	satellitePosition(id: ID!): SatellitePosition
	message(channel: String!): String!
}

type SatellitePosition {
	id: ID!
	name: String!
	latitude: Float!
	longitude: Float!
}
`

// Define your resolvers
type Resolver struct{}

// SatellitePosition resolver
func (r *Resolver) SatellitePosition(ctx context.Context, args struct{ ID graphql.ID }) *SatellitePositionResolver {
	// Returning the SatellitePositionResolver with hardcoded data
	return &SatellitePositionResolver{
		IDField:        args.ID, // Use the ID type here
		NameField:      "Satellite " + string(args.ID),
		LatitudeField:  40.7128,
		LongitudeField: -74.0060,
	}
}

// Message resolver for Redis Pub/Sub
func (r *Resolver) Message(ctx context.Context, args struct{ Channel string }) (string, error) {
	// Subscribe to the specified Redis channel
	sub := rdb.Subscribe(ctx, args.Channel)
	defer sub.Close()

	// Wait for a message from the channel
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}

	return msg.Payload, nil
}

// SatellitePositionResolver for the data model
type SatellitePositionResolver struct {
	IDField        graphql.ID // Use graphql.ID instead of string
	NameField      string
	LatitudeField  float64
	LongitudeField float64
}

// Implementing methods to resolve fields for SatellitePosition
func (r *SatellitePositionResolver) ID() graphql.ID {
	return r.IDField // Return the ID type
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

	// Set up the HTTP server
	http.Handle("/query", &relay.Handler{Schema: schema})
	log.Printf("GraphQL Gateway is running on http://localhost:%s/query", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}

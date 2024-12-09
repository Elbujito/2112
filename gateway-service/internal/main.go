package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Elbujito/2112/graphql-api/graph"
)

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	appPort := os.Getenv("GATEWAY_PORT")

	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}
	if appPort == "" {
		appPort = "4000"
	}

	resolver := NewResolver()
	initRedis(redisHost, redisPort)

	// Start the TLE subscription and position generation
	go subscribeToTleAndGeneratePositions(resolver)

	// GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	log.Printf("🚀 Server is running at http://localhost:%s/", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Elbujito/2112/src/graphql-api/go/graph"
	"github.com/rs/cors" // Add CORS middleware package
)

func main() {
	// Get environment variables or set defaults
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	appPort := getEnv("GATEWAY_PORT", "4000")

	log.Printf("Using Redis host: %s, port: %s", redisHost, redisPort)
	log.Printf("Using Gateway port: %s", appPort)

	// Initialize Redis
	initRedis(redisHost, redisPort)

	// Initialize Resolver
	resolver := NewResolver()

	// Start Redis subscription in a goroutine
	go subscribeToRedisForPositionUpdates(resolver)

	// Create GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Setup HTTP handlers
	mux := http.NewServeMux()
	mux.Handle("/query", srv)
	mux.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	// Add CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust based on your frontend's origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	// Start HTTP server with graceful shutdown
	server := &http.Server{
		Addr:    ":" + appPort,
		Handler: corsMiddleware,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 GraphQL endpoint: http://localhost:%s/query", appPort)
		log.Printf("🎮 Playground available at: http://localhost:%s/", appPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for termination signal for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting.")
}

// Helper to get environment variables with a default fallback
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

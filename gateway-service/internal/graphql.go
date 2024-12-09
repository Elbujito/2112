package main

import (
	"context"
	"sync"

	model "github.com/Elbujito/2112/graphql-api/graph/model"
)

// Resolver struct for custom resolvers
type Resolver struct{}

// Mutex for safe concurrent access
var graphqlMutex = &sync.Mutex{}

// SatellitePosition resolver
func (r *Resolver) SatellitePosition(ctx context.Context, id string) (*model.SatellitePosition, error) {
	graphqlMutex.Lock()
	defer graphqlMutex.Unlock()

	position, exists := satelliteData[id]
	if !exists {
		return nil, nil // Return nil if not found
	}
	return &position, nil
}

// SatelliteTle resolver
func (r *Resolver) SatelliteTle(ctx context.Context, id string) (*model.SatelliteTle, error) {
	graphqlMutex.Lock()
	defer graphqlMutex.Unlock()

	tle, exists := satelliteTleData[id]
	if !exists {
		return nil, nil // Return nil if not found
	}
	return &tle, nil
}

// MessageHistory resolver
func (r *Resolver) MessageHistory(ctx context.Context) ([]string, error) {
	graphqlMutex.Lock()
	defer graphqlMutex.Unlock()
	return messageHistory, nil
}

// MessageReceived subscription resolver
func (r *Resolver) MessageReceived(ctx context.Context) (<-chan string, error) {
	ch := make(chan string)

	go func() {
		for _, msg := range messageHistory {
			ch <- msg
		}
		close(ch)
	}()

	return ch, nil
}

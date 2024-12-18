package main

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
)

type subscriptionResolver struct {
	*Resolver
}

// SatellitePositionUpdated resolves the subscription for real-time position updates
func (s *subscriptionResolver) SatellitePositionUpdated(ctx context.Context, id string) (<-chan *model.SatellitePosition, error) {
	s.Mutex.Lock()
	ch := make(chan *model.SatellitePosition, 1)
	s.PositionSubscribers[id] = ch
	s.Mutex.Unlock()

	// Clean up subscription when the context is canceled
	go func() {
		<-ctx.Done()
		s.Mutex.Lock()
		delete(s.PositionSubscribers, id)
		close(ch)
		s.Mutex.Unlock()
	}()

	return ch, nil
}

// SatelliteVisibilityUpdated resolves the subscription for real-time visibility updates
func (s *subscriptionResolver) SatelliteVisibilityUpdated(ctx context.Context, latitude float64, longitude float64, startTime string, endTime string) (<-chan []*model.TileVisibility, error) {
	s.Mutex.Lock()
	key := createVisibilityKey(latitude, longitude)
	ch := make(chan []*model.TileVisibility, 1)
	s.VisibilitySubscribers[key] = ch
	s.Mutex.Unlock()

	// Clean up subscription when the context is canceled
	go func() {
		<-ctx.Done()
		s.Mutex.Lock()
		delete(s.VisibilitySubscribers, key)
		close(ch)
		s.Mutex.Unlock()
	}()

	return ch, nil
}

// Helper to generate a unique key for visibility subscribers based on latitude and longitude
func createVisibilityKey(latitude float64, longitude float64) string {
	return fmt.Sprintf("%.6f:%.6f", latitude, longitude)
}

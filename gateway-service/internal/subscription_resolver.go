package main

import (
	"context"

	"github.com/Elbujito/2112/graphql-api/graph/model"
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

package main

import (
	"context"

	"github.com/Elbujito/2112/graphql-api/graph/model"
)

type subscriptionResolver struct {
	*Resolver
}

func (s *subscriptionResolver) SatellitePositionUpdated(ctx context.Context, id string) (<-chan *model.SatellitePosition, error) {
	s.Mutex.Lock()
	ch := make(chan *model.SatellitePosition, 1)
	s.PositionSubscribers[id] = ch
	s.Mutex.Unlock()

	go func() {
		<-ctx.Done()
		s.Mutex.Lock()
		delete(s.PositionSubscribers, id)
		close(ch)
		s.Mutex.Unlock()
	}()

	return ch, nil
}

package main

import (
	"context"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
)

type subscriptionResolver struct {
	*CustomResolver
}

// SatellitePositionUpdated resolves the subscription for real-time position updates
func (s *subscriptionResolver) SatellitePositionUpdated(ctx context.Context, uid string, id string) (<-chan *model.SatellitePosition, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Initialize the user's position subscription map if it doesn't exist
	if _, exists := s.PositionSubscribers[uid]; !exists {
		s.PositionSubscribers[uid] = make(map[string]chan *model.SatellitePosition)
	}

	// Create a new channel for the satellite position updates
	ch := make(chan *model.SatellitePosition, 1)
	s.PositionSubscribers[uid][id] = ch

	// Clean up subscription when the context is canceled
	go func() {
		<-ctx.Done()
		s.cleanupPositionSubscriber(uid, id, ch)
	}()

	return ch, nil
}

// SatelliteVisibilityUpdated resolves the subscription for real-time visibility updates
func (s *subscriptionResolver) SatelliteVisibilityUpdated(ctx context.Context, uid string, userLocation model.UserLocationInput, startTime string, endTime string) (<-chan []*model.SatelliteVisibility, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Create a new channel for the visibility updates
	ch := make(chan []*model.SatelliteVisibility, 1)
	s.VisibilitySubscribers[uid] = ch

	// Clean up subscription when the context is canceled
	go func() {
		<-ctx.Done()
		s.cleanupVisibilitySubscriber(uid, ch)
	}()

	return ch, nil
}

// cleanupPositionSubscriber removes a position subscription and closes the channel
func (s *subscriptionResolver) cleanupPositionSubscriber(uid, id string, ch chan *model.SatellitePosition) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if userSubs, ok := s.PositionSubscribers[uid]; ok {
		delete(userSubs, id)
		if len(userSubs) == 0 {
			delete(s.PositionSubscribers, uid)
		}
		close(ch)
	}
}

// cleanupVisibilitySubscriber removes a visibility subscription and closes the channel
func (s *subscriptionResolver) cleanupVisibilitySubscriber(uid string, ch chan []*model.SatelliteVisibility) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	delete(s.VisibilitySubscribers, uid)
	close(ch)
}

package main

import (
	"sync"

	"github.com/Elbujito/2112/src/graphql-api/go/graph"
	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
	"github.com/go-redis/redis/v8"
)

type CustomResolver struct {
	graph.Resolver
	// Redis client instance
	rdb *redis.Client
	// Map for position subscriptions, keyed by uid and satellite ID
	PositionSubscribers map[string]map[string]chan *model.SatellitePosition
	// Map for visibility subscriptions, keyed by uid
	VisibilitySubscribers map[string]chan []*model.SatelliteVisibility
	Mutex                 sync.Mutex
}

// NewCustomResolver initializes a new CustomResolver instance with a Redis client
func NewCustomResolver(redisClient *redis.Client) *CustomResolver {
	return &CustomResolver{
		rdb:                   redisClient,
		PositionSubscribers:   make(map[string]map[string]chan *model.SatellitePosition),
		VisibilitySubscribers: make(map[string]chan []*model.SatelliteVisibility),
	}
}

// NotifyPositionSubscribers sends a position update to the relevant subscribers
func (r *CustomResolver) NotifyPositionSubscribers(uid string, position *model.SatellitePosition) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if userSubs, ok := r.PositionSubscribers[uid]; ok {
		if ch, ok := userSubs[position.ID]; ok {
			ch <- position
		}
	}
}

// NotifyVisibilitySubscribers sends a visibility update to the relevant subscribers
func (r *CustomResolver) NotifyVisibilitySubscribers(uid string, visibilities []*model.SatelliteVisibility) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if ch, ok := r.VisibilitySubscribers[uid]; ok {
		ch <- visibilities
	}
}

// AddPositionSubscriber adds a subscriber for satellite position updates
func (r *CustomResolver) AddPositionSubscriber(uid string, satelliteID string, ch chan *model.SatellitePosition) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if _, ok := r.PositionSubscribers[uid]; !ok {
		r.PositionSubscribers[uid] = make(map[string]chan *model.SatellitePosition)
	}
	r.PositionSubscribers[uid][satelliteID] = ch
}

// RemovePositionSubscriber removes a subscriber for satellite position updates
func (r *CustomResolver) RemovePositionSubscriber(uid string, satelliteID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if userSubs, ok := r.PositionSubscribers[uid]; ok {
		delete(userSubs, satelliteID)
		if len(userSubs) == 0 {
			delete(r.PositionSubscribers, uid)
		}
	}
}

// AddVisibilitySubscriber adds a subscriber for satellite visibility updates
func (r *CustomResolver) AddVisibilitySubscriber(uid string, ch chan []*model.SatelliteVisibility) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.VisibilitySubscribers[uid] = ch
}

// RemoveVisibilitySubscriber removes a subscriber for satellite visibility updates
func (r *CustomResolver) RemoveVisibilitySubscriber(uid string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.VisibilitySubscribers, uid)
}

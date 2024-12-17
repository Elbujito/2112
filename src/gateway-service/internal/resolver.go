package main

import (
	"sync"

	"github.com/Elbujito/2112/src/graphql-api/go/graph"
	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
)

type Resolver struct {
	PositionSubscribers   map[string]chan *model.SatellitePosition
	VisibilitySubscribers map[string]chan []*model.TileVisibility
	Mutex                 sync.Mutex
}

func NewResolver() *Resolver {
	return &Resolver{
		PositionSubscribers: make(map[string]chan *model.SatellitePosition),
	}
}

func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{r}
}

func (r *Resolver) NotifySubscribers(position *model.SatellitePosition) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	for id, ch := range r.PositionSubscribers {
		if id == position.ID {
			ch <- position
		}
	}
}

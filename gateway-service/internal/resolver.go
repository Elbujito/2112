package main

import (
	"sync"

	"github.com/Elbujito/2112/graphql-api/graph"
	"github.com/Elbujito/2112/graphql-api/graph/model"
)

type Resolver struct {
	PositionSubscribers map[string]chan *model.SatellitePosition
	Mutex               sync.Mutex
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

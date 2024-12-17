package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
)

// SatellitePosition is the resolver for the satellitePosition field.
func (r *queryResolver) SatellitePosition(ctx context.Context, id string) (*model.SatellitePosition, error) {
	panic(fmt.Errorf("not implemented: SatellitePosition - satellitePosition"))
}

// SatelliteTle is the resolver for the satelliteTle field.
func (r *queryResolver) SatelliteTle(ctx context.Context, id string) (*model.SatelliteTle, error) {
	panic(fmt.Errorf("not implemented: SatelliteTle - satelliteTle"))
}

// SatellitePositionsInRange is the resolver for the satellitePositionsInRange field.
func (r *queryResolver) SatellitePositionsInRange(ctx context.Context, id string, startTime string, endTime string) ([]*model.SatellitePosition, error) {
	panic(fmt.Errorf("not implemented: SatellitePositionsInRange - satellitePositionsInRange"))
}

// RequestSatelliteVisibilitiesInRange is the resolver for the requestSatelliteVisibilitiesInRange field.
func (r *queryResolver) RequestSatelliteVisibilitiesInRange(ctx context.Context, latitude float64, longitude float64, startTime string, endTime string) (bool, error) {
	panic(fmt.Errorf("not implemented: RequestSatelliteVisibilitiesInRange - requestSatelliteVisibilitiesInRange"))
}

// SatellitePositionUpdated is the resolver for the satellitePositionUpdated field.
func (r *subscriptionResolver) SatellitePositionUpdated(ctx context.Context, id string) (<-chan *model.SatellitePosition, error) {
	panic(fmt.Errorf("not implemented: SatellitePositionUpdated - satellitePositionUpdated"))
}

// SatelliteVisibilityUpdated is the resolver for the satelliteVisibilityUpdated field.
func (r *subscriptionResolver) SatelliteVisibilityUpdated(ctx context.Context, latitude float64, longitude float64, startTime string, endTime string) (<-chan []*model.TileVisibility, error) {
	panic(fmt.Errorf("not implemented: SatelliteVisibilityUpdated - satelliteVisibilityUpdated"))
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

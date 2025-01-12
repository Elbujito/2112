package domain

import (
	"context"
	"time"
)

// Context represents a logical grouping for satellites, TLEs, and mappings in the domain layer.
type Context struct {
	ModelBase
	Name        string // Unique name of the context
	Description string // Optional description of the context
	ActivatedAt *time.Time
}

// ContextSatellite represents the relationship between Context and Satellite in the domain layer.
type ContextSatellite struct {
	ContextID   string // ID of the Context
	SatelliteID string // ID of the Satellite
}

// ContextTLE represents the relationship between Context and TLE in the domain layer.
type ContextTLE struct {
	ContextID string // ID of the Context
	TLEID     string // ID of the TLE
}

// ContextTile represents the relationship between a Context and a Tile in the domain layer.
type ContextTile struct {
	ContextID string // ID of the associated Context
	TileID    string // ID of the associated Tile
}

// ContextRepository defines the interface for context operations.
type ContextRepository interface {
	Save(ctx context.Context, context Context) error
	Update(ctx context.Context, context Context) error
	FindByID(ctx context.Context, id string) (Context, error)
	FindAll(ctx context.Context) ([]Context, error)
	DeleteByID(ctx context.Context, id string) error
	FindBySatelliteID(ctx context.Context, satelliteID string) ([]Context, error)
	AssignSatellite(ctx context.Context, contextID, satelliteID string) error
	RemoveSatellite(ctx context.Context, contextID, satelliteID string) error
}

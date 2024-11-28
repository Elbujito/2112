package domain

import (
	"context"

	"github.com/Elbujito/2112/pkg/fx/polygon"
	"github.com/google/uuid"
)

// TileRepository defines the interface for Tile repository operations.
type TileRepository interface {
	FindByQuadkey(ctx context.Context, key polygon.Quadkey) (*Tile, error)
	FindAll(ctx context.Context) ([]Tile, error)
	Save(ctx context.Context, tile Tile) error
	Update(ctx context.Context, tile Tile) error
	Upsert(ctx context.Context, tile Tile) error
	DeleteByQuadKey(ctx context.Context, key polygon.Quadkey) error
}

// Tile represents the domain entity Tile
type Tile struct {
	ID      string // Unique identifier
	Polygon polygon.Polygon
}

// NewTile constructor
func NewTile(polygon polygon.Polygon) Tile {
	return Tile{
		ID:      uuid.NewString(),
		Polygon: polygon,
	}
}

package domain

import (
	"context"
	"time"

	"github.com/Elbujito/2112/pkg/fx/polygon"
	"github.com/google/uuid"
)

// TileRepository defines the interface for Tile repository operations.
type TileRepository interface {
	FindByQuadkey(ctx context.Context, key polygon.Quadkey) (*Tile, error)                         // Find a tile by Quadkey
	FindBySpatialLocation(ctx context.Context, lat, lon float64) (*Tile, error)                    // Find a tile by spatial location
	FindTilesInRegion(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]Tile, error) // Find tiles intersecting a region
	FindAll(ctx context.Context) ([]Tile, error)                                                   // Retrieve all tiles
	Save(ctx context.Context, tile Tile) error                                                     // Save a new tile
	Update(ctx context.Context, tile Tile) error                                                   // Update an existing tile
	Upsert(ctx context.Context, tile Tile) error                                                   // Upsert (insert or update) a tile
	DeleteByQuadkey(ctx context.Context, key polygon.Quadkey) error                                // Delete a tile by Quadkey
	DeleteBySpatialLocation(ctx context.Context, lat, lon float64) error                           // Delete a tile by spatial location
}

// Tile represents the domain entity Tile
type Tile struct {
	ID        string // Unique identifier (UUID)
	CreatedAt time.Time
	UpdatedAt time.Time
	Quadkey   string          // Quadkey representing the tile
	ZoomLevel int             // Zoom level of the tile
	CenterLat float64         // Center latitude of the tile
	CenterLon float64         // Center longitude of the tile
	NbFaces   int             // Number of faces in the tile's geometry
	Radius    float64         // Radius of the tile (in meters or other unit)
	Vertices  []polygon.Point // Vertices representing the boundary of the tile
}

// NewTile constructor
// NewTile constructor
func NewTile(polygon polygon.Polygon) Tile {
	return Tile{
		ID:        uuid.NewString(), // Generate a new UUID
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Quadkey:   polygon.Center.Key(),     // Extract quadkey from the polygon center
		ZoomLevel: polygon.Center.Level,     // Use the zoom level from the center
		CenterLat: polygon.Center.Latitude,  // Use center latitude
		CenterLon: polygon.Center.Longitude, // Use center longitude
		NbFaces:   polygon.NbFaces,          // Number of faces in the tile
		Radius:    polygon.Radius,           // Tile radius
		Vertices:  polygon.Boundaries,       // Boundary vertices
	}
}
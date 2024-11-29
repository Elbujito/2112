package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TileSatelliteMappingRepository interface {
	FindByNoradIDAndTile(ctx context.Context, noradID string, tileID string) ([]TileSatelliteMapping, error)
	FindAll(ctx context.Context) ([]TileSatelliteMapping, error)
	Save(ctx context.Context, visibility TileSatelliteMapping) error
	Update(ctx context.Context, visibility TileSatelliteMapping) error
	Delete(ctx context.Context, id string) error
	SaveBatch(ctx context.Context, visibilities []TileSatelliteMapping) error
	FindAllVisibleTilesByNoradIDSortedByAOSTime(
		ctx context.Context,
		noradID string,
	) ([]TileSatelliteInfo, error)
}

// TileSatelliteMapping represents the domain entity TileSatelliteMapping
type TileSatelliteMapping struct {
	ID           string // Unique identifier
	CreatedAt    time.Time
	UpdatedAt    time.Time
	NoradID      string
	TileID       string
	Aos          time.Time
	MaxElevation float64
}

// NewVisibility constructor
func NewVisibility(noradID string,
	tileID string,
	startTime time.Time,
	maxElevation float64) TileSatelliteMapping {
	return TileSatelliteMapping{
		ID:           uuid.NewString(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		NoradID:      noradID,
		TileID:       tileID,
		Aos:          startTime,
		MaxElevation: maxElevation,
	}

}

// TileSatelliteInfo represents the aggregated data of a tile and satellite, sorted by AOS time.
type TileSatelliteInfo struct {
	TileID           string    // The ID of the tile
	TileQuadkey      string    // The Quadkey of the tile
	TileCenterLat    float64   // Latitude of the tile center
	TileCenterLon    float64   // Longitude of the tile center
	TileZoomLevel    int       // Zoom level of the tile
	SatelliteID      string    // The ID of the satellite
	SatelliteNoradID string    // The NORAD ID of the satellite
	AOS              time.Time // The Acquisition of Signal (AOS) time
}

package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MappingRepository interface {
	FindByNoradIDAndTile(ctx context.Context, noradID string, tileID string) ([]TileSatelliteMapping, error)
	FindAll(ctx context.Context) ([]TileSatelliteMapping, error)
	Save(ctx context.Context, visibility TileSatelliteMapping) error
	Update(ctx context.Context, visibility TileSatelliteMapping) error
	Delete(ctx context.Context, id string) error
	SaveBatch(ctx context.Context, visibilities []TileSatelliteMapping) error
	FindSatellitesForTiles(ctx context.Context, tileIDs []string) ([]Satellite, error)
	FindAllVisibleTilesByNoradIDSortedByAOSTime(
		ctx context.Context,
		noradID string,
	) ([]TileSatelliteInfo, error)
	ListSatellitesMappingWithPagination(ctx context.Context, page int, pageSize int, search *SearchRequest) ([]TileSatelliteInfo, int64, error)
}

// TileSatelliteMapping represents the domain entity TileSatelliteMapping
type TileSatelliteMapping struct {
	ID                    string // Unique identifier
	CreatedAt             time.Time
	UpdatedAt             time.Time
	NoradID               string
	TileID                string
	IntersectionLongitude float64
	IntersectionLatitude  float64
}

// NewMapping constructor
func NewMapping(noradID string,
	tileID string, intersection Point) TileSatelliteMapping {
	return TileSatelliteMapping{
		ID:                    uuid.NewString(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		NoradID:               noradID,
		TileID:                tileID,
		IntersectionLongitude: intersection.Longitude,
		IntersectionLatitude:  intersection.Latitude,
	}

}

// TileSatelliteInfo represents the aggregated data of a tile and satellite, sorted by AOS time.
type TileSatelliteInfo struct {
	TileID           string  // The ID of the tile
	TileQuadkey      string  // The Quadkey of the tile
	TileCenterLat    float64 // Latitude of the tile center
	TileCenterLon    float64 // Longitude of the tile center
	TileZoomLevel    int     // Zoom level of the tile
	SatelliteID      string  // The ID of the satellite
	SatelliteNoradID string  // The NORAD ID of the satellite
	Intersection     Point
}

type Point struct {
	Longitude float64
	Latitude  float64
}

package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/polygon"
)

// Tile Model
type Tile struct {
	ModelBase
	Quadkey        string  `gorm:"size:256;unique;not null"`   // Unique identifier for the tile (Quadkey)
	ZoomLevel      int     `gorm:"not null"`                   // Zoom level for the tile
	CenterLat      float64 `gorm:"not null"`                   // Center latitude of the tile
	CenterLon      float64 `gorm:"not null"`                   // Center longitude of the tile
	SpatialIndex   string  `gorm:"type:geometry;spatialIndex"` // Geometry column for spatial queries
	NbFaces        int     `gorm:"not null"`                   // Number of faces in the tile's shape
	Radius         float64 `gorm:"not null"`                   // Radius of the tile in meters
	BoundariesJSON string  `gorm:"type:json"`                  // Serialized JSON of the boundary vertices of the tile
}

// Validate validates the fields of the Tile model.
func (t *Tile) Validate() error {
	// Check required fields
	if t.Quadkey == "" {
		return errors.New("quadkey is required")
	}
	if t.ZoomLevel < 0 {
		return errors.New("zoom level must be non-negative")
	}
	if t.NbFaces <= 0 {
		return errors.New("number of faces must be greater than 0")
	}
	if t.Radius <= 0 {
		return errors.New("radius must be greater than 0")
	}
	if t.CenterLat < -90 || t.CenterLat > 90 {
		return errors.New("center latitude must be between -90 and 90")
	}
	if t.CenterLon < -180 || t.CenterLon > 180 {
		return errors.New("center longitude must be between -180 and 180")
	}

	// Validate BoundariesJSON if present
	if t.BoundariesJSON != "" {
		var temp []polygon.Point
		if err := json.Unmarshal([]byte(t.BoundariesJSON), &temp); err != nil {
			return fmt.Errorf("boundaries JSON must be valid: %w", err)
		}
	}

	return nil
}

// MapToForm converts the Tile model to a TileForm structure.
func (t *Tile) MapToForm() *TileForm {
	return &TileForm{
		FormBase: FormBase{
			ID:        t.ID,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
		Quadkey:   t.Quadkey,
		ZoomLevel: t.ZoomLevel,
		CenterLat: t.CenterLat,
		CenterLon: t.CenterLon,
	}
}

// MapFromDomain converts a domain.Tile to a models.Tile.
func MapFromDomain(domainTile domain.Tile) Tile {
	// Serialize boundaries
	boundariesJSON, err := json.Marshal(domainTile.Vertices)
	if err != nil {
		boundariesJSON = []byte("[]") // Default to empty array on failure
	}

	// Convert to model
	return Tile{
		ModelBase: ModelBase{
			ID: domainTile.ID,
		},
		Quadkey:        domainTile.Quadkey,
		ZoomLevel:      domainTile.ZoomLevel,
		CenterLat:      domainTile.CenterLat,
		CenterLon:      domainTile.CenterLon,
		NbFaces:        domainTile.NbFaces,
		Radius:         domainTile.Radius,
		BoundariesJSON: string(boundariesJSON),
		SpatialIndex:   createGeometryFromBoundaries(domainTile.Vertices), // Generate spatial index geometry
	}
}

// MapToDomain converts a models.Tile to a domain.Tile.
func MapToDomain(modelTile Tile) domain.Tile {
	var boundaries []polygon.Point

	// Deserialize boundaries
	err := json.Unmarshal([]byte(modelTile.BoundariesJSON), &boundaries)
	if err != nil {
		boundaries = nil // Default to nil if deserialization fails
	}

	// Convert to domain
	return domain.Tile{
		ID:        modelTile.ID,
		Quadkey:   modelTile.Quadkey,
		ZoomLevel: modelTile.ZoomLevel,
		CenterLat: modelTile.CenterLat,
		CenterLon: modelTile.CenterLon,
		NbFaces:   modelTile.NbFaces,
		Radius:    modelTile.Radius,
		Vertices:  boundaries,
	}
}

// createGeometryFromBoundaries generates WKT (Well-Known Text) representation of the boundaries for spatial indexing.
func createGeometryFromBoundaries(vertices []polygon.Point) string {
	if len(vertices) == 0 {
		return "POLYGON(EMPTY)" // Return an empty polygon if no vertices are provided
	}

	wkt := "POLYGON(("
	for i, vertex := range vertices {
		if i > 0 {
			wkt += ","
		}
		wkt += fmt.Sprintf("%f %f", vertex.Longitude, vertex.Latitude)
	}
	wkt += fmt.Sprintf(",%f %f", vertices[0].Longitude, vertices[0].Latitude) // Close the polygon
	wkt += "))"
	return wkt
}

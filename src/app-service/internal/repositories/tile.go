package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/data"
	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"gorm.io/gorm"
)

type TileRepository struct {
	db *data.Database
}

// NewTileRepository creates a new instance of TileRepository.
func NewTileRepository(db *data.Database) domain.TileRepository {
	return &TileRepository{db: db}
}

// FindByQuadkey retrieves a Tile by its quadkey.
func (r *TileRepository) FindByQuadkey(ctx context.Context, quadkey string) (*domain.Tile, error) {
	var tile models.Tile
	result := r.db.DbHandler.Where("quadkey = ?", quadkey).First(&tile)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	tileMapped := models.MapToDomain(tile)
	return &tileMapped, nil
}

// FindBySpatialLocation retrieves a Tile by a geographical location using spatial indexing.
func (r *TileRepository) FindBySpatialLocation(ctx context.Context, lat, lon float64) (*domain.Tile, error) {
	var tile models.Tile
	result := r.db.DbHandler.Raw(`
		SELECT *
		FROM tiles
		WHERE ST_Contains(spatial_index, ST_SetSRID(ST_Point(?, ?), 4326))
		LIMIT 1
	`, lon, lat).Scan(&tile)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	tileMapped := models.MapToDomain(tile)
	return &tileMapped, nil
}

// FindAll retrieves all Tiles.
func (r *TileRepository) FindAll(ctx context.Context) ([]domain.Tile, error) {
	var tiles []models.Tile
	result := r.db.DbHandler.Find(&tiles)
	if result.Error != nil {
		return nil, result.Error
	}

	// Map models to domain
	var domainTiles []domain.Tile
	for _, t := range tiles {
		domainTiles = append(domainTiles, models.MapToDomain(t))
	}
	return domainTiles, nil
}

// Save creates a new Tile record.
func (r *TileRepository) Save(ctx context.Context, tile domain.Tile) error {
	modelTile := models.MapFromDomain(tile)
	return r.db.DbHandler.Create(&modelTile).Error
}

// Update modifies an existing Tile record.
func (r *TileRepository) Update(ctx context.Context, tile domain.Tile) error {
	modelTile := models.MapFromDomain(tile)
	return r.db.DbHandler.Save(&modelTile).Error
}

// DeleteByQuadkey removes a Tile record by its quadkey.
func (r *TileRepository) DeleteByQuadkey(ctx context.Context, key string) error {
	return r.db.DbHandler.Where("quadkey = ?", key).Delete(&models.Tile{}).Error
}

// Upsert inserts or updates a Tile record in the database.
func (r *TileRepository) Upsert(ctx context.Context, tile domain.Tile) error {
	// Check for an existing tile using the spatial location or quadkey
	existingTile, err := r.FindByQuadkey(ctx, tile.Quadkey)
	if err != nil {
		return err
	}

	if existingTile != nil {
		// Update if the tile exists
		return r.Update(ctx, tile)
	}

	// Save if the tile doesn't exist
	return r.Save(ctx, tile)
}

// DeleteBySpatialLocation removes a Tile record by its geographical location.
func (r *TileRepository) DeleteBySpatialLocation(ctx context.Context, lat, lon float64) error {
	var tile models.Tile
	result := r.db.DbHandler.Raw(`
		SELECT *
		FROM tiles
		WHERE ST_Contains(spatial_index, ST_SetSRID(ST_Point(?, ?), 4326))
		LIMIT 1
	`, lon, lat).Scan(&tile)
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("tile not found")
	}

	return r.db.DbHandler.Delete(&tile).Error
}

// FindTilesInRegion retrieves tiles that intersect a given bounding box.
func (r *TileRepository) FindTilesInRegion(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]domain.Tile, error) {
	var tiles []models.Tile
	result := r.db.DbHandler.Raw(`
		SELECT *
		FROM tiles
		WHERE ST_Intersects(spatial_index, ST_MakeEnvelope(?, ?, ?, ?, 4326))
	`, minLon, minLat, maxLon, maxLat).Scan(&tiles)
	if result.Error != nil {
		return nil, result.Error
	}

	// Map models to domain
	var domainTiles []domain.Tile
	for _, t := range tiles {
		domainTiles = append(domainTiles, models.MapToDomain(t))
	}
	return domainTiles, nil
}

func (r *TileRepository) FindTilesVisibleFromLine(ctx context.Context, sat domain.Satellite, points []domain.SatellitePosition) ([]domain.TileSatelliteMapping, error) {
	if len(points) < 2 {
		return nil, fmt.Errorf("at least two points are required to create a line")
	}

	// Construct a WKT (Well-Known Text) representation of the LINESTRING
	wktPoints := make([]string, len(points))
	for i, point := range points {
		wktPoints[i] = fmt.Sprintf("%f %f", point.Longitude, point.Latitude)
	}
	lineString := fmt.Sprintf("LINESTRING(%s)", strings.Join(wktPoints, ", "))

	// Query to find intersecting tiles and calculate interest points
	query := `
        WITH line_geom AS (
            SELECT ST_GeomFromText(?, 4326) AS geom
        )
        SELECT 
            tiles.*,
            ST_AsText(ST_PointOnSurface(ST_Intersection(line_geom.geom, spatial_index))) AS intersection_geom
        FROM tiles, line_geom
        WHERE ST_Intersects(spatial_index, line_geom.geom)
    `

	// Execute the query
	var results []struct {
		models.Tile
		IntersectionGeom string `gorm:"column:intersection_geom"`
	}
	result := r.db.DbHandler.Raw(query, lineString).Scan(&results)
	if result.Error != nil {
		log.Printf("Error executing query: %v\n", result.Error)
		return nil, result.Error
	}

	if len(results) == 0 {
		log.Printf("No tiles intersect the line for satellite %s\n", sat.NoradID)
		return nil, nil
	}

	// Map results to domain with interest points
	var mappings []domain.TileSatelliteMapping
	for _, res := range results {
		tile := models.MapToDomain(res.Tile)

		// Parse intersection geometry
		interestPoint, err := parseIntersectionGeometry(res.IntersectionGeom)
		if err != nil {
			log.Printf("Failed to parse intersection geometry for TileID %s: %v\n", tile.ID, err)
			continue // Skip invalid intersections
		}

		mapping := domain.TileSatelliteMapping{
			NoradID:               sat.NoradID,
			TileID:                tile.ID,
			IntersectionLongitude: interestPoint.Longitude,
			IntersectionLatitude:  interestPoint.Latitude,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

func parseIntersectionGeometry(wkt string) (domain.Point, error) {
	if wkt == "" || !strings.HasPrefix(wkt, "POINT(") || !strings.HasSuffix(wkt, ")") {
		return domain.Point{}, fmt.Errorf("invalid WKT format: %s", wkt)
	}

	coordinates := strings.TrimPrefix(strings.TrimSuffix(wkt, ")"), "POINT(")
	parts := strings.Fields(coordinates)
	if len(parts) != 2 {
		return domain.Point{}, fmt.Errorf("invalid WKT coordinates: %s", coordinates)
	}

	longitude, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return domain.Point{}, fmt.Errorf("error parsing longitude: %w", err)
	}

	latitude, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return domain.Point{}, fmt.Errorf("error parsing latitude: %w", err)
	}

	return domain.Point{
		Longitude: longitude,
		Latitude:  latitude,
	}, nil
}

// FindTilesIntersectingLocation retrieves all tiles that intersect the user's location.
func (r *TileRepository) FindTilesIntersectingLocation(ctx context.Context, lat, lon, radius float64) ([]domain.Tile, error) {
	var tiles []models.Tile

	// Query tiles that intersect the user's location with the given radius
	result := r.db.DbHandler.Raw(`
		SELECT *
		FROM tiles
		WHERE ST_DWithin(
			spatial_index,
			ST_MakePoint(?, ?)::geography,
			?
		)
	`, lon, lat, radius).Scan(&tiles)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(tiles) == 0 {
		return nil, errors.New("no tiles intersecting the user's location")
	}

	// Map tiles to domain objects
	var domainTiles []domain.Tile
	for _, tile := range tiles {
		domainTiles = append(domainTiles, models.MapToDomain(tile))
	}

	return domainTiles, nil
}

// SaveBatch allows batch insertion of tiles for optimized performance.
func (r *TileRepository) SaveBatch(ctx context.Context, tiles []domain.Tile) error {
	modelTiles := make([]models.Tile, len(tiles))
	for i, t := range tiles {
		modelTiles[i] = models.MapFromDomain(t)
	}
	return r.db.DbHandler.Create(&modelTiles).Error
}

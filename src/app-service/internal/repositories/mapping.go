package repository

import (
	"context"

	"github.com/Elbujito/2112/src/app-service/internal/data"
	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"gorm.io/gorm"
)

type TileSatelliteMappingRepository struct {
	db *data.Database
}

// NewTileSatelliteMappingRepository creates a new instance of VisibilityRepository.
func NewTileSatelliteMappingRepository(db *data.Database) domain.MappingRepository {
	return &TileSatelliteMappingRepository{db: db}
}

// FindByNoradIDAndTile retrieves Visibility records by NORAD ID and Tile ID.
func (r *TileSatelliteMappingRepository) FindByNoradIDAndTile(ctx context.Context, noradID, tileID string) ([]domain.TileSatelliteMapping, error) {
	var visibilities []domain.TileSatelliteMapping
	result := r.db.DbHandler.Where("norad_id = ? AND tile_id = ?", noradID, tileID).Find(&visibilities)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return visibilities, nil
}

// FindAll retrieves all Visibility records.
func (r *TileSatelliteMappingRepository) FindAll(ctx context.Context) ([]domain.TileSatelliteMapping, error) {
	var visibilities []domain.TileSatelliteMapping
	result := r.db.DbHandler.Find(&visibilities)
	if result.Error != nil {
		return nil, result.Error
	}
	return visibilities, nil
}

// Save creates a new Visibility record.
func (r *TileSatelliteMappingRepository) Save(ctx context.Context, visibility domain.TileSatelliteMapping) error {
	// Avoid duplicate visibility records with unique constraints if applicable
	return r.db.DbHandler.Create(&visibility).Error
}

// SaveBatch creates multiple Visibility records in a batch operation.
func (r *TileSatelliteMappingRepository) SaveBatch(ctx context.Context, visibilities []domain.TileSatelliteMapping) error {
	// Use GORM batch insert
	return r.db.DbHandler.CreateInBatches(visibilities, 100).Error
}

// Update modifies an existing Visibility record.
func (r *TileSatelliteMappingRepository) Update(ctx context.Context, visibility domain.TileSatelliteMapping) error {
	return r.db.DbHandler.Save(&visibility).Error
}

// UpdateBatch updates multiple Visibility records in a batch.
func (r *TileSatelliteMappingRepository) UpdateBatch(ctx context.Context, visibilities []domain.TileSatelliteMapping) error {
	tx := r.db.DbHandler.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, visibility := range visibilities {
		if err := tx.Save(&visibility).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// Delete removes a Visibility record by its ID.
func (r *TileSatelliteMappingRepository) Delete(ctx context.Context, id string) error {
	return r.db.DbHandler.Where("id = ?", id).Delete(&domain.TileSatelliteMapping{}).Error
}

// DeleteBatch deletes multiple Visibility records by their IDs.
func (r *TileSatelliteMappingRepository) DeleteBatch(ctx context.Context, ids []string) error {
	return r.db.DbHandler.Where("id IN ?", ids).Delete(&domain.TileSatelliteMapping{}).Error
}

// FindByNoradID retrieves all Visibility records for a given NORAD ID.
func (r *TileSatelliteMappingRepository) FindByNoradID(ctx context.Context, noradID string) ([]domain.TileSatelliteMapping, error) {
	var visibilities []domain.TileSatelliteMapping
	result := r.db.DbHandler.Where("norad_id = ?", noradID).Find(&visibilities)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return visibilities, nil
}

// FindByTileID retrieves all Visibility records for a given Tile ID.
func (r *TileSatelliteMappingRepository) FindByTileID(ctx context.Context, tileID string) ([]domain.TileSatelliteMapping, error) {
	var visibilities []domain.TileSatelliteMapping
	result := r.db.DbHandler.Where("tile_id = ?", tileID).Find(&visibilities)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return visibilities, nil
}

// FindAllVisibleTilesByNoradIDSortedByAOSTime retrieves all tiles visible for a given NORAD ID,
// and aggregates them with the satellite information, sorted by AOS time.
func (r *TileSatelliteMappingRepository) FindAllVisibleTilesByNoradIDSortedByAOSTime(
	ctx context.Context,
	noradID string,
) ([]domain.TileSatelliteInfo, error) {
	var tileMappings []domain.TileSatelliteMapping
	result := r.db.DbHandler.
		Where("norad_id = ?", noradID).
		Order("aos ASC"). // Order by AOS in ascending order (earliest AOS first)
		Find(&tileMappings)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	// Collect the tile IDs from the tileMappings
	tileIDs := make([]string, len(tileMappings))
	for i, mapping := range tileMappings {
		tileIDs[i] = mapping.TileID
	}

	var modelTiles []models.Tile
	// Fetch the corresponding tiles by their IDs
	result = r.db.DbHandler.
		Where("id IN ?", tileIDs).
		Find(&modelTiles)

	if result.Error != nil {
		return nil, result.Error
	}

	// Create a map to hold the TileSatelliteInfo objects
	tileSatelliteInfos := make([]domain.TileSatelliteInfo, len(tileMappings))

	// Loop over the tileMappings and aggregate the data
	for i, mapping := range tileMappings {
		// Map the model tile to domain tile
		modelTile := findTileByID(modelTiles, mapping.TileID)

		tileSatelliteInfos[i] = domain.TileSatelliteInfo{
			TileID:           modelTile.ID,
			TileQuadkey:      modelTile.Quadkey,
			TileCenterLat:    modelTile.CenterLat,
			TileCenterLon:    modelTile.CenterLon,
			TileZoomLevel:    modelTile.ZoomLevel,
			SatelliteID:      mapping.NoradID,
			SatelliteNoradID: mapping.NoradID,
			AOS:              mapping.Aos,
		}
	}

	return tileSatelliteInfos, nil
}

// Helper function to find the tile by ID in the fetched model tiles
func findTileByID(modelTiles []models.Tile, tileID string) models.Tile {
	for _, tile := range modelTiles {
		if tile.ID == tileID {
			return tile
		}
	}
	return models.Tile{}
}

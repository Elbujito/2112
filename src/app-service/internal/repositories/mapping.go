package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Elbujito/2112/src/app-service/internal/data"
	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"github.com/google/uuid"
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

// SaveBatch creates multiple Visibility records in a batch operation with "ON CONFLICT DO NOTHING".
func (r *TileSatelliteMappingRepository) SaveBatch(ctx context.Context, visibilities []domain.TileSatelliteMapping) error {
	if len(visibilities) == 0 {
		return nil // No records to insert
	}

	const batchSize = 100 // Define batch size to limit query size

	// Process in batches
	for i := 0; i < len(visibilities); i += batchSize {
		end := i + batchSize
		if end > len(visibilities) {
			end = len(visibilities)
		}
		batch := visibilities[i:end]

		// Construct query dynamically
		placeholders := make([]string, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*5) // 5 fields per record (including ID)

		for j, v := range batch {
			if v.ID == "" {
				v.ID = uuid.NewString() // Generate a new UUID if ID is not set
			}
			placeholders[j] = "(?, ?, ?, ?, ?)"
			valueArgs = append(valueArgs, v.ID, v.NoradID, v.TileID, v.CreatedAt, v.UpdatedAt)
		}

		query := `
            INSERT INTO tile_satellite_mappings (id, norad_id, tile_id, created_at, updated_at)
            VALUES %s
            ON CONFLICT ON CONSTRAINT unique_norad_tile_mapping DO NOTHING
        `
		formattedQuery := fmt.Sprintf(query, strings.Join(placeholders, ", "))

		// Execute batch insert
		if err := r.db.DbHandler.WithContext(ctx).Exec(formattedQuery, valueArgs...).Error; err != nil {
			return fmt.Errorf("failed to save batch: %w", err)
		}
	}

	return nil
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

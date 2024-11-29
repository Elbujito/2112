package repository

import (
	"context"

	"github.com/Elbujito/2112/internal/data"
	"github.com/Elbujito/2112/internal/domain"
	"gorm.io/gorm"
)

type VisibilityRepository struct {
	db *data.Database
}

// NewVisibilityRepository creates a new instance of VisibilityRepository.
func NewVisibilityRepository(db *data.Database) domain.VisibilityRepository {
	return &VisibilityRepository{db: db}
}

// FindByNoradIDAndTile retrieves Visibility records by NORAD ID and Tile ID.
func (r *VisibilityRepository) FindByNoradIDAndTile(ctx context.Context, noradID, tileID string) ([]domain.Visibility, error) {
	var visibilities []domain.Visibility
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
func (r *VisibilityRepository) FindAll(ctx context.Context) ([]domain.Visibility, error) {
	var visibilities []domain.Visibility
	result := r.db.DbHandler.Find(&visibilities)
	if result.Error != nil {
		return nil, result.Error
	}
	return visibilities, nil
}

// Save creates a new Visibility record.
func (r *VisibilityRepository) Save(ctx context.Context, visibility domain.Visibility) error {
	// Avoid duplicate visibility records with unique constraints if applicable
	return r.db.DbHandler.Create(&visibility).Error
}

// SaveBatch creates multiple Visibility records in a batch operation.
func (r *VisibilityRepository) SaveBatch(ctx context.Context, visibilities []domain.Visibility) error {
	// Use GORM batch insert
	return r.db.DbHandler.CreateInBatches(visibilities, 100).Error
}

// Update modifies an existing Visibility record.
func (r *VisibilityRepository) Update(ctx context.Context, visibility domain.Visibility) error {
	return r.db.DbHandler.Save(&visibility).Error
}

// UpdateBatch updates multiple Visibility records in a batch.
func (r *VisibilityRepository) UpdateBatch(ctx context.Context, visibilities []domain.Visibility) error {
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
func (r *VisibilityRepository) Delete(ctx context.Context, id string) error {
	return r.db.DbHandler.Where("id = ?", id).Delete(&domain.Visibility{}).Error
}

// DeleteBatch deletes multiple Visibility records by their IDs.
func (r *VisibilityRepository) DeleteBatch(ctx context.Context, ids []string) error {
	return r.db.DbHandler.Where("id IN ?", ids).Delete(&domain.Visibility{}).Error
}

// FindByNoradID retrieves all Visibility records for a given NORAD ID.
func (r *VisibilityRepository) FindByNoradID(ctx context.Context, noradID string) ([]domain.Visibility, error) {
	var visibilities []domain.Visibility
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
func (r *VisibilityRepository) FindByTileID(ctx context.Context, tileID string) ([]domain.Visibility, error) {
	var visibilities []domain.Visibility
	result := r.db.DbHandler.Where("tile_id = ?", tileID).Find(&visibilities)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return visibilities, nil
}

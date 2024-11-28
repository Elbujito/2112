package repository

import (
	"context"

	"github.com/Elbujito/2112/internal/data"
	"github.com/Elbujito/2112/internal/domain"
)

type VisibilityRepository struct {
	db *data.Database
}

// NewVisibilityRepository creates a new instance of VisibilityRepository.
func NewVisibilityRepository(db *data.Database) domain.VisibilityRepository {
	return &VisibilityRepository{db: db}
}

// FindByNoradIDAndTile retrieves Visibility records by NORAD ID and Tile ID.
func (r *VisibilityRepository) FindByNoradIDAndTile(ctx context.Context, noradID string, tileID string) ([]domain.Visibility, error) {
	var visibilities []domain.Visibility
	result := r.db.DbHandler.Where("norad_id = ? AND tile_id = ?", noradID, tileID).Find(&visibilities)
	return visibilities, result.Error
}

// FindAll retrieves all Visibility records.
func (r *VisibilityRepository) FindAll(ctx context.Context) ([]domain.Visibility, error) {
	var visibilities []domain.Visibility
	result := r.db.DbHandler.Find(&visibilities)
	return visibilities, result.Error
}

// Save creates a new Visibility record.
func (r *VisibilityRepository) Save(ctx context.Context, visibility domain.Visibility) error {
	return r.db.DbHandler.Create(&visibility).Error
}

// Update modifies an existing Visibility record.
func (r *VisibilityRepository) Update(ctx context.Context, visibility domain.Visibility) error {
	return r.db.DbHandler.Save(&visibility).Error
}

// Delete removes a Visibility record by its ID.
func (r *VisibilityRepository) Delete(ctx context.Context, id string) error {
	return r.db.DbHandler.Where("id = ?", id).Delete(&domain.Visibility{}).Error
}

package repository

import (
	"context"
	"errors"

	"github.com/Elbujito/2112/internal/data"
	"github.com/Elbujito/2112/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SatelliteRepository struct {
	db *data.Database
}

// NewSatelliteRepository creates a new instance of SatelliteRepository.
func NewSatelliteRepository(db *data.Database) domain.SatelliteRepository {
	return &SatelliteRepository{db: db}
}

// FindByNoradID retrieves a satellite by its NORAD ID.
func (r *SatelliteRepository) FindByNoradID(ctx context.Context, noradID string) (domain.Satellite, error) {
	var satellite domain.Satellite
	result := r.db.DbHandler.Where("norad_id = ?", noradID).First(&satellite)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.Satellite{}, nil
	}
	return satellite, result.Error
}

// FindAll retrieves all satellites.
func (r *SatelliteRepository) FindAll(ctx context.Context) ([]domain.Satellite, error) {
	var satellites []domain.Satellite
	result := r.db.DbHandler.Find(&satellites)
	return satellites, result.Error
}

// Save creates a new satellite record.
func (r *SatelliteRepository) Save(ctx context.Context, satellite domain.Satellite) error {
	return r.db.DbHandler.Create(&satellite).Error
}

// Update modifies an existing satellite record.
func (r *SatelliteRepository) Update(ctx context.Context, satellite domain.Satellite) error {
	return r.db.DbHandler.Save(&satellite).Error
}

// DeleteByNoradID removes a satellite record by its NoradID.
func (r *SatelliteRepository) DeleteByNoradID(ctx context.Context, noradID string) error {
	return r.db.DbHandler.Where("noradID = ?", noradID).Delete(&domain.Satellite{}).Error
}

// SaveBatch performs a batch upsert (insert or update) for satellites.
func (r *SatelliteRepository) SaveBatch(ctx context.Context, satellites []domain.Satellite) error {
	if len(satellites) == 0 {
		return nil // Nothing to save
	}

	// Use Gen's support for ON CONFLICT upsert
	return r.db.DbHandler.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "norad_id"}}, // Define the unique constraint column
			UpdateAll: true,                                // Update all fields in case of conflict
		}).
		CreateInBatches(satellites, 100).Error // Batch size: 100
}

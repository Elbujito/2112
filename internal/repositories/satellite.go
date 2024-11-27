package repository

import (
	"errors"

	"github.com/Elbujito/2112/internal/domain"
	"gorm.io/gorm"
)

// satelliteRepositoryImpl is the concrete implementation of the SatelliteRepository interface.
type satelliteRepositoryImpl struct {
	db *gorm.DB
}

// NewSatelliteRepository creates a new instance of SatelliteRepository.
func NewSatelliteRepository(db *gorm.DB) domain.SatelliteRepository {
	return &satelliteRepositoryImpl{db: db}
}

// FindByNoradID retrieves a satellite by its NORAD ID.
func (r *satelliteRepositoryImpl) FindByNoradID(noradID string) (domain.Satellite, error) {
	var satellite domain.Satellite
	result := r.db.Where("norad_id = ?", noradID).First(&satellite)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.Satellite{}, nil
	}
	return satellite, result.Error
}

// Find retrieves a satellite by its ID.
func (r *satelliteRepositoryImpl) Find(id string) (domain.Satellite, error) {
	var satellite domain.Satellite
	result := r.db.Where("id = ?", id).First(&satellite)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.Satellite{}, nil
	}
	return satellite, result.Error
}

// FindAll retrieves all satellites.
func (r *satelliteRepositoryImpl) FindAll() ([]domain.Satellite, error) {
	var satellites []domain.Satellite
	result := r.db.Find(&satellites)
	return satellites, result.Error
}

// Save creates a new satellite record.
func (r *satelliteRepositoryImpl) Save(satellite domain.Satellite) error {
	return r.db.Create(&satellite).Error
}

// Update modifies an existing satellite record.
func (r *satelliteRepositoryImpl) Update(satellite domain.Satellite) error {
	return r.db.Save(&satellite).Error
}

// Delete removes a satellite record by its ID.
func (r *satelliteRepositoryImpl) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Satellite{}).Error
}

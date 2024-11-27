package repository

import (
	"errors"

	"github.com/Elbujito/2112/internal/data/models"
	"github.com/Elbujito/2112/internal/domain"
	"gorm.io/gorm"
)

// tleRepositoryImpl is the concrete implementation of the TLERepository interface.
type tleRepositoryImpl struct {
	db *gorm.DB
}

// NewTLERepository creates a new instance of TLERepository.
func NewTLERepository(db *gorm.DB) domain.TLERepository {
	return &tleRepositoryImpl{db: db}
}

// mapToDomainTLE converts a models.TLE to a domain.TLE.
func mapToDomainTLE(model models.TLE) domain.TLE {
	return domain.TLE{
		ID:      model.ID,
		NoradID: model.NoradID,
		Line1:   model.Line1,
		Line2:   model.Line2,
		Epoch:   model.Epoch,
	}
}

// mapToModelTLE converts a domain.TLE to a models.TLE.
func mapToModelTLE(domainTLE domain.TLE) models.TLE {
	return models.TLE{
		ModelBase: models.ModelBase{ID: domainTLE.ID},
		NoradID:   domainTLE.NoradID,
		Line1:     domainTLE.Line1,
		Line2:     domainTLE.Line2,
		Epoch:     domainTLE.Epoch,
	}
}

// FindByNoradID retrieves all TLEs for a given NORAD ID.
func (r *tleRepositoryImpl) FindByNoradID(noradID string) ([]domain.TLE, error) {
	var modelTLEs []models.TLE
	result := r.db.Where("norad_id = ?", noradID).Find(&modelTLEs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	// Map models to domain
	var domainTLEs []domain.TLE
	for _, modelTLE := range modelTLEs {
		domainTLEs = append(domainTLEs, mapToDomainTLE(modelTLE))
	}
	return domainTLEs, nil
}

// FindAll retrieves all TLE records from the database.
func (r *tleRepositoryImpl) FindAll() ([]domain.TLE, error) {
	var modelTLEs []models.TLE
	result := r.db.Find(&modelTLEs)
	if result.Error != nil {
		return nil, result.Error
	}

	// Map models to domain
	var domainTLEs []domain.TLE
	for _, modelTLE := range modelTLEs {
		domainTLEs = append(domainTLEs, mapToDomainTLE(modelTLE))
	}
	return domainTLEs, nil
}

// Save inserts a new TLE record into the database.
func (r *tleRepositoryImpl) Save(tle domain.TLE) error {
	modelTLE := mapToModelTLE(tle)
	return r.db.Create(&modelTLE).Error
}

// Update modifies an existing TLE record in the database.
func (r *tleRepositoryImpl) Update(tle domain.TLE) error {
	modelTLE := mapToModelTLE(tle)
	return r.db.Save(&modelTLE).Error
}

// Upsert inserts or updates a TLE record in the database.
func (r *tleRepositoryImpl) Upsert(tle domain.TLE) error {
	existingTLEs, err := r.FindByNoradID(tle.NoradID)
	if err != nil {
		return err
	}
	if len(existingTLEs) > 0 {
		// Update the first matching record
		existingTLE := existingTLEs[0]
		existingTLE.Line1 = tle.Line1
		existingTLE.Line2 = tle.Line2
		existingTLE.Epoch = tle.Epoch
		return r.Update(existingTLE)
	}
	// Insert as a new record
	return r.Save(tle)
}

// Delete removes a TLE record by its ID.
func (r *tleRepositoryImpl) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.TLE{}).Error
}

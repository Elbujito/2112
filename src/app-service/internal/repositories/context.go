package repository

import (
	"context"

	"github.com/Elbujito/2112/src/app-service/internal/data"
	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"gorm.io/gorm"
)

// ContextRepository manages context data access.
type ContextRepository struct {
	db *data.Database
}

// NewContextRepository creates a new ContextRepository instance.
func NewContextRepository(db *data.Database) domain.ContextRepository {
	return &ContextRepository{db: db}
}

// Save creates a new context record.
func (r *ContextRepository) Save(ctx context.Context, context domain.Context) error {
	model := models.MapToContextModel(context)
	return r.db.DbHandler.Create(&model).Error
}

// Update modifies an existing context record.
func (r *ContextRepository) Update(ctx context.Context, context domain.Context) error {
	model := models.MapToContextModel(context)
	return r.db.DbHandler.Save(&model).Error
}

// FindByID retrieves a context by its ID.
func (r *ContextRepository) FindByID(ctx context.Context, id string) (domain.Context, error) {
	var model models.Context
	result := r.db.DbHandler.First(&model, "id = ? AND deleted_at IS NULL", id)
	if result.Error != nil {
		return domain.Context{}, result.Error
	}
	return models.MapToContextDomain(model), nil
}

// FindAll retrieves all contexts.
func (r *ContextRepository) FindAll(ctx context.Context) ([]domain.Context, error) {
	var results []models.Context
	result := r.db.DbHandler.Find(&results, "deleted_at IS NULL")
	if result.Error != nil {
		return nil, result.Error
	}

	var contexts []domain.Context
	for _, model := range results {
		context := models.MapToContextDomain(model)
		contexts = append(contexts, context)
	}
	return contexts, nil
}

// DeleteByID marks a context record as deleted.
func (r *ContextRepository) DeleteByID(ctx context.Context, id string) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindBySatelliteID retrieves contexts associated with a satellite.
func (r *ContextRepository) FindBySatelliteID(ctx context.Context, satelliteID string) ([]domain.Context, error) {
	var query []models.Context
	result := r.db.DbHandler.Table("contexts").
		Joins("JOIN context_satellites ON contexts.id = context_satellites.context_id").
		Where("context_satellites.satellite_id = ?", satelliteID).
		Find(&query)

	if result.Error != nil {
		return nil, result.Error
	}

	var contexts []domain.Context
	for _, model := range query {
		contexts = append(contexts, models.MapToContextDomain(model))
	}
	return contexts, nil
}

// AssignSatellite associates a satellite with a context.
func (r *ContextRepository) AssignSatellite(ctx context.Context, contextID, satelliteID string) error {
	contextSatellite := models.ContextSatellite{
		ContextID:   contextID,
		SatelliteID: satelliteID,
	}
	return r.db.DbHandler.Create(&contextSatellite).Error
}

// RemoveSatellite removes the association between a satellite and a context.
func (r *ContextRepository) RemoveSatellite(ctx context.Context, contextID, satelliteID string) error {
	return r.db.DbHandler.Where("context_id = ? AND satellite_id = ?", contextID, satelliteID).
		Delete(&models.ContextSatellite{}).Error
}

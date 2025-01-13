package services

import (
	"context"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

// ContextService definition
type ContextService struct {
	repo domain.GameContextRepository
}

// NewContextService creates a new instance of ContextService.
func NewContextService(repo domain.GameContextRepository) ContextService {
	return ContextService{repo: repo}
}

// Create creates a new GameContext.
func (c *ContextService) Create(ctx context.Context, context domain.GameContext) (domain.GameContext, error) {
	err := c.repo.Save(ctx, context)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// Update updates an existing GameContext.
func (c *ContextService) Update(ctx context.Context, context domain.GameContext) (domain.GameContext, error) {
	err := c.repo.Update(ctx, context)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// GetByUniqueName retrieves a GameContext by its unique name.
func (c *ContextService) GetByUniqueName(ctx context.Context, name domain.GameContextName) (domain.GameContext, error) {
	context, err := c.repo.FindByUniqueName(ctx, name)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// DeleteByUniqueName deletes a GameContext by its unique name.
func (c *ContextService) DeleteByUniqueName(ctx context.Context, name domain.GameContextName) error {
	err := c.repo.DeleteByUniqueName(ctx, string(name))
	if err != nil {
		return err
	}
	return nil
}

// ActiveContext activates a GameContext by its unique name.
func (c *ContextService) ActiveContext(ctx context.Context, name domain.GameContextName) error {
	err := c.repo.ActivateContext(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

// DisableContext deactivates a GameContext by its unique name.
func (c *ContextService) DisableContext(ctx context.Context, name domain.GameContextName) error {
	err := c.repo.DesactiveContext(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

// GetActiveContext retrieves the currently active GameContext.
func (c *ContextService) GetActiveContext(ctx context.Context) (domain.GameContext, error) {
	context, err := c.repo.GetActiveContext(ctx)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// GetAllContexts retrieves all GameContexts.
func (c *ContextService) GetAllContexts(ctx context.Context) ([]domain.GameContext, error) {
	contexts, err := c.repo.FindAll(ctx)
	if err != nil {
		return []domain.GameContext{}, err
	}
	return contexts, nil
}

// FindBySatelliteID retrieves all GameContexts associated with a specific satellite.
func (c *ContextService) FindBySatelliteID(ctx context.Context, satelliteID domain.SatelliteID) (domain.GameContext, error) {
	context, err := c.repo.FindActiveBySatelliteID(ctx, satelliteID)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// AssignSatellite associates a satellite with a GameContext.
func (c *ContextService) AssignSatellite(ctx context.Context, name domain.GameContextName, satelliteID domain.SatelliteID) error {
	err := c.repo.AssignSatellite(ctx, name, satelliteID)
	if err != nil {
		return err
	}
	return nil
}

// RemoveSatellite removes the association between a satellite and a GameContext.
func (c *ContextService) RemoveSatellite(ctx context.Context, name domain.GameContextName, satelliteID domain.SatelliteID) error {
	err := c.repo.RemoveSatellite(ctx, name, satelliteID)
	if err != nil {
		return err
	}
	return nil
}

// FindAllWithPagination retrieves all contexts with pagination and optional filtering by name or description.
func (c *ContextService) FindAllWithPagination(ctx context.Context, page int, pageSize int, wildcard string) ([]domain.GameContext, error) {
	contexts, err := c.repo.FindAllWithPagination(ctx, page, pageSize, wildcard)
	if err != nil {
		return []domain.GameContext{}, err
	}
	return contexts, nil
}

// Setters and Unsetters for timestamps

func (c *ContextService) SetActivatedAt(ctx context.Context, name domain.GameContextName, activatedAt time.Time) error {
	return c.repo.SetActivatedAt(ctx, name, activatedAt)
}

func (c *ContextService) UnsetActivatedAt(ctx context.Context, name domain.GameContextName) error {
	return c.repo.UnsetActivatedAt(ctx, name)
}

func (c *ContextService) SetDesactivatedAt(ctx context.Context, name domain.GameContextName, desactivatedAt time.Time) error {
	return c.repo.SetDesactivatedAt(ctx, name, desactivatedAt)
}

func (c *ContextService) UnsetDesactivatedAt(ctx context.Context, name domain.GameContextName) error {
	return c.repo.UnsetDesactivatedAt(ctx, name)
}

func (c *ContextService) SetTriggerGeneratedMappingAt(ctx context.Context, name domain.GameContextName, timestamp time.Time) error {
	return c.repo.SetTriggerGeneratedMappingAt(ctx, name, timestamp)
}

func (c *ContextService) UnsetTriggerGeneratedMappingAt(ctx context.Context, name domain.GameContextName) error {
	return c.repo.UnsetTriggerGeneratedMappingAt(ctx, name)
}

func (c *ContextService) SetTriggerImportedTLEAt(ctx context.Context, name domain.GameContextName, timestamp time.Time) error {
	return c.repo.SetTriggerImportedTLEAt(ctx, name, timestamp)
}

func (c *ContextService) UnsetTriggerImportedTLEAt(ctx context.Context, name domain.GameContextName) error {
	return c.repo.UnsetTriggerImportedTLEAt(ctx, name)
}

func (c *ContextService) SetTriggerImportedSatelliteAt(ctx context.Context, name domain.GameContextName, timestamp time.Time) error {
	return c.repo.SetTriggerImportedSatelliteAt(ctx, name, timestamp)
}

func (c *ContextService) UnsetTriggerImportedSatelliteAt(ctx context.Context, name domain.GameContextName) error {
	return c.repo.UnsetTriggerImportedSatelliteAt(ctx, name)
}

package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/clients/redis"
	"github.com/Elbujito/2112/src/app-service/internal/data"
	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xtime"
)

// TleRepository implements the TLERepository interface with caching.
type TleRepository struct {
	db          *data.Database
	redisClient *redis.RedisClient
	cacheTTL    time.Duration
}

// NewTLERepository initializes the repository with a cache TTL.
func NewTLERepository(db *data.Database, redisClient *redis.RedisClient, cacheTTL time.Duration) TleRepository {
	return TleRepository{db: db, redisClient: redisClient, cacheTTL: cacheTTL}
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
		NoradID: domainTLE.NoradID,
		Line1:   domainTLE.Line1,
		Line2:   domainTLE.Line2,
		Epoch:   domainTLE.Epoch,
	}
}

// Cache-Aside Get: Check cache first, fallback to database, and update cache.
func (r *TleRepository) GetTle(ctx context.Context, id string) (domain.TLE, error) {
	key := fmt.Sprintf("satellite:tle:%s", id)

	// Check Redis cache
	data, err := r.redisClient.HGetAll(ctx, key)
	if err == nil && len(data) > 0 {
		epoch, parseErr := xtime.ParseEpoch(data["epoch"])
		if parseErr == nil {
			return domain.TLE{
				ID:    id,
				Line1: data["line_1"],
				Line2: data["line_2"],
				Epoch: epoch,
			}, nil
		}
	}

	// Fallback to database
	var modelTLE models.TLE
	result := r.db.DbHandler.First(&modelTLE, "id = ?", id)
	if result.Error != nil {
		return domain.TLE{}, result.Error
	}

	// Map model to domain
	tle := mapToDomainTLE(modelTLE)

	// Update Redis cache
	cacheData := map[string]interface{}{
		"line_1": tle.Line1,
		"line_2": tle.Line2,
		"epoch":  tle.Epoch,
		"id":     tle.NoradID,
	}
	if err := r.redisClient.HSet(ctx, key, cacheData); err != nil {
		log.Printf("Failed to update Redis cache for key %s: %v\n", key, err)
	}
	err = r.redisClient.Expire(ctx, key, r.cacheTTL)
	if err != nil {
		return tle, err
	}
	return tle, nil
}

// Cache-Aside Save: Save to the database and update the cache.
func (r *TleRepository) SaveTle(ctx context.Context, tle domain.TLE) error {
	// Save to database
	modelTLE := mapToModelTLE(tle)
	if err := r.db.DbHandler.Create(&modelTLE).Error; err != nil {
		return err
	}

	// Update cache
	key := fmt.Sprintf("satellite:tle:%s", tle.ID)
	cacheData := map[string]interface{}{
		"line_1": tle.Line1,
		"line_2": tle.Line2,
		"epoch":  tle.Epoch,
		"id":     tle.NoradID,
	}
	if err := r.redisClient.HSet(ctx, key, cacheData); err != nil {
		log.Printf("Failed to update Redis cache for key %s: %v\n", key, err)
	}
	err := r.redisClient.Expire(ctx, key, r.cacheTTL)
	if err != nil {
		return err
	}
	return nil
}

// Cache-Aside Update: Update the database and refresh the cache.
func (r *TleRepository) UpdateTle(ctx context.Context, tle domain.TLE) error {
	// Update database
	modelTLE := mapToModelTLE(tle)
	if err := r.db.DbHandler.Save(&modelTLE).Error; err != nil {
		return err
	}

	// Refresh cache
	key := fmt.Sprintf("satellite:tle:%s", tle.ID)
	cacheData := map[string]interface{}{
		"line_1": tle.Line1,
		"line_2": tle.Line2,
		"epoch":  tle.Epoch,
		"id":     tle.NoradID,
	}
	if err := r.redisClient.HSet(ctx, key, cacheData); err != nil {
		log.Printf("Failed to update Redis cache for key %s: %v\n", key, err)
	}
	err := r.redisClient.Expire(ctx, key, r.cacheTTL)
	if err != nil {
		return err
	}
	return nil
}

// Cache-Aside Delete: Remove from the database and invalidate the cache.
func (r *TleRepository) DeleteTle(ctx context.Context, id string) error {
	// Delete from database
	if err := r.db.DbHandler.Delete(&models.TLE{}, "id = ?", id).Error; err != nil {
		return err
	}

	key := fmt.Sprintf("satellite:tle:%s", id)
	if err := r.redisClient.Del(ctx, key); err != nil {
		log.Printf("Failed to delete Redis cache for key %s: %v\n", key, err)
	}
	return nil
}

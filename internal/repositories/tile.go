package repository

import (
	"context"
	"errors"

	"github.com/Elbujito/2112/internal/data"
	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/polygon"
	"gorm.io/gorm"
)

type TileRepository struct {
	db *data.Database
}

// NewTileRepository creates a new instance of TileRepository.
func NewTileRepository(db *data.Database) domain.TileRepository {
	return &TileRepository{db: db}
}

// FindByQuadkey retrieves a Tile by its quadkey.
func (r *TileRepository) FindByQuadkey(ctx context.Context, quadkey polygon.Quadkey) (*domain.Tile, error) {
	var tile domain.Tile
	result := r.db.DbHandler.Where("quadkey = ?", quadkey.Key()).First(&tile)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tile, result.Error
}

// FindAll retrieves all Tiles.
func (r *TileRepository) FindAll(ctx context.Context) ([]domain.Tile, error) {
	var tiles []domain.Tile
	result := r.db.DbHandler.Find(&tiles)
	return tiles, result.Error
}

// Save creates a new Tile record.
func (r *TileRepository) Save(ctx context.Context, tile domain.Tile) error {
	return r.db.DbHandler.Create(&tile).Error
}

// Update modifies an existing Tile record.
func (r *TileRepository) Update(ctx context.Context, tile domain.Tile) error {
	return r.db.DbHandler.Save(&tile).Error
}

// DeleteByQuadKey removes a Tile record by its quadkey.
func (r *TileRepository) DeleteByQuadKey(ctx context.Context, key polygon.Quadkey) error {
	return r.db.DbHandler.Where("quadkey = ?", key.Key()).Delete(&domain.Tile{}).Error
}

// Upsert inserts or updates a TLE record in the database.
func (r *TileRepository) Upsert(ctx context.Context, tile domain.Tile) error {
	existingTile, err := r.FindByQuadkey(ctx, tile.Polygon.Center)
	if err != nil {
		return err
	}
	if existingTile != nil {
		return r.Update(ctx, tile)
	}
	return r.Save(ctx, tile)
}

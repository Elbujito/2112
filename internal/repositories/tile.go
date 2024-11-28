package repository

import (
	"context"
	"errors"

	"github.com/Elbujito/2112/internal/data"
	"github.com/Elbujito/2112/internal/data/models"
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
	var tile models.Tile

	result := r.db.DbHandler.Where("quadkey = ?", quadkey.Key()).First(&tile)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	tileMapped := MapToDomain(tile)
	return &tileMapped, result.Error
}

// FindAll retrieves all Tiles.
func (r *TileRepository) FindAll(ctx context.Context) ([]domain.Tile, error) {
	var tiles []models.Tile
	result := r.db.DbHandler.Find(&tiles)
	if result.Error != nil {
		return nil, result.Error
	}

	// Map models to domain
	var domainTiles []domain.Tile
	for _, t := range tiles {
		domainTiles = append(domainTiles, MapToDomain(t))
	}
	return domainTiles, nil
}

// Save creates a new Tile record.
func (r *TileRepository) Save(ctx context.Context, tile domain.Tile) error {
	modelTile := MapToModel(tile)
	return r.db.DbHandler.Create(&modelTile).Error
}

// Update modifies an existing Tile record.
func (r *TileRepository) Update(ctx context.Context, tile domain.Tile) error {
	modelTile := MapToModel(tile)
	return r.db.DbHandler.Save(&modelTile).Error
}

// DeleteByQuadKey removes a Tile record by its quadkey.
func (r *TileRepository) DeleteByQuadKey(ctx context.Context, key polygon.Quadkey) error {
	return r.db.DbHandler.Where("quadkey = ?", key.Key()).Delete(&models.Tile{}).Error
}

// Upsert inserts or updates a Tile record in the database.
func (r *TileRepository) Upsert(ctx context.Context, tile domain.Tile) error {
	existingTile, err := r.FindByQuadkey(ctx, polygon.Quadkey{
		Lat:   tile.CenterLat,
		Long:  tile.CenterLon,
		Level: tile.ZoomLevel,
	})
	if err != nil {
		return err
	}
	if existingTile != nil {
		return r.Update(ctx, tile)
	}
	return r.Save(ctx, tile)
}

// MapToDomain maps a models.Tile to a domain.Tile.
func MapToDomain(modelTile models.Tile) domain.Tile {
	return domain.Tile{
		ID:        modelTile.ID,
		Quadkey:   modelTile.Quadkey,
		ZoomLevel: modelTile.ZoomLevel,
		CenterLat: modelTile.CenterLat,
		CenterLon: modelTile.CenterLon,
	}
}

// MapToModel maps a domain.Tile to a models.Tile.
func MapToModel(domainTile domain.Tile) models.Tile {
	return models.Tile{
		ModelBase: models.ModelBase{
			ID: domainTile.ID,
		},
		Quadkey:   domainTile.Quadkey,
		ZoomLevel: domainTile.ZoomLevel,
		CenterLat: domainTile.CenterLat,
		CenterLon: domainTile.CenterLon,
	}
}

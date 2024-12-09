package services

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

type TileService struct {
	repo domain.TileRepository // Interface for Tile repository
}

// NewTileService creates a new instance of TileService.
func NewTileService(repo domain.TileRepository) TileService {
	return TileService{repo: repo}
}

// FindAllTiles retrieves all tiles from the repository.
func (s *TileService) FindAllTiles(ctx context.Context) ([]domain.Tile, error) {
	tiles, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tiles: %w", err)
	}

	if len(tiles) == 0 {
		return nil, fmt.Errorf("no tiles found")
	}

	return tiles, nil
}

// GetTilesInRegion fetches tiles that intersect with the given bounding box.
func (s *TileService) GetTilesInRegion(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]domain.Tile, error) {
	// Validate input
	if minLat >= maxLat || minLon >= maxLon {
		return nil, fmt.Errorf("invalid bounding box coordinates")
	}

	// Call the repository to fetch tiles
	tiles, err := s.repo.FindTilesInRegion(ctx, minLat, minLon, maxLat, maxLon)
	if err != nil {
		return nil, fmt.Errorf("error fetching tiles in region: %w", err)
	}

	return tiles, nil
}

package services

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

type TileService struct {
	repo        domain.TileRepository
	mappingRepo domain.MappingRepository
}

// NewTileService creates a new instance of TileService.
func NewTileService(repo domain.TileRepository, mappingRepo domain.MappingRepository) TileService {
	return TileService{repo: repo, mappingRepo: mappingRepo}
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

// ListSatellitesMappingWithPagination retrieves satellites with pagination and includes a flag indicating if a TLE is present.
func (s *TileService) ListSatellitesMappingWithPagination(ctx context.Context, page int, pageSize int, search *domain.SearchRequest) ([]domain.TileSatelliteInfo, int64, error) {
	// Validate inputs
	if page <= 0 {
		return nil, 0, fmt.Errorf("page must be greater than 0")
	}
	if pageSize <= 0 {
		return nil, 0, fmt.Errorf("pageSize must be greater than 0")
	}

	// Fetch satellites with pagination and TLE flag
	mappings, totalRecords, err := s.mappingRepo.ListSatellitesMappingWithPagination(ctx, page, pageSize, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve satellites mapping with paginations: %w", err)
	}

	return mappings, totalRecords, nil
}

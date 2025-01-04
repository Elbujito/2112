package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
)

type TileService struct {
	repo          domain.TileRepository
	tleRepo       repository.TleRepository
	satelliteRepo domain.SatelliteRepository
	mappingRepo   domain.MappingRepository
}

// NewTileService creates a new instance of TileService.
func NewTileService(
	tileRepo domain.TileRepository,
	tleRepo repository.TleRepository,
	satelliteRepo domain.SatelliteRepository,
	mappingRepo domain.MappingRepository,
) TileService {
	return TileService{
		repo:          tileRepo,
		tleRepo:       tleRepo,
		satelliteRepo: satelliteRepo,
		mappingRepo:   mappingRepo,
	}
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

func (s *TileService) GetSatelliteMappingsByNoradID(ctx context.Context, noradID string) ([]domain.TileSatelliteInfo, error) {
	mappings, err := s.mappingRepo.GetSatelliteMappingsByNoradID(ctx, noradID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve satellites mapping for [%s]: %w", noradID, err)
	}

	return mappings, nil
}

// RecomputeMappings deletes existing mappings for a NORAD ID and computes new ones.
func (s *TileService) RecomputeMappings(ctx context.Context, noradID string, startTime, endTime time.Time) error {
	log.Printf("Recomputing mappings for NORAD ID: %s\n", noradID)

	// Step 1: Delete existing mappings
	if err := s.mappingRepo.DeleteMappingsByNoradID(ctx, noradID); err != nil {
		return fmt.Errorf("failed to delete existing mappings for NORAD ID [%s]: %w", noradID, err)
	}
	log.Printf("Deleted existing mappings for NORAD ID: %s\n", noradID)

	// Step 2: Fetch satellite data
	satellite, err := s.satelliteRepo.FindByNoradID(ctx, noradID)
	if err != nil {
		return fmt.Errorf("failed to fetch satellite for NORAD ID [%s]: %w", noradID, err)
	}

	// Step 3: Fetch satellite positions
	positions, err := s.tleRepo.QuerySatellitePositions(ctx, satellite.NoradID, startTime, endTime)
	if err != nil {
		return fmt.Errorf("failed to fetch satellite positions for NORAD ID [%s]: %w", noradID, err)
	}

	// Ensure there are enough positions to compute mappings
	if len(positions) < 2 {
		log.Printf("Not enough positions to compute mappings for NORAD ID: %s\n", noradID)
		return nil
	}

	// Step 4: Compute new mappings
	mappings, err := s.repo.FindTilesVisibleFromLine(ctx, satellite, positions)
	if err != nil {
		return fmt.Errorf("failed to compute tile mappings for NORAD ID [%s]: %w", noradID, err)
	}

	// Step 5: Save new mappings
	if err := s.mappingRepo.SaveBatch(ctx, mappings); err != nil {
		return fmt.Errorf("failed to save new mappings for NORAD ID [%s]: %w", noradID, err)
	}

	log.Printf("Recomputed and saved %d mappings for NORAD ID: %s\n", len(mappings), noradID)
	return nil
}

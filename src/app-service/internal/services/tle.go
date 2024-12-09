package services

import (
	"context"

	"github.com/Elbujito/2112/src/app-service/internal/api/mappers"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

type celestrackClient interface {
	FetchTLEFromSatCatByCategory(ctx context.Context, category string) ([]*mappers.RawTLE, error)
	FetchSatelliteMetadata(ctx context.Context) ([]*mappers.SatelliteMetadata, error)
}

type TleService struct {
	celestrackClient celestrackClient
}

// NewTleService creates a new instance of TleService.
func NewTleService(celestrackClient celestrackClient) TleService {
	return TleService{celestrackClient: celestrackClient}
}

func (s *TleService) FetchTLEFromSatCatByCategory(ctx context.Context, category string) ([]domain.TLE, error) {
	rawTLEs, err := s.celestrackClient.FetchTLEFromSatCatByCategory(ctx, category)
	if err != nil {
		return []domain.TLE{}, err
	}

	tles := make([]domain.TLE, len(rawTLEs))
	for idx, raw := range rawTLEs {
		tle, err := domain.NewTLE(
			raw.NoradID,
			raw.Line1,
			raw.Line2,
		)

		if err != nil {
			return []domain.TLE{}, err
		}
		tles[idx] = tle
	}

	return tles, err
}

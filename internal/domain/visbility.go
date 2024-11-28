package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type VisibilityRepository interface {
	FindByNoradIDAndTile(ctx context.Context, noradID string, tileID string) ([]Visibility, error)
	FindAll(ctx context.Context) ([]Visibility, error)
	Save(ctx context.Context, visibility Visibility) error
	Update(ctx context.Context, visibility Visibility) error
	Delete(ctx context.Context, id string) error
}

// Visibility represents the domain entity Visibility
type Visibility struct {
	ID           string // Unique identifier
	NoradID      string
	TileID       string
	StartTime    time.Time
	EndTime      time.Time
	MaxElevation float64
}

// NewVisibility constructor
func NewVisibility(noradID string,
	tileID string,
	startTime time.Time,
	endTime time.Time,
	maxElevation float64) Visibility {
	return Visibility{
		ID:           uuid.NewString(),
		NoradID:      noradID,
		TileID:       tileID,
		StartTime:    startTime,
		EndTime:      endTime,
		MaxElevation: maxElevation,
	}

}

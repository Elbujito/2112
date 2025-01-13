package models

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

// TLE Model
type TLE struct {
	ModelBase
	NoradID string    `gorm:"size:255;not null;index"` // Foreign key to Satellite table via Norad ID
	Line1   string    `gorm:"size:255;not null"`
	Line2   string    `gorm:"size:255;not null"`
	Epoch   time.Time `gorm:"not null"` // Time associated with the TLE
}

// MapToDomain converts a models.Tile to a domain.Tile.
func MapToTLEDomain(t TLE) domain.TLE {
	// Convert to domain
	return domain.TLE{
		ModelBase: domain.ModelBase{
			ID:          t.ID,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   &t.UpdatedAt,
			DeleteAt:    t.DeleteAt,
			ProcessedAt: t.ProcessedAt,
			IsActive:    t.IsActive,
			IsFavourite: t.IsFavourite,
			DisplayName: t.DisplayName,
		},
		NoradID: t.NoradID,
		Line1:   t.Line1,
		Line2:   t.Line2,
		Epoch:   t.Epoch,
	}
}

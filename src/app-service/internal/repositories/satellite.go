package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Elbujito/2112/internal/data"
	"github.com/Elbujito/2112/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SatelliteWithTLE struct {
	ID         string     `gorm:"column:id"`
	Name       string     `gorm:"column:name"`
	NoradID    string     `gorm:"column:norad_id"`
	Owner      string     `gorm:"column:owner"`
	LaunchDate *time.Time `gorm:"column:launch_date"`
	Apogee     *float64   `gorm:"column:apogee"`
	Perigee    *float64   `gorm:"column:perigee"`
	Line1      *string    `gorm:"column:line1"`
	Line2      *string    `gorm:"column:line2"`
	UpdatedAt  *time.Time `gorm:"column:updated_at"`
}

type SatelliteRepository struct {
	db *data.Database
}

// NewSatelliteRepository creates a new instance of SatelliteRepository.
func NewSatelliteRepository(db *data.Database) domain.SatelliteRepository {
	return &SatelliteRepository{db: db}
}

// FindByNoradID retrieves a satellite by its NORAD ID.
func (r *SatelliteRepository) FindByNoradID(ctx context.Context, noradID string) (domain.Satellite, error) {
	var satellite domain.Satellite
	result := r.db.DbHandler.Where("norad_id = ?", noradID).First(&satellite)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.Satellite{}, nil
	}
	return satellite, result.Error
}

// FindAll retrieves all satellites.
func (r *SatelliteRepository) FindAll(ctx context.Context) ([]domain.Satellite, error) {
	var satellites []domain.Satellite
	result := r.db.DbHandler.Find(&satellites)
	return satellites, result.Error
}

// Save creates a new satellite record.
func (r *SatelliteRepository) Save(ctx context.Context, satellite domain.Satellite) error {
	return r.db.DbHandler.Create(&satellite).Error
}

// Update modifies an existing satellite record.
func (r *SatelliteRepository) Update(ctx context.Context, satellite domain.Satellite) error {
	return r.db.DbHandler.Save(&satellite).Error
}

// DeleteByNoradID removes a satellite record by its NoradID.
func (r *SatelliteRepository) DeleteByNoradID(ctx context.Context, noradID string) error {
	return r.db.DbHandler.Where("noradID = ?", noradID).Delete(&domain.Satellite{}).Error
}

// SaveBatch performs a batch upsert (insert or update) for satellites.
func (r *SatelliteRepository) SaveBatch(ctx context.Context, satellites []domain.Satellite) error {
	if len(satellites) == 0 {
		return nil // Nothing to save
	}

	// Use Gen's support for ON CONFLICT upsert
	return r.db.DbHandler.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "norad_id"}}, // Define the unique constraint column
			UpdateAll: true,                                // Update all fields in case of conflict
		}).
		CreateInBatches(satellites, 100).Error // Batch size: 100
}
func (r *SatelliteRepository) FindAllWithPagination(ctx context.Context, page int, pageSize int, searchRequest *domain.SearchRequest) ([]domain.Satellite, int64, error) {
	var results []SatelliteWithTLE
	var totalRecords int64

	// Calculate the offset
	offset := (page - 1) * pageSize

	// Build the base query with LEFT JOIN to include optional TLE data
	query := r.db.DbHandler.Table("satellites").
		Select(`
			satellites.id, satellites.name, satellites.norad_id, satellites.owner, 
			satellites.launch_date, satellites.apogee, satellites.perigee, 
			tles.line1, tles.line2, tles.updated_at`).
		Joins("LEFT JOIN tles ON satellites.norad_id = tles.norad_id")

	// Apply case-insensitive wildcard filter if provided
	if searchRequest != nil && searchRequest.Wildcard != "" {
		wildcard := "%" + searchRequest.Wildcard + "%"
		query = query.Where(
			"LOWER(satellites.norad_id) LIKE LOWER(?) OR LOWER(satellites.name) LIKE LOWER(?)",
			wildcard, wildcard,
		)
	}

	// Count the total number of records
	countQuery := query.Session(&gorm.Session{}) // Clone the query to avoid side effects
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve satellites with pagination
	paginationQuery := query.Limit(pageSize).Offset(offset)
	if err := paginationQuery.Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	// Map the results to the domain.Satellite structure
	var satellites []domain.Satellite
	for _, result := range results {
		satellite := domain.Satellite{
			ID:         result.ID,
			Name:       result.Name,
			NoradID:    result.NoradID,
			Owner:      result.Owner,
			LaunchDate: result.LaunchDate,
			Apogee:     result.Apogee,
			Perigee:    result.Perigee,
		}

		// Add TLE data if available
		if result.Line1 != nil && result.Line2 != nil && result.UpdatedAt != nil {
			satellite.TleUpdatedAt = result.UpdatedAt
		}

		satellites = append(satellites, satellite)
	}

	return satellites, totalRecords, nil
}

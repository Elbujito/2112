package repository

import (
	"context"

	"github.com/Elbujito/2112/src/templates/go-server/internal/data"
	"github.com/Elbujito/2112/src/templates/go-server/internal/domain"
	"gorm.io/gorm/clause"
)

type TestRepository struct {
	db *data.Database
}

// NewTestRepository creates a new instance of TestRepository.
func NewTestRepository(db *data.Database) domain.TestRepository {
	return &TestRepository{db: db}
}

// FindAll retrieves all tests.
func (r *TestRepository) FindAll(ctx context.Context) ([]domain.Test, error) {
	var tests []domain.Test
	result := r.db.DbHandler.Find(&tests)
	return tests, result.Error
}

// Save creates a new test record.
func (r *TestRepository) Save(ctx context.Context, test domain.Test) error {
	return r.db.DbHandler.Create(&test).Error
}

// Update modifies an existing test record.
func (r *TestRepository) Update(ctx context.Context, test domain.Test) error {
	return r.db.DbHandler.Save(&test).Error
}

// SaveBatch performs a batch upsert (insert or update) for tests.
func (r *TestRepository) SaveBatch(ctx context.Context, tests []domain.Test) error {
	if len(tests) == 0 {
		return nil // Nothing to save
	}

	// Use Gen's support for ON CONFLICT upsert
	return r.db.DbHandler.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "your_clause"}}, // Define the unique constraint column
			UpdateAll: true,                                   // Update all fields in case of conflict
		}).
		CreateInBatches(tests, 100).Error // Batch size: 100
}

package migrations

import (
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migrations definition
var Migrations *gormigrate.Gormigrate

// MigrationsList definition
var MigrationsList = []*gormigrate.Migration{}

// Init init
func Init(db *gorm.DB) {
	Migrations = gormigrate.New(db, gormigrate.DefaultOptions, MigrationsList)
}

// AddMigration definition
func AddMigration(migration *gormigrate.Migration) {
	MigrationsList = append(MigrationsList, migration)
}

// Migrate definition
func Migrate() error {
	return Migrations.Migrate()
}

// Rollback definition
func Rollback() error {
	return Migrations.RollbackLast()
}

// AutoMigrateAndLog definition
func AutoMigrateAndLog(db *gorm.DB, model interface{}, id string) error {
	if err := db.AutoMigrate(model); err != nil {
		logFail(id, err)
		return err
	}
	logSuccess(id)
	return nil
}

func logSuccess(id string, rollback ...bool) {
	if len(rollback) > 0 && rollback[0] {
		logger.Info("Rolled back migration: %s", id)
		return
	}
	logger.Info("Applied migration: %s", id)
}

func logFail(id string, err error, rollback ...bool) {
	if len(rollback) > 0 && rollback[0] {
		logger.Error("Failed to rollback migration: %s, error: %s", id, err)
		return
	}
	logger.Error("Failed to apply migration: %s, error: %s", id, err)
}

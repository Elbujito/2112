package migrations

import (
	"github.com/Elbujito/2112/template/go-server/internal/data/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "2024112604_create",
		Migrate: func(db *gorm.DB) error {
			type Test struct {
				models.ModelBase
				Name string `gorm:"size:255;not null"` // Satellite name
			}

			// AutoMigrate all tables
			if err := db.AutoMigrate(&Test{}); err != nil {
				return err
			}

			return nil
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable("test")
		},
	}

	AddMigration(m)
}

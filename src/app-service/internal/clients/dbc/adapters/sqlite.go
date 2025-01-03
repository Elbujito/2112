package adapters

import (
	"fmt"
	"reflect"

	"github.com/Elbujito/2112/src/app-service/internal/config/features"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	Adapters.AppendAdapter(xconstants.DB_PLATFORM_SQLITE, &SQLiteAdapter{
		requiredConfig: []string{"Name"},
	})
}

type SQLiteAdapter struct {
	IAdapter
	config         features.DatabaseConfig
	requiredConfig []string
}

func (a *SQLiteAdapter) SetConfig(config features.DatabaseConfig) {
	a.config = config
}

func (a *SQLiteAdapter) GetDriver() (gorm.Dialector, error) {
	dsn, _ := a.GetDSN()
	return sqlite.Open(dsn), nil
}

func (a *SQLiteAdapter) GetServerDriver() (gorm.Dialector, error) {
	dsn, _ := a.GetServerDSN()
	return sqlite.Open(dsn), nil
}

func (a *SQLiteAdapter) GetDSN() (string, error) {
	return a.config.Name, nil
}

func (a *SQLiteAdapter) GetServerDSN() (string, error) {
	return a.config.Name, nil
}

func (a *SQLiteAdapter) GetDbCreateStatement() (string, error) {
	return "", nil
}

func (a *SQLiteAdapter) GetDbDropStatement() (string, error) {
	return "", nil
}

func (a *SQLiteAdapter) ValidateConfig() error {
	vOfConfig := reflect.ValueOf(a.config)
	if vOfConfig.Kind() == reflect.Ptr {
		vOfConfig = vOfConfig.Elem()
	}
	for _, requiredField := range a.requiredConfig {
		for i := 0; i < vOfConfig.NumField(); i++ {
			configField := vOfConfig.Field(i)
			if configField.Kind() == reflect.String {
				if vOfConfig.Type().Field(i).Name == requiredField &&
					vOfConfig.Field(i).Interface().(string) == "" {
					return fmt.Errorf("database adapter requirements not satisfied. missing required field: %s", vOfConfig.Type().Field(i).Name)
				}
			}
		}
	}
	return nil
}

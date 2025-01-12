package adapters

import (
	"fmt"
	"reflect"

	"github.com/Elbujito/2112/src/app-service/internal/config/features"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	Adapters.AppendAdapter(xconstants.DB_PLATFORM_MYSQL, &MySQLAdapter{
		requiredConfig: []string{"Host", "Port", "User", "Password", "Name"},
	})
}

type MySQLAdapter struct {
	IAdapter
	config         features.DatabaseConfig
	requiredConfig []string
}

func (a *MySQLAdapter) SetConfig(config features.DatabaseConfig) {
	a.config = config
}

func (a *MySQLAdapter) GetDriver() (gorm.Dialector, error) {
	dsn, _ := a.GetDSN()
	return mysql.Open(dsn), nil
}

func (a *MySQLAdapter) GetServerDriver() (gorm.Dialector, error) {
	dsn, _ := a.GetServerDSN()
	return mysql.Open(dsn), nil
}

func (a *MySQLAdapter) GetDSN() (string, error) {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		a.config.User,
		a.config.Password,
		a.config.Host,
		a.config.Port,
		a.config.Name), nil
}

func (a *MySQLAdapter) GetServerDSN() (string, error) {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		a.config.User,
		a.config.Password,
		a.config.Host,
		a.config.Port), nil
}

func (a *MySQLAdapter) GetDbCreateStatement() (string, error) {
	return "CREATE DATABASE IF NOT EXISTS ", nil
}

func (a *MySQLAdapter) GetDbDropStatement() (string, error) {
	return "DROP DATABASE IF EXISTS ", nil
}

func (a *MySQLAdapter) ValidateConfig() error {
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

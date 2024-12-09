package adapters

import (
	"github.com/Elbujito/2112/src/app-service/internal/config/features"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"gorm.io/gorm"
)

var Adapters = &Adapter{
	defaultPlatform: xconstants.DEFAULT_DB_PLATFORM,
	currentPlatform: xconstants.DEFAULT_DB_PLATFORM,
	adapters:        make(map[string]IAdapter),
}

type IAdapter interface {
	SetConfig(config features.DatabaseConfig)
	GetDriver() (gorm.Dialector, error)
	GetServerDriver() (gorm.Dialector, error)
	GetDSN() (string, error)
	GetServerDSN() (string, error)
	GetDbCreateStatement() (string, error)
	GetDbDropStatement() (string, error)
	ValidateConfig() error
}

type Adapter struct {
	IAdapter
	adapters        map[string]IAdapter
	defaultPlatform string
	currentPlatform string
	config          features.DatabaseConfig
}

func (a *Adapter) SetConfig(config features.DatabaseConfig) {
	a.config = config
	a.currentPlatform = config.Platform
	for _, adapter := range a.adapters {
		adapter.SetConfig(a.config)
	}
}

func (a *Adapter) GetDriver() (gorm.Dialector, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDriver()
	}
	return nil, xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

func (a *Adapter) GetServerDriver() (gorm.Dialector, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDriver()
	}
	return nil, xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

func (a *Adapter) GetDSN() (string, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDSN()
	}
	return "", xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

func (a *Adapter) GetServerDSN() (string, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetServerDSN()
	}
	return "", xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

func (a *Adapter) AppendAdapter(name string, adapter IAdapter) {
	a.adapters[name] = adapter
}

func (a *Adapter) GetDbCreateStatement() (string, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDbCreateStatement()
	}
	return "", xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

func (a *Adapter) GetDbDropStatement() (string, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDbDropStatement()
	}
	return "", xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

func (a *Adapter) ValidateConfig() error {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.ValidateConfig()
	}
	return xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

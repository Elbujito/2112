package proc

import (
	"time"

	"github.com/Elbujito/2112/internal/clients/logger"
	"github.com/Elbujito/2112/internal/clients/service"
)

func StartWatcher() {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	interval := fx.IntFromStr(config.WatcherSleepInterval)

	go func() {
		// This is a sample watcher
		// Command execution goes here ...

		logger.Info("Watcher started")
		for {
			// Watcher logic goes here ...
			logger.Info("Watcher running...")
			// Break the loop after 10 iterations
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}

	}()
}

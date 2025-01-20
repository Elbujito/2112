package proc

import (
	"io"
	"os"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/clients/service"
	log "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xutils"
)

func StartWatcher() {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	interval := xutils.IntFromStr(config.WatcherSleepInterval)

	var logWriter io.Writer
	logWriter = os.Stdout
	logger, err := log.NewLogger(logWriter, log.DebugLevel, log.LoggerTypes.Logrus())
	if err != nil {
		panic(err)
	}
	log.SetDefaultLogger(logger)
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

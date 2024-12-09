package proc

import (
	"context"
	"fmt"
	"log"

	"github.com/Elbujito/2112/src/app-service/internal/clients/celestrack"
	propagator "github.com/Elbujito/2112/src/app-service/internal/clients/propagate"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/data"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	"github.com/Elbujito/2112/src/app-service/internal/tasks"
	"github.com/Elbujito/2112/src/app-service/internal/tasks/handlers"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xutils"
)

func TaskExec(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a task name")
		return
	}
	taskName := args[0]
	taskArgs := xutils.ResolveArgs(args[1:])

	database := data.NewDatabase()

	propagteClient := propagator.NewPropagatorClient(config.Env)

	// Assuming you have a service or repository to fetch tiles by NORAD ID
	tleRepo := repository.NewTLERepository(&database)
	celestrackClient := celestrack.NewCelestrackClient(config.Env)
	satelliteRepo := repository.NewSatelliteRepository(&database)
	visibilityRepo := repository.NewTileSatelliteMappingRepository(&database)
	tileRepo := repository.NewTileRepository(&database)

	tleService := services.NewTleService(celestrackClient)
	satService := services.NewSatelliteService(tleRepo, propagteClient, celestrackClient, satelliteRepo)

	monitor, err := tasks.NewTaskMonitor(satelliteRepo, tleRepo, tileRepo, visibilityRepo, tleService, satService)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = monitor.Process(ctx, handlers.TaskName(taskName), taskArgs)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

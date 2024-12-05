package proc

import (
	"context"
	"fmt"
	"log"

	"github.com/Elbujito/2112/internal/clients/celestrack"
	propagator "github.com/Elbujito/2112/internal/clients/propagate"
	"github.com/Elbujito/2112/internal/config"
	"github.com/Elbujito/2112/internal/data"
	repository "github.com/Elbujito/2112/internal/repositories"
	"github.com/Elbujito/2112/internal/services"
	"github.com/Elbujito/2112/internal/tasks"
	"github.com/Elbujito/2112/pkg/fx"
)

func TaskExec(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a task name")
		return
	}
	taskName := args[0]
	taskArgs := fx.ResolveArgs(args[1:])

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

	err = monitor.Process(ctx, tasks.TaskName(taskName), taskArgs)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

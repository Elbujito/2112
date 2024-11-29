package proc

import (
	"context"
	"fmt"
	"log"

	"github.com/Elbujito/2112/internal/clients/celestrack"
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

	celestrackClient := celestrack.CelestrackClient{}

	satelliteRepo := repository.NewSatelliteRepository(&database)
	tleRepo := repository.NewTLERepository(&database)
	visibilityRepo := repository.NewTileSatelliteMappingRepository(&database)
	tileRepo := repository.NewTileRepository(&database)

	tleService := services.NewTleService(&celestrackClient)

	monitor, err := tasks.NewTaskMonitor(satelliteRepo, tleRepo, tileRepo, visibilityRepo, tleService)
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

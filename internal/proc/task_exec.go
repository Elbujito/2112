package proc

import (
	"context"
	"fmt"
	"log"

	"github.com/Elbujito/2112/internal/data"
	"github.com/Elbujito/2112/internal/repository"
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

	satelliteRepo := repository.NewSatelliteRepository(&database)
	tleRepo := repository.NewTLERepository(database.DbHandler)
	// tleService := services.NewTleService(celestrack.FetchCategoryTLEHandler)

	monitor, err := tasks.NewTaskMonitor(satelliteRepo, tleRepo, services.TleService{})
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

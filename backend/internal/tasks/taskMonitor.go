package tasks

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/internal/services"
	"github.com/Elbujito/2112/internal/tasks/handlers"
)

type TaskName string

type TaskHandler interface {
	GetTask() handlers.Task
	Run(ctx context.Context, args map[string]string) error
}

type TaskMonitor struct {
	Tasks map[TaskName]TaskHandler
}

func NewTaskMonitor(satelliteRepo domain.SatelliteRepository, tleRepo domain.TLERepository, tileRepo domain.TileRepository, visibilityRepo domain.TileSatelliteMappingRepository, tleService services.TleService, satelliteService services.SatelliteService) (TaskMonitor, error) {

	tleProvisionHandler := handlers.NewTLEProvisionHandler(
		satelliteRepo,
		tleRepo,
		&tleService,
	)

	tileProvisionHandler := handlers.NewTileProvisionHandler(
		tileRepo,
	)

	computeVisibilitiesHandler := handlers.NewComputeVisibilitiesHandler(
		tileRepo,
		tleRepo,
		satelliteRepo,
		visibilityRepo,
	)

	visbilityBySatelliteHorizonHandle := handlers.NewVisbilityBySatelliteHorizonHandler(
		tileRepo,
		tleRepo,
		satelliteRepo,
		visibilityRepo,
	)

	satelliteProvisionHandler := handlers.NewSatelliteProvisionHandler(
		satelliteRepo,
		&satelliteService,
	)

	tasks := map[TaskName]TaskHandler{
		TaskName("fetchAndUpsertTLE"):                   &tleProvisionHandler,
		TaskName("fetchAndStoreTiles"):                  &tileProvisionHandler,
		TaskName("execComputeVisibilitiesTask"):         &computeVisibilitiesHandler,
		TaskName("execVisbilityBySatelliteHorizonTask"): &visbilityBySatelliteHorizonHandle,
		TaskName("fetchAndUpsertSatellite"):             &satelliteProvisionHandler,
	}
	return TaskMonitor{
		Tasks: tasks,
	}, nil
}

func (t *TaskMonitor) Process(ctx context.Context, taskName TaskName, args map[string]string) error {
	handler, err := t.GetMatchingTask(taskName)
	if err != nil {
		return err
	}
	return handler.Run(ctx, args)
}

func (t *TaskMonitor) GetMatchingTask(taskName TaskName) (task TaskHandler, err error) {
	hh, ok := t.Tasks[taskName]
	if !ok {
		return task, fmt.Errorf("task no found for [%s]", taskName)
	}
	return hh, nil
}

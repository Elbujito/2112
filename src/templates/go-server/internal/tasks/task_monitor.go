package tasks

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/src/template/go-server/internal/domain"
	"github.com/Elbujito/2112/src/template/go-server/internal/tasks/handlers"
)

type TaskHandler interface {
	GetTask() handlers.Task
	Run(ctx context.Context, args map[string]string) error
}

type TaskMonitor struct {
	Tasks map[handlers.TaskName]TaskHandler
}

func NewTaskMonitor(testRepo domain.TestRepository) (TaskMonitor, error) {

	testHandler := handlers.NewTestHandler(
		testRepo,
	)

	tasks := map[handlers.TaskName]TaskHandler{
		testHandler.GetTask().Name: &testHandler,
	}
	return TaskMonitor{
		Tasks: tasks,
	}, nil
}

func (t *TaskMonitor) Process(ctx context.Context, taskName handlers.TaskName, args map[string]string) error {
	handler, err := t.GetMatchingTask(taskName)
	if err != nil {
		return err
	}
	return handler.Run(ctx, args)
}

func (t *TaskMonitor) GetMatchingTask(taskName handlers.TaskName) (task TaskHandler, err error) {
	hh, ok := t.Tasks[taskName]
	if !ok {
		return task, fmt.Errorf("task no found for [%s]", taskName)
	}
	return hh, nil
}

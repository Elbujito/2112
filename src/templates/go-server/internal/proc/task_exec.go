package proc

import (
	"context"
	"fmt"
	"log"

	"github.com/Elbujito/2112/src/template/go-server/internal/data"
	repository "github.com/Elbujito/2112/src/template/go-server/internal/repositories"
	"github.com/Elbujito/2112/src/template/go-server/internal/tasks"
	"github.com/Elbujito/2112/src/template/go-server/internal/tasks/handlers"
	"github.com/Elbujito/2112/src/template/go-server/pkg/fx/xutils"
)

func TaskExec(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a task name")
		return
	}
	taskName := args[0]
	taskArgs := xutils.ResolveArgs(args[1:])

	database := data.NewDatabase()

	testRepo := repository.NewTestRepository(&database)

	monitor, err := tasks.NewTaskMonitor(testRepo)
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

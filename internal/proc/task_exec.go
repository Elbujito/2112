package proc

import (
	"fmt"

	"github.com/Elbujito/2112/internal/tasks"
	"github.com/Elbujito/2112/pkg/fx"
)

func TaskExec(args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a task name")
		return
	}
	taskName := args[0]
	task := tasks.Tasks.GetTask(taskName)
	if task == nil {
		fmt.Println("Task not found")
		return
	}
	taskArgs := fx.ResolveArgs(args[1:])
	if err := task.Execute(taskArgs); err != nil {
		fmt.Println(err)
	}
}

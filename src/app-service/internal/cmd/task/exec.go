package task

import (
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/proc"
	"github.com/spf13/cobra"
)

// ExecCmd creates the `exec` subcommand
func ExecCmd(serviceComponent interface{}) *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Start exec task",
		Long:  `Start the exec task.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Executing a task...")
			proc.TaskExec(cmd.Context(), args)
		},
	}
}

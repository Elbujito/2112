package task

import (
	"github.com/Elbujito/2112/src/templates/go-server/internal/proc"

	"github.com/spf13/cobra"
)

// ExecCmd represents the exec command
var ExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Start exec task",
	Long:  `Start the exec task.`,
	Run:   execExecCmd,
}

func init() {
	// This is auto executed upon start
	// Initialization processes can go here ...
}

func execExecCmd(cmd *cobra.Command, args []string) {
	proc.TaskExec(cmd.Context(), args)
}

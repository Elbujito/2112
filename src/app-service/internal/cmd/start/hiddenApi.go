package start

import (
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/proc"

	"github.com/spf13/cobra"
)

// HiddenApiCmd represents the hiddenApi command
var HiddenApiCmd = &cobra.Command{
	Use:   "hiddenApi",
	Short: "Start hidden API service",
	Long:  `Start hidden API web server.`,
	Run:   execHiddenApiCmd,
}

func init() {
	// This is auto executed upon start
	// Initialization processes can go here ...
}

func execHiddenApiCmd(cmd *cobra.Command, args []string) {
	// Command execution goes here ...
	if config.StartWatcherFlag {
		go WatcherCmd.Run(cmd, args)
	}
	proc.StartHiddenApi()
}

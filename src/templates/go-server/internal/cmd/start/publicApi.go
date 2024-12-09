package start

import (
	"github.com/Elbujito/2112/src/template/go-server/internal/config"
	"github.com/Elbujito/2112/src/template/go-server/internal/proc"

	"github.com/spf13/cobra"
)

// PublicApiCmd represents the publicApi command
var PublicApiCmd = &cobra.Command{
	Use:   "publicApi",
	Short: "Start public API service",
	Long:  `Start public API web server.`,
	Run:   execPublicApiCmd,
}

func init() {
	// This is auto executed upon start
	// Initialization processes can go here ...
}

func execPublicApiCmd(cmd *cobra.Command, args []string) {
	// Command execution goes here ...
	if config.StartWatcherFlag {
		go WatcherCmd.Run(cmd, args)
	}
	proc.StartPublicApi()
}

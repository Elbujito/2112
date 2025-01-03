package db

import (
	"github.com/Elbujito/2112/src/templates/go-server/internal/proc"

	"github.com/spf13/cobra"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create database",
	Long:  `Create dababase.`,
	Run:   execCreateCmd,
}

func init() {
	// This is auto executed upon start
	// Initialization processes can go here ...
}

func execCreateCmd(cmd *cobra.Command, args []string) {
	// Command execution goes here ...
	proc.DBCreate()
}

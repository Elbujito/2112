package db

import (
	"github.com/Elbujito/2112/src/template/go-server/internal/proc"

	"github.com/spf13/cobra"
)

// MigrateCmd represents the migrate command
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Long:  `Run database migrations.`,
	Run:   execMigrateCmd,
}

func init() {
	// This is auto executed upon start
	// Initialization processes can go here ...
}

func execMigrateCmd(cmd *cobra.Command, args []string) {
	// Command execution goes here ...
	proc.DBMigrate()
}

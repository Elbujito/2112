package cmd

import (
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/cmd/db"
	"github.com/Elbujito/2112/src/app-service/internal/proc"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db <option>",
	Short: "Start db related operations",
	Long: `Start a database operation.
Please key in an option to start. Type 'db -h' for more information.

Popular options are:
- db migrate
- db rollback`,
}

func init() {
	// This is auto executed upon start
	// Initialization processes can go here ...

	// Register sub commands
	dbCmd.AddCommand(db.CreateCmd)
	dbCmd.AddCommand(db.DropCmd)
	dbCmd.AddCommand(db.MigrateCmd)
	dbCmd.AddCommand(db.RollbackCmd)
	dbCmd.AddCommand(db.SeedCmd)

	// Register persistent function for all sub commands
	dbCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(cmd, args)
		execDbPersistentPreRun(cmd)
	}

	// Register db command
	rootCmd.AddCommand(dbCmd)
}

func execDbPersistentPreRun(cmd *cobra.Command) {
	logger.Debug("Executing db persistent pre run ...")

	proc.InitDbClient()
	proc.ConfigureClients()

	ca := cmd.CalledAs()

	// here we don't always establish the connection because different
	// commands may require different database connections
	// so we only establish the connection when calling migrate, rollback or seed
	switch ca {
	case xconstants.NAME_CMD_DB_MIGRATE,
		xconstants.NAME_CMD_DB_ROLLBACK,
		xconstants.NAME_CMD_DB_SEED:
		proc.InitDbConnection()
		proc.InitModels()
	}

	// You can initialize other features here ...
	// this will run before any command, make sure to put only global initializations here
	// to avoid running into nil pointers or undefined variables
	// ...

}

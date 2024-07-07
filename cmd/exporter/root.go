package exporter

import (
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/exporter"
	"github.com/castus/speedcube-events/logger"
	"github.com/spf13/cobra"
)

var log = logger.Default()
var exportLocally bool
var exportForFrontend bool
var persistDatabase bool

func Setup() *cobra.Command {
	cmd.Flags().BoolVarP(&exportLocally, "local", "l", false, "Export database locally")
	cmd.Flags().BoolVarP(&exportForFrontend, "frontend", "f", false, "Export database to S3 for frontend use")
	cmd.Flags().BoolVarP(&persistDatabase, "database", "d", false, "Export database to DynamoDB")

	return cmd
}

var cmd = &cobra.Command{
	Use:   "export",
	Short: "Export DynamoDB database as JSON to S3 or to a file",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.Database{}
		database.Initialize()
		if persistDatabase {
			exporter.PersistDatabase(database)
		} else if exportForFrontend {
			exporter.ExportForFrontend(database)
		} else {
			exporter.ExportLocal(database)
		}
	},
}

package exporter

import (
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/exporter"
	"github.com/castus/speedcube-events/logger"
	"github.com/spf13/cobra"
)

var log = logger.Default()
var exportLocally bool

func Setup() *cobra.Command {
	cmd.Flags().BoolVarP(&exportLocally, "local", "l", false, "Export database locally")

	return cmd
}

var cmd = &cobra.Command{
	Use:   "export",
	Short: "Export DynamoDB database as JSON to S3 or to a file",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.Database{}
		database.Initialize()
		dbCompetitions := database.GetAll()
		if exportLocally {
			exporter.ExportLocal(dbCompetitions)
		} else {
			exporter.Export(dbCompetitions)
		}
	},
}

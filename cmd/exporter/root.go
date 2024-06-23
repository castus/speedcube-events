package exporter

import (
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/exporter"
	"github.com/castus/speedcube-events/logger"
	"github.com/spf13/cobra"
)

var log = logger.Default()

var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export DynamoDB database as JSON to S3",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := db.GetClient()
		if err != nil {
			log.Error("Couldn't get database client", err)
			panic(err)
		}

		dbCompetitions, err := db.AllItems(c)
		if err != nil {
			log.Error("Couldn't fetch items from database", err)
			panic(err)
		}
		exporter.Export(dbCompetitions)
	},
}

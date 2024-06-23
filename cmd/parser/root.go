package parser

import (
	"github.com/castus/speedcube-events/externalParser"
	"github.com/castus/speedcube-events/logger"
	"github.com/spf13/cobra"
)

var log = logger.Default()

var Cmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse external data files from S3 and update database",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Running external parsing...")
		externalParser.Run()
	},
}

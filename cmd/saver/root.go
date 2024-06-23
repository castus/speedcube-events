package saver

import (
	"github.com/castus/speedcube-events/exporter"
	"github.com/castus/speedcube-events/logger"
	"github.com/spf13/cobra"
)

var log = logger.Default()

var Cmd = &cobra.Command{
	Use:   "save",
	Short: "Save a source of truth page to disc",
	Run: func(cmd *cobra.Command, args []string) {
		exporter.SaveWebpageAsFile("kalendarz-imprez.html")
		log.Info("Webpage saved on disk")
	},
}

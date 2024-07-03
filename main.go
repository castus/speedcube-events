package main

import (
	"fmt"
	"os"

	"github.com/castus/speedcube-events/cmd/exporter"
	"github.com/castus/speedcube-events/cmd/parser"
	"github.com/castus/speedcube-events/cmd/saver"
	"github.com/castus/speedcube-events/cmd/scraper"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "speedcube",
		Short: "speedcube - CLI to manage Speedcube events project",
		Long:  `The Speedcube events project's purpose is to create a place when all Poland Speedcube Events can be displayed with details like: the main competition, all competitions, limit of competitors, currently registered competitors etc.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Type help to see all available commands")
		},
	}

	rootCmd.AddCommand(saver.Cmd)
	rootCmd.AddCommand(parser.Cmd)
	rootCmd.AddCommand(exporter.Setup())
	rootCmd.AddCommand(scraper.Setup())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	"os"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "timex",
	Short: "timex is a simple time tracking CLI",
	Long:  `A fast and simple command-line tool to help you track time spent on your projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action is to show help
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(database.InitDatabase)
}

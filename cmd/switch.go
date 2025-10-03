package cmd

import (
	"fmt"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch [project_name]",
	Short: "Switch the current timer to a new project",
	Long:  `Stop the current timer and start a new one for the given project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newProject := args[0]
		oldProject, err := database.SwitchTimer(newProject)
		if err != nil {
			fmt.Printf("Error switching timer: %v\n", err)
			return
		}
		fmt.Printf("Switched from '%s' to '%s'\n", *oldProject, newProject)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}

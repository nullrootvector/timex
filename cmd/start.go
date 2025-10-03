package cmd

import (
	"fmt"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [project_name]",
	Short: "Start a timer for a project",
	Long:  `Start a new time entry for a given project. If the project does not exist, it will be created.`, 
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		err := database.StartTimer(projectName)
		if err != nil {
			fmt.Printf("Error starting timer: %v\n", err)
			return
		}
		fmt.Printf("Started timer for project '%s'\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}


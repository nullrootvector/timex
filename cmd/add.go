package cmd

import (
	"fmt"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [project_name]",
	Short: "Add a new project",
	Long:  `Add a new project to the list of projects to track.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		err := database.AddProject(projectName)
		if err != nil {
			fmt.Printf("Could not add project '%s': %v\n", projectName, err)
			return
		}
		fmt.Printf("Added project '%s'\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}


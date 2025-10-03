package cmd

import (
	"fmt"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [project_name]",
	Short: "Remove a project",
	Long:  `Remove a project. This will only work if the project has no time entries associated with it.`, 
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		err := database.RemoveProject(projectName)
		if err != nil {
			fmt.Printf("Error removing project: %v\n", err)
			return
		}
		fmt.Printf("Removed project '%s'\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}


package cmd

import (
	"fmt"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all projects that have been tracked.`,
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := database.GetProjects()
		if err != nil {
			fmt.Println("Could not retrieve projects:", err)
			return
		}

		if len(projects) == 0 {
			fmt.Println("No projects found. Add one with 'timex add <project>'")
			return
		}

		fmt.Println("Projects:")
		for _, project := range projects {
			fmt.Printf("- %s\n", project)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

package cmd

import (
	"fmt"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var note string

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current timer",
	Long:  `Stop the currently running time entry.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, err := database.StopTimer(note)
		if err != nil {
			fmt.Printf("Error stopping timer: %v\n", err)
			return
		}
		fmt.Printf("Stopped timer for project '%s'\n", *projectName)
	},
}

func init() {
	stopCmd.Flags().StringVarP(&note, "note", "m", "", "Add a note to the time entry")
	rootCmd.AddCommand(stopCmd)
}

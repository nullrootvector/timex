package cmd

import (
	"fmt"
	"time"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the current timer",
	Long:  `Shows the project for the currently running timer and the elapsed time.`,
	Run: func(cmd *cobra.Command, args []string) {
		info, err := database.GetActiveTimerInfo()
		if err != nil {
			fmt.Printf("Error getting status: %v\n", err)
			return
		}

		if info == nil {
			fmt.Println("No timer is currently running.")
			return
		}

		duration := time.Since(info.StartTime)
		fmt.Printf("Timer for project '%s' has been running for %s\n", info.ProjectName, duration.Round(time.Second))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

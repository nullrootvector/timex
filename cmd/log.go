package cmd

import (
	"fmt"
	"time"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var from, to string

var logCmd = &cobra.Command{
	Use:   "log [project_name]",
	Short: "Log a time entry manually",
	Long:  `Manually add a time entry for a project with a specific start and end time.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		// Basic validation
		if from == "" || to == "" {
			fmt.Println("Error: --from and --to flags are required")
			return
		}

		// Parse times - assumes HH:MM format for the current day
		now := time.Now()
		layout := "15:04"

		startTime, err := time.ParseInLocation(layout, from, now.Location())
		if err != nil {
			fmt.Printf("Error parsing --from time: %v\n", err)
			return
		}
		startTime = time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, now.Location())

		endTime, err := time.ParseInLocation(layout, to, now.Location())
		if err != nil {
			fmt.Printf("Error parsing --to time: %v\n", err)
			return
		}
		endTime = time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, now.Location())

		if !endTime.After(startTime) {
			fmt.Println("Error: --to time must be after --from time")
			return
		}

		err = database.LogTime(projectName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error logging time: %v\n", err)
			return
		}

		fmt.Printf("Logged %s for project '%s' from %s to %s\n", endTime.Sub(startTime), projectName, from, to)
	},
}

func init() {
	logCmd.Flags().StringVar(&from, "from", "", "Start time in HH:MM format")
	logCmd.Flags().StringVar(&to, "to", "", "End time in HH:MM format")
	rootCmd.AddCommand(logCmd)
}

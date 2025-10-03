package cmd

import (
	"fmt"
	"time"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of time entries",
	Long:  `Generate a report of time entries, filterable by time range.`,
	Run: func(cmd *cobra.Command, args []string) {
		today, _ := cmd.Flags().GetBool("today")
		week, _ := cmd.Flags().GetBool("week")

		var start, end time.Time
		now := time.Now()

		if today {
			start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			end = start.Add(24 * time.Hour)
		} else if week {
			weekday := now.Weekday()
			start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -int(weekday))
			end = start.AddDate(0, 0, 7)
		} else {
			// Default to all time
			start = time.Time{}
			end = now.Add(24 * time.Hour) // A bit in the future to catch all
		}

		entries, err := database.GetTimeEntries(start, end)
		if err != nil {
			fmt.Printf("Error getting report: %v\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No time entries found for the selected period.")
			return
		}

		printReport(entries)
	},
}

func printReport(entries []database.TimeEntry) {
	projectTotals := make(map[string]time.Duration)
	var totalDuration time.Duration

	for _, entry := range entries {
		projectTotals[entry.ProjectName] += entry.Duration
		totalDuration += entry.Duration
	}

	fmt.Println("--- Time Report ---")
	for name, duration := range projectTotals {
		fmt.Printf("Project: %-20s Total Time: %s\n", name, duration.Round(time.Second))
	}
	fmt.Println("---")
	fmt.Printf("Total Tracked Time: %s\n", totalDuration.Round(time.Second))
	fmt.Println("-------------------")
}

func init() {
	reportCmd.Flags().Bool("today", false, "Report on entries from today")
	reportCmd.Flags().Bool("week", false, "Report on entries from this week")
	rootCmd.AddCommand(reportCmd)
}

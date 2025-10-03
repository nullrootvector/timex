package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nullrootvector/timex/database"
	"github.com/spf13/cobra"
)

var format string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export time entries to a file",
	Long:  `Export time entries to a machine-readable format like JSON or CSV.`,
	Run: func(cmd *cobra.Command, args []string) {
		if format != "json" && format != "csv" {
			fmt.Println("Error: --format must be either 'json' or 'csv'")
			return
		}

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
			start = time.Time{}
			end = now.Add(24 * time.Hour)
		}

		entries, err := database.GetTimeEntries(start, end)
		if err != nil {
			fmt.Printf("Error getting entries for export: %v\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No entries to export for the selected period.")
			return
		}

		switch format {
		case "json":
			exportJSON(entries)
		case "csv":
			exportCSV(entries)
		}
	},
}

func exportJSON(entries []database.TimeEntry) {
	jsonBytes, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Printf("Error exporting to JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
}

func exportCSV(entries []database.TimeEntry) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	header := []string{"Project", "Start Time", "End Time", "Duration (s)", "Notes"}
	writer.Write(header)

	for _, entry := range entries {
		record := []string{
			entry.ProjectName,
			entry.StartTime.Format(time.RFC3339),
			entry.EndTime.Format(time.RFC3339),
			fmt.Sprintf("%.0f", entry.Duration.Seconds()),
			entry.Notes,
		}
		writer.Write(record)
	}
}

func init() {
	exportCmd.Flags().StringVar(&format, "format", "", "Export format: 'json' or 'csv'")
	exportCmd.MarkFlagRequired("format")
	exportCmd.Flags().Bool("today", false, "Export entries from today")
	exportCmd.Flags().Bool("week", false, "Export entries from this week")
	rootCmd.AddCommand(exportCmd)
}

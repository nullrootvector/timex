package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ActiveTimerInfo struct {
	ProjectName string
	StartTime   time.Time
}

type TimeEntry struct {
	ProjectName string
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Notes       string
}

var DB *sql.DB

// querier can be a *sql.DB or a *sql.Tx
type querier interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func InitDatabase() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get user home directory: ", err)
	}

	dbPath := filepath.Join(home, ".timex.db")

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	createTables()
}

func createTables() {
	projectsTable := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	`
	_, err := DB.Exec(projectsTable)
	if err != nil {
		log.Fatal("Failed to create projects table: ", err)
	}

	timeEntriesTable := `
	CREATE TABLE IF NOT EXISTS time_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		start_time DATETIME NOT NULL,
		end_time DATETIME,
		notes TEXT,
		FOREIGN KEY (project_id) REFERENCES projects (id)
	);
	`
	_, err = DB.Exec(timeEntriesTable)
	if err != nil {
		log.Fatal("Failed to create time_entries table: ", err)
	}
}

func GetProjects() ([]string, error) {
	rows, err := DB.Query("SELECT name FROM projects ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		projects = append(projects, name)
	}

	return projects, nil
}

func AddProject(name string) error {
	return addProject(DB, name)
}

func addProject(q querier, name string) error {
	_, err := q.Exec("INSERT INTO projects (name) VALUES (?)", name)
	return err
}


func getProjectID(name string) (int, error) {
	return getProjectIDTx(DB, name)
}

func getProjectIDTx(q querier, name string) (int, error) {
	var id int
	err := q.QueryRow("SELECT id FROM projects WHERE name = ?", name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Project doesn't exist, create it
			if err := addProject(q, name); err != nil {
				return 0, err
			}
			// Get the ID of the newly created project
			return getProjectIDTx(q, name)
		}
		return 0, err
	}
	return id, nil
}

func GetActiveTimer() (*string, error) {
	var projectName string
	err := DB.QueryRow(`
		SELECT p.name FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.end_time IS NULL
	`).Scan(&projectName)

	if err != nil {
		if err == sql.ErrNoRows {
			// No active timer
			return nil, nil
		}
		return nil, err
	}

	return &projectName, nil
}

func StartTimer(projectName string) error {
	activeTimer, err := GetActiveTimer()
	if err != nil {
		return err
	}
	if activeTimer != nil {
		return errors.New("a timer is already running for project: " + *activeTimer)
	}

	projectID, err := getProjectID(projectName)
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO time_entries (project_id, start_time) VALUES (?, ?)", projectID, time.Now())
	return err
}

func StopTimer(note string) (*string, error) {
	var activeTimerID int
	var projectName string

	err := DB.QueryRow(`
		SELECT te.id, p.name FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.end_time IS NULL
	`).Scan(&activeTimerID, &projectName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no timer is currently running")
		}
		return nil, err
	}

	_, err = DB.Exec("UPDATE time_entries SET end_time = ?, notes = ? WHERE id = ?", time.Now(), note, activeTimerID)
	return &projectName, nil
}

func GetActiveTimerInfo() (*ActiveTimerInfo, error) {
	var info ActiveTimerInfo
	err := DB.QueryRow(`
		SELECT p.name, te.start_time FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.end_time IS NULL
	`).Scan(&info.ProjectName, &info.StartTime)

	if err != nil {
		if err == sql.ErrNoRows {
			// No active timer
			return nil, nil
		}
		return nil, err
	}

	return &info, nil
}

func GetTimeEntries(start time.Time, end time.Time) ([]TimeEntry, error) {
	rows, err := DB.Query(`
		SELECT p.name, te.start_time, te.end_time, te.notes
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.end_time IS NOT NULL AND te.start_time >= ? AND te.end_time <= ?
		ORDER BY te.start_time DESC
	`, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var entry TimeEntry
		var notes sql.NullString
		if err := rows.Scan(&entry.ProjectName, &entry.StartTime, &entry.EndTime, &notes); err != nil {
			return nil, err
		}
		entry.Duration = entry.EndTime.Sub(entry.StartTime)
		if notes.Valid {
			entry.Notes = notes.String
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func SwitchTimer(newProjectName string) (*string, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}

	// 1. Stop the current timer
	var activeTimerID int
	var oldProjectName string
	err = tx.QueryRow(`
		SELECT te.id, p.name FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.end_time IS NULL
	`).Scan(&activeTimerID, &oldProjectName)

	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			return nil, errors.New("no timer is currently running")
		}
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("UPDATE time_entries SET end_time = ? WHERE id = ?", time.Now(), activeTimerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. Start a new timer
	newProjectID, err := getProjectIDTx(tx, newProjectName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO time_entries (project_id, start_time) VALUES (?, ?)", newProjectID, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &oldProjectName, tx.Commit()
}

func RemoveProject(projectName string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Get project ID
	var projectID int
	err = tx.QueryRow("SELECT id FROM projects WHERE name = ?", projectName).Scan(&projectID)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return errors.New("project not found")
		}
		return err
	}

	// Check for time entries
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM time_entries WHERE project_id = ?", projectID).Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}

	if count > 0 {
		tx.Rollback()
		return fmt.Errorf("cannot remove project with %d associated time entries", count)
	}

	// Delete the project
	_, err = tx.Exec("DELETE FROM projects WHERE id = ?", projectID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func LogTime(projectName string, startTime time.Time, endTime time.Time) error {
	projectID, err := getProjectID(projectName)
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO time_entries (project_id, start_time, end_time) VALUES (?, ?, ?)", projectID, startTime, endTime)
	return err
}
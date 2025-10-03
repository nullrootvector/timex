
package database

import (
	"database/sql"
	"errors"
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

var DB *sql.DB

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
	_, err := DB.Exec("INSERT INTO projects (name) VALUES (?)", name)
	return err
}

func getProjectID(name string) (int, error) {
	var id int
	err := DB.QueryRow("SELECT id FROM projects WHERE name = ?", name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Project doesn't exist, create it
			if err := AddProject(name); err != nil {
				return 0, err
			}
			// Get the ID of the newly created project
			return getProjectID(name)
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
	if err != nil {
		return nil, err
	}

	return &projectName, nil
}

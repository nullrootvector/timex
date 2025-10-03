
package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

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

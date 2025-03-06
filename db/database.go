package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	stmt := `
	CREATE TABLE IF NOT EXISTS todos (
    	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    	title TEXT NOT NULL,
    	description TEXT NOT NULL,
		notes TEXT,
    	due_date TIMESTAMP,
    	status TEXT DEFAULT 'pending'
	);`
	_, err = DB.Exec(stmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, stmt)
	}
}

// CloseDB closes the database connection
func CloseDB() {
	DB.Close()
}

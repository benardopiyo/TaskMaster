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
        CREATE TABLE IF NOT EXISTS users (
                id TEXT NOT NULL PRIMARY KEY, -- Changed to TEXT for UUID
                username TEXT NOT NULL UNIQUE,
                password TEXT NOT NULL
        );
        CREATE TABLE IF NOT EXISTS profiles (
                id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                user_id TEXT NOT NULL, -- Changed to TEXT for UUID
                name TEXT,
                email TEXT,
                image_path TEXT, -- Added image_path
                FOREIGN KEY (user_id) REFERENCES users(id)
        );
        CREATE TABLE IF NOT EXISTS todos (
                id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                user_id TEXT NOT NULL, -- Changed to TEXT for UUID
                title TEXT NOT NULL,
                description TEXT NOT NULL,
                notes TEXT,
                due_date TIMESTAMP,
                status TEXT DEFAULT 'pending',
                FOREIGN KEY (user_id) REFERENCES users(id)
        );
        CREATE TABLE IF NOT EXISTS completed_todos (
                id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                user_id TEXT NOT NULL, -- Changed to TEXT for UUID
                title TEXT NOT NULL,
                description TEXT NOT NULL,
                notes TEXT,
                due_date TIMESTAMP,
                completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (user_id) REFERENCES users(id)
        );
        `
	_, err = DB.Exec(stmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, stmt)
	}
}

// CloseDB closes the database connection
func CloseDB() {
	DB.Close()
}

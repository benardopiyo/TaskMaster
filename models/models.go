package models

import "time"

// User represents a user in the system
type User struct {
	ID       string
	Username string
	Password string
}

// Profile represents a user profile
type Profile struct {
	ID        int
	UserID    string
	Name      string
	Email     string
	ImagePath string // image_path
}

// Todo represents a task in the todo list
type Todo struct {
	ID          int
	UserID      string
	Title       string
	Description string
	Notes       string
	DueDate     time.Time
	Status      string
}

// CompletedTodo represents a completed task
type CompletedTodo struct {
	ID          int
	UserID      string
	Title       string
	Description string
	Notes       string
	DueDate     time.Time
	CompletedAt time.Time
}

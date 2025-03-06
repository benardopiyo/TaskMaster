package models

import "time"

// Todo represents a task in the todo list
type Todo struct {
	ID          int
	Title       string
	Description string
	Notes       string
	DueDate     time.Time
	Status      string
}

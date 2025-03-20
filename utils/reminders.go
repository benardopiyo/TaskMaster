package utils

import (
	"log"
	"time"

	"todo-app/db"
)

// Check for due tasks and send reminders
func CheckReminders() {
	rows, err := db.DB.Query("SELECT id, title, due_date FROM todos WHERE status = 'pending'")
	if err != nil {
		log.Println("Error fetching reminders:", err)
		return
	}
	defer rows.Close()

	now := time.Now()
	for rows.Next() {
		var id int
		var title string
		var dueDate time.Time

		err := rows.Scan(&id, &title, &dueDate)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// Calculate the time remaining until the due date
		timeRemaining := dueDate.Sub(now)

		// If the task is due within 3 days
		if timeRemaining <= 3*24*time.Hour {
			// Send an alert message to the user
			log.Printf("Alert: Task '%s' is due in less than 3 days!\n", title)
			// You can implement additional logic here to display the alert message on the index.html page
		}
	}
}

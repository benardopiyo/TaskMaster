package utils

import (
	"fmt"
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

		// If the task is due within the next hour
		if dueDate.Sub(now) <= time.Hour {
			fmt.Printf("Reminder: Task '%s' is due soon!\n", title)
			SendEmail("benardopiyo13@gmail.com", title, dueDate) // Send email alert
		}
	}
}

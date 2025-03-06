package utils

import (
	"fmt"
	"net/smtp"
	"time"
)

// SendEmail sends a reminder email
func SendEmail(to string, taskTitle string, dueDate time.Time) {
	from := "benardopiyo.com"
	password := "benkopiyo"

	msg := fmt.Sprintf("Subject: Task Reminder\n\nYour task '%s' is due at %s!", taskTitle, dueDate.Format("2006-01-02 15:04"))
	auth := smtp.PlainAuth("", from, password, "smtp.example.com")

	err := smtp.SendMail("smtp.example.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println("Failed to send email:", err)
	}
}

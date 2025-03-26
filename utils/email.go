package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"time"
)

// Email Credentials
const (
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
	senderEmail  = "benardopiyo13@gmail.com"
	senderPass   = "benkopiyo"
)

// SendEmail sends an email reminder
func SendEmail(to, taskTitle string, dueDate time.Time) {
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)

	subject := "Subject: Task Reminder\n"
	body := fmt.Sprintf("Your task '%s' is due on %s!\n", taskTitle, dueDate.Format("2006-01-02"))
	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{to}, msg)
	if err != nil {
		log.Println("❌ Failed to send email:", err)
		return
	}
	log.Printf("✅ Reminder email sent to %s for task: %s\n", to, taskTitle)
}

package main

import (
	"fmt"
	"log"
	"math"
	"net/smtp"
	"time"

	"github.com/muhammadolammi/rentradar/internal/database"
)

func sendEmail(to, subject, body string, smtpModel SMTPModel) error {
	from := smtpModel.UserName
	password := smtpModel.Password
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpModel.Server)
	err := smtp.SendMail("smtp.example.com:587", auth, from, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}

func sendWithRetry(senderFunction func() error, retries int) error {
	var err error
	for i := range retries {
		err = senderFunction()
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}
	return err
}

// Dispatch notification based on contact type
func SendNotification(notification Notification, smtpMdel SMTPModel) error {
	var senderFunc func() error

	switch notification.ContactMethod {
	case "email":
		senderFunc = func() error {
			return sendEmail(notification.Contact, notification.Subject, notification.Body, smtpMdel)
		}
	case "whatsapp":
		senderFunc = func() error {
			// replace with actual WhatsApp sending logic
			log.Println("Simulating WhatsApp send...")
			return fmt.Errorf("WhatsApp not implemented yet")
		}
	case "sms":
		senderFunc = func() error {
			// replace with actual SMS sending logic
			log.Println("Simulating SMS send...")
			return fmt.Errorf("SMS not implemented yet")
		}
	default:
		return fmt.Errorf("unknown contact method: %s", notification.ContactMethod)
	}

	// Call generic retry logic
	return sendWithRetry(senderFunc, 3)
}

// Notification  Model Helper
func DbNotificationToModelsNotification(dbNotification database.Notification) Notification {
	return Notification{
		ID:            dbNotification.ID,
		UserID:        dbNotification.UserID,
		Status:        dbNotification.Status,
		ListingID:     dbNotification.ListingID,
		SentAt:        dbNotification.SentAt,
		ContactMethod: dbNotification.ContactMethod,
		Contact:       dbNotification.Contact,
		Subject:       dbNotification.Subject,
		Body:          dbNotification.Body,
	}
}

func DbNotificationsToModelsNotifications(dbNotifications []database.Notification) []Notification {
	notications := []Notification{}
	for _, dbNotification := range dbNotifications {
		notications = append(notications, DbNotificationToModelsNotification(dbNotification))
	}
	return notications
}

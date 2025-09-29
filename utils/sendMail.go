package utils

import (
	"gopkg.in/gomail.v2"
	"log"
	"github.com/joho/godotenv"
	"os"
	
)

func SendEmail(to string, subject string, body string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var APPS_PASS string = os.Getenv("APPS_PASS_EMAIL")
	var EMAIL_SENDER string = os.Getenv("EMAIL_SENDER")
	m := gomail.NewMessage()
	m.SetHeader("From", EMAIL_SENDER) 
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, EMAIL_SENDER, APPS_PASS) 
	// APPS_PASS bukan pass email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully to:", to)
	return nil
}

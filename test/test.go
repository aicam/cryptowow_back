package main

import (
	"log"
	"net/smtp"
)

func main() {
	// Configuration
	from := "doNotReply@cryptowow.com"
	password := "super_secret_password"
	to := []string{"ali.risheh876@gmail.com"}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("My super secret message.")

	// Create authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send actual message
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
}

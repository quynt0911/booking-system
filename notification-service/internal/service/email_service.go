package service

import "log"

func SendEmail(to, subject, body string) error {
	log.Printf("Sending email to %s with subject '%s'", to, subject)
	return nil
}
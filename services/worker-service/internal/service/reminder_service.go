package service

import "log"

func SendReminder(userID, time string) {
	log.Printf("Sending reminder to user %s at %s", userID, time)
}
package service

import (
	"fmt"
	"net/smtp"

	"github.com/your-org/booking-system/services/notification-service/internal/config"
)

type EmailService struct {
	cfg config.Config
}

func NewEmailService(cfg config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) Send(to, subject, body string) error {
	from := s.cfg.SMTPFrom
	msg := []byte(fmt.Sprintf("Subject: %s\r\nFrom: %s\r\nTo: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n%s", subject, from, to, body))
	host := s.cfg.SMTPHost
	port := s.cfg.SMTPPort
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, host)
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

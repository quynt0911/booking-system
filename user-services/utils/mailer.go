package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

func SendWelcomeEmail(to string, name string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "your_email@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Chào mừng bạn đến với hệ thống tư vấn")
	m.SetBody("text/plain", fmt.Sprintf("Xin chào %s, cảm ơn bạn đã đăng ký!", name))

	d := gomail.NewDialer("smtp.gmail.com", 587, "your_email@gmail.com", "your_app_password")

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Gửi email lỗi:", err)
		return
	}

	fmt.Println("📨 Đã gửi email cho", to)
}

package interfaces

import (
	"login-api/models"
	"os"

	"gopkg.in/gomail.v2"
)

type SMTPEmailService struct {
	dialer *gomail.Dialer
}

func NewSMTPEmailService() *SMTPEmailService {
	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)
	return &SMTPEmailService{dialer: dialer}
}

func (s *SMTPEmailService) SendWelcomeEmail(user models.UserCreatedEvent) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your-app@example.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Welcome to Our Platform!")
	m.SetBody("text/html", `
        <h1>Welcome to Our Platform!</h1>
        <p>Dear `+user.Name+`,</p>
        <p>Thank you for registering with us. We're excited to have you on board!</p>
        <p>If you have any questions, feel free to reach out to our support team.</p>
        <p>Best regards,<br>Your App Team</p>
    `)
	return s.dialer.DialAndSend(m)
}

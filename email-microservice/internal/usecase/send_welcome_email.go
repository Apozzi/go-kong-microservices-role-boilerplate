package usecase

import "login-api/models"

type EmailService interface {
	SendWelcomeEmail(user models.UserCreatedEvent) error
}

type SendWelcomeEmailUseCase struct {
	EmailService EmailService
}

func (u *SendWelcomeEmailUseCase) Execute(event models.UserCreatedEvent) error {
	return u.EmailService.SendWelcomeEmail(event)
}

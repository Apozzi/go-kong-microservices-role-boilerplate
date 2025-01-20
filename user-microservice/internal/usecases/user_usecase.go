package usecases

import (
	"encoding/json"
	"errors"
	"login-api/internal/repositories"
	"login-api/models"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type UserUseCase struct {
	repo     repositories.UserRepository
	amqpChan *amqp091.Channel
}

func NewUserUseCase(repo repositories.UserRepository, amqpChan *amqp091.Channel) *UserUseCase {
	return &UserUseCase{
		repo:     repo,
		amqpChan: amqpChan,
	}
}

func (uc *UserUseCase) GetAll() ([]models.User, error) {
	return uc.repo.FindAll()
}

func (uc *UserUseCase) GetByID(id uint64) (*models.User, error) {
	return uc.repo.GetUserWithRoles(id)
}

func (uc *UserUseCase) Create(user *models.User) error {
	existingUser, _ := uc.repo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("user already registered")
	}

	user.RegisterDate = time.Now()
	err := uc.repo.Create(user)
	if err != nil {
		return err
	}

	event := models.UserCreatedEvent{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = uc.amqpChan.Publish(
		"user_events",
		"user.created",
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)

	return err
}

func (uc *UserUseCase) Update(id uint64, user *models.User) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.ID = id
	return uc.repo.Update(user)
}

func (uc *UserUseCase) Delete(id uint64) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return uc.repo.Delete(id)
}

func (uc *UserUseCase) AddRoleToUser(userID uint64, roleID string) error {
	user, err := uc.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return uc.repo.AddRole(user.ID, roleID)
}

func (uc *UserUseCase) RemoveRoleFromUser(userID uint64, roleID string) error {
	user, err := uc.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return uc.repo.RemoveRole(user.ID, roleID)
}

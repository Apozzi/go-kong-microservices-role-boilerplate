package handlers

import (
	"encoding/json"
	"log"
	"login-api/internal/usecase"
	"login-api/models"

	"github.com/rabbitmq/amqp091-go"
)

type EventHandler struct {
	useCase *usecase.SendWelcomeEmailUseCase
}

func NewEventHandler(useCase *usecase.SendWelcomeEmailUseCase) *EventHandler {
	return &EventHandler{useCase: useCase}
}

func (h *EventHandler) HandleMessage(msg amqp091.Delivery) {
	var event models.UserCreatedEvent
	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		log.Printf("Error deserializing message: %v", err)
		msg.Nack(false, false)
		return
	}

	err = h.useCase.Execute(event)
	if err != nil {
		log.Printf("Error processing event: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
	log.Printf("Successfully processed event for user: %s", event.Email)
}

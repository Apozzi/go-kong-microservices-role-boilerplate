package main

import (
	"log"
	"login-api/internal/handlers"
	"login-api/internal/interfaces"
	"login-api/internal/usecase"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() (*amqp091.Connection, *amqp091.Channel) {
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		"user_events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	return conn, ch
}

func main() {
	godotenv.Load()
	conn, ch := ConnectRabbitMQ()
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"welcome_email_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name,
		"user.created",
		"user_events",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	emailService := interfaces.NewSMTPEmailService()
	useCase := &usecase.SendWelcomeEmailUseCase{EmailService: emailService}
	eventHandler := handlers.NewEventHandler(useCase)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			eventHandler.HandleMessage(msg)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// @title           Login Role Boilerplate
// @version         1.0
// @description     Exemplo de servidor com login e autenticação com roles.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Anderson Rodrigo Pozzi
// @contact.url    adeveloper.com.br
// @contact.email  eanderea1@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api/v1

package main

import (
	"log"
	config "login-api/internal/config"
	controllers "login-api/internal/controllers"
	"login-api/internal/repositories"
	routers "login-api/internal/routers"
	"login-api/internal/usecases"
	middleware "login-api/middleware"
	"login-api/models"
	"os"
	"reflect"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	_ "login-api/docs"
	/*
		Adicionar os modulos
		go get -u github.com/swaggo/gin-swagger
		go get -u github.com/swaggo/files
	*/
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AutoMigrate(db *gorm.DB) {
	models := []interface{}{
		&models.User{},
		&models.Role{},
		&models.UserRole{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Printf("Failed to migrate model %v: %v", reflect.TypeOf(model), err)
		} else {
			log.Printf("Successfully migrated model %v", reflect.TypeOf(model))
		}
	}
}

func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	// Adicionar RabbitMq go get github.com/rabbitmq/amqp091-go
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
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
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}

func main() {
	godotenv.Load()
	db := config.Connect()
	AutoMigrate(db)
	var secretKey string = os.Getenv("SECRET_KEY")
	jwtAuth, _ := middleware.NewJWTTokenMaker(secretKey)
	controllers.Initialize(config.Connect(), jwtAuth)
	router := gin.Default()
	rabbitmqConn, rabbitmqChan, err := setupRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to setup RabbitMQ: %v", err)
	}
	defer rabbitmqConn.Close()
	defer rabbitmqChan.Close()

	// Repositories
	userRepo := repositories.NewGormUserRepository(db)

	// Use Cases
	userUseCase := usecases.NewUserUseCase(userRepo, rabbitmqChan)

	// Controllers
	userController := controllers.NewUserController(userUseCase)

	routers.Routers(router, userController)
	// Para acessar o swagger: http://localhost:8081/swagger/index.html#/
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := os.Getenv("PORT")
	router.Run(":" + port)
}

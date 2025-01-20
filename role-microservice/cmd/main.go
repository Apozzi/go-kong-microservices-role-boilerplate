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
		&models.Role{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Printf("Failed to migrate model %v: %v", reflect.TypeOf(model), err)
		} else {
			log.Printf("Successfully migrated model %v", reflect.TypeOf(model))
		}
	}
}

func main() {
	godotenv.Load()
	db := config.Connect()
	AutoMigrate(db)
	var secretKey string = os.Getenv("SECRET_KEY")
	jwtAuth, _ := middleware.NewJWTTokenMaker(secretKey)
	controllers.Initialize(config.Connect(), jwtAuth)
	router := gin.Default()

	// Repository
	roleRepo := repositories.NewGormRoleRepository(db)

	// Use Case
	roleUseCase := usecases.NewRoleUseCase(roleRepo)

	// Controller
	roleController := controllers.NewRoleController(roleUseCase)

	routers.Routers(router, roleController)
	// Para acessar o swagger: http://localhost:8081/swagger/index.html#/
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := os.Getenv("PORT")
	router.Run(":" + port)
}

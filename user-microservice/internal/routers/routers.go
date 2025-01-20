package routers

import (
	"login-api/internal/controllers"
	middleware "login-api/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	ROLE_ADMIN    = "Admin"
	ROLE_MODIFIER = "Modifier"
	ROLE_WATCHER  = "Watcher"
)

func Routers(router *gin.Engine, userController *controllers.UserController) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Version 1")
	})
	router.POST("login", controllers.Login)

	privateRoute := router.Group("")
	privateRoute.Use(controllers.Authenticate)
	{
		privateRoute.GET("verifyToken", controllers.VerifyToken)

		userRoutes := privateRoute.Group("user")
		{
			userRoutes.GET("", middleware.RequireRoles(ROLE_ADMIN), userController.GetUsers)
			userRoutes.GET(":id", middleware.RequireRoles(ROLE_ADMIN), userController.GetUser)
			userRoutes.POST("", middleware.RequireRoles(ROLE_ADMIN), userController.CreateUser)
			userRoutes.PUT(":id", middleware.RequireRoles(ROLE_ADMIN), userController.UpdateUser)
			userRoutes.DELETE(":id", middleware.RequireRoles(ROLE_ADMIN), userController.DeleteUser)

			userRoutes.POST(":userId/roles/:roleId", middleware.RequireRoles(ROLE_ADMIN), userController.AddRoleToUser)
			userRoutes.DELETE("remove/:userId/roles/:roleId", middleware.RequireRoles(ROLE_ADMIN), userController.RemoveRoleFromUser)
		}
	}
}

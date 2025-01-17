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

func Routers(router *gin.Engine) {
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
			userRoutes.GET("", middleware.RequireRoles(ROLE_ADMIN), controllers.GetUsers)
			userRoutes.GET(":id", middleware.RequireRoles(ROLE_ADMIN), controllers.GetUser)
			userRoutes.POST("", middleware.RequireRoles(ROLE_ADMIN), controllers.PostUser)
			userRoutes.PUT(":id", middleware.RequireRoles(ROLE_ADMIN), controllers.PutUser)
			userRoutes.DELETE(":id", middleware.RequireRoles(ROLE_ADMIN), controllers.DeleteUser)

			userRoutes.POST(":userId/roles/:roleId", middleware.RequireRoles(ROLE_ADMIN), controllers.AddRoleToUser)
			userRoutes.DELETE("remove/:userId/roles/:roleId", middleware.RequireRoles(ROLE_ADMIN), controllers.RemoveRoleFromUser)
		}
	}
}

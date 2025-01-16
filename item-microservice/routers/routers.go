package routers

import (
	"login-api/controllers"
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

	privateRoute := router.Group("")
	privateRoute.Use(controllers.Authenticate)
	{

		itemRoutes := privateRoute.Group("item")
		{
			itemRoutes.GET("", middleware.RequireRoles(ROLE_ADMIN, ROLE_MODIFIER, ROLE_WATCHER), controllers.GetItems)
			itemRoutes.GET(":id", middleware.RequireRoles(ROLE_ADMIN, ROLE_MODIFIER, ROLE_WATCHER), controllers.GetItem)
			itemRoutes.POST("", middleware.RequireRoles(ROLE_ADMIN, ROLE_MODIFIER), controllers.PostItem)
			itemRoutes.PUT(":id", middleware.RequireRoles(ROLE_ADMIN, ROLE_MODIFIER), controllers.PutItem)
			itemRoutes.DELETE(":id", middleware.RequireRoles(ROLE_ADMIN, ROLE_MODIFIER), controllers.DeleteItem)
		}
	}
}

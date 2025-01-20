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

func Routers(router *gin.Engine, roleController *controllers.RoleController) {
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
		roleRoutes := privateRoute.Group("role")
		{
			roleRoutes.GET("", middleware.RequireRoles(ROLE_ADMIN), roleController.GetRoles)
			roleRoutes.GET(":id", middleware.RequireRoles(ROLE_ADMIN), roleController.GetRole)
			roleRoutes.POST("", middleware.RequireRoles(ROLE_ADMIN), roleController.CreateRole)
			roleRoutes.PUT(":id", middleware.RequireRoles(ROLE_ADMIN), roleController.UpdateRole)
			roleRoutes.DELETE(":id", middleware.RequireRoles(ROLE_ADMIN), roleController.DeleteRole)
		}
	}
}

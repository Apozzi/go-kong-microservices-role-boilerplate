// @title Login API
// @version 1.0
// @description API for user authentication and JWT token generation.
// @host localhost:8081
// @BasePath /

package controllers

import (
	"time"

	models "login-api/models"

	"github.com/gin-gonic/gin"
)

// Login authenticates a user.
// @Summary Authenticate user
// @Description Validates user credentials and returns a JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param userLogin body models.UserLogin true "User credentials"
// @Success 200 {object} map[string]interface{} "JWT token and user information"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Failed to create token"
// @Router /login [post]
func Login(c *gin.Context) {
	var userLogin models.UserLogin
	c.BindJSON(&userLogin)
	var user models.User
	result := db.Preload("Roles").Where("email = ? AND password = ?", userLogin.Username, userLogin.Password).First(&user)

	if result.Error != nil {
		c.JSON(401, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	if user.ID != 0 {
		roles := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roles[i] = role.Name
		}

		token, err := auth.CreateTokenWithRoles(user.Email, roles, time.Hour*12)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to create token",
			})
			return
		}
		c.JSON(200, gin.H{
			"token":   token,
			"user_id": user.ID,
			"roles":   roles,
		})
	} else {
		c.JSON(401, gin.H{
			"error": "Invalid credentials",
		})
	}
}

// VerifyToken checks the validity of a JWT token.
// @Summary Verify token
// @Description Checks if the provided JWT token is valid.
// @Tags Authentication
// @Produce json
// @Success 200 {object} map[string]bool "Token is valid"
// @Router /verify-token [get]
func VerifyToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": true,
	})
}

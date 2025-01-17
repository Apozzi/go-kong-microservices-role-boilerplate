package controllers

import (
	middleware "login-api/middleware"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Response representa uma resposta genérica da API
type Response struct {
	Data    interface{}       `json:"data,omitempty"`
	Message string            `json:"message,omitempty"`
	Error   string            `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

var db *gorm.DB
var auth middleware.Auth

func Initialize(dbConnection *gorm.DB, authService middleware.Auth) {
	db = dbConnection
	auth = authService
}

func Authenticate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	payload, err := auth.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("Email", payload.Username)
	c.Set("Roles", payload.Roles)
	c.Next()
}

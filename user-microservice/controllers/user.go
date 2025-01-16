// @title User Management API
// @version 1.0
// @description API for managing users and roles
// @host localhost:8082
// @BasePath /

package controllers

import (
	"database/sql"
	"fmt"
	models "login-api/models"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// @Summary Get all users
// @Description Get a list of all users with their roles
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]models.User} "Success"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	result := db.Joins("JOIN user_roles ur ON ur.user_id = users.id").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Preload("Roles").Find(&users)
	if result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Failed to retrieve users",
		})
		return
	}

	c.JSON(200, Response{
		Data: users,
	})
}

// @Summary Get user by ID
// @Description Get a single user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} Response{data=models.User} "Success"
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	db.Where("id_user = @idUser", sql.Named("idUser", id)).Find(&user)
	c.JSON(200, Response{
		Data: user,
	})
}

// @Summary Create new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 200 {object} Response{data=models.User} "Success"
// @Failure 400 {object} ErrorResponse "Validation Error"
// @Failure 409 {object} ErrorResponse "User Already Exists"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users [post]
func PostUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			switch fieldName {
			case "Name":
				errorMessages["name"] = "Name is required"
			case "Role":
				errorMessages["role"] = "Role is required"
			case "Email":
				if fieldErr.Tag() == "required" {
					errorMessages["email"] = "Email is required"
				} else if fieldErr.Tag() == "email" {
					errorMessages["email"] = "Invalid email format"
				}
			}
		}
		c.JSON(400, ErrorResponse{
			Errors: errorMessages,
		})
		return
	}
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(409, ErrorResponse{ // Verifica cadastro de usuários repetidos através do Email.
			Error: "User already registered",
		})
		return
	}
	user.RegisterDate = time.Now()
	if result := db.Create(&user); result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Failed to create user",
		})
		return
	}
	c.JSON(200, Response{
		Data: user,
	})
}

// @Summary Update user
// @Description Update an existing user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "Updated user information"
// @Success 200 {object} Response{data=models.User} "Success"
// @Router /users/{id} [put]
func PutUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	c.BindJSON(&user)
	user.ID, _ = strconv.ParseUint(id, 10, 64)
	db.Save(&user)
	c.JSON(200, Response{
		Data: user,
	})
}

// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} Response{data=models.User} "Success"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	db.Where("id_user = @idUser", sql.Named("idUser", id)).Delete(&user)
	c.JSON(200, Response{
		Data: user,
	})
}

func getUserAndRole(userID, roleID int) (models.User, models.Role, error) {
	var user models.User
	var role models.Role
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return user, role, fmt.Errorf("user not found")
	}
	if err := db.Where("id = ?", roleID).First(&role).Error; err != nil {
		return user, role, fmt.Errorf("role not found")
	}

	return user, role, nil
}

// @Summary Add role to user
// @Description Associate a role with a user
// @Tags user-roles
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param roleId path string true "Role ID"
// @Success 200 {object} Response{message=string,user=models.User,role=models.Role} "Success"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 404 {object} ErrorResponse "User or Role Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/{userId}/roles/{roleId} [post]
func AddRoleToUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("roleId"))
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid role ID"})
		return
	}

	user, role, err := getUserAndRole(userID, roleID)
	if err != nil {
		c.JSON(404, ErrorResponse{Error: err.Error()})
		return
	}

	if err := db.Model(&user).Association("Roles").Append(&role); err != nil {
		c.JSON(500, ErrorResponse{Error: "Failed to associate role with user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Role successfully added to user",
		"user":    user,
		"role":    role,
	})
}

// @Summary Remove role from user
// @Description Remove a role association from a user
// @Tags user-roles
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param roleId path string true "Role ID"
// @Success 200 {object} Response{message=string,user=models.User,role=models.Role} "Success"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 404 {object} ErrorResponse "User or Role Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/{userId}/roles/{roleId} [delete]
func RemoveRoleFromUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("roleId"))
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid role ID"})
		return
	}

	user, role, err := getUserAndRole(userID, roleID)
	if err != nil {
		c.JSON(404, ErrorResponse{Error: err.Error()})
		return
	}

	if err := db.Model(&user).Association("Roles").Delete(&role); err != nil {
		c.JSON(500, ErrorResponse{Error: "Failed to remove role from user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Role successfully removed from user",
		"user":    user,
		"role":    role,
	})
}

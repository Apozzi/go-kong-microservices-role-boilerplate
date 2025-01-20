// @title User Management API
// @version 1.0
// @description API for managing users and roles
// @host localhost:8082
// @BasePath /

package controllers

import (
	"login-api/internal/usecases"
	models "login-api/models"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase *usecases.UserUseCase
}

func NewUserController(userUseCase *usecases.UserUseCase) *UserController {
	return &UserController{userUseCase: userUseCase}
}

// @Summary Get all users
// @Description Get a list of all users with their roles
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]domains.User} "Success"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users [get]
func (ctrl *UserController) GetUsers(c *gin.Context) {
	users, err := ctrl.userUseCase.GetAll()
	if err != nil {
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
// @Success 200 {object} Response{data=domains.User} "Success"
// @Failure 404 {object} ErrorResponse{error=string} "User Not Found"
// @Router /users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	user, err := ctrl.userUseCase.GetByID(id)
	if err != nil {
		c.JSON(404, ErrorResponse{
			Error: "User not found",
		})
		return
	}

	c.JSON(200, Response{
		Data: user,
	})
}

// @Summary Create new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body domains.User true "User information"
// @Success 200 {object} Response{data=domains.User} "Success"
// @Failure 400 {object} ErrorResponse "Validation Error"
// @Failure 409 {object} ErrorResponse "User Already Exists"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				switch fieldErr.Field() {
				case "Name":
					errorMessages["name"] = "Name is required"
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
	}

	err := ctrl.userUseCase.Create(&user)
	if err != nil {
		if err.Error() == "user already registered" {
			c.JSON(409, ErrorResponse{
				Error: "User already registered",
			})
			return
		}
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
// @Param user body domains.User true "Updated user information"
// @Success 200 {object} Response{data=domains.User} "Success"
// @Failure 404 {object} ErrorResponse{error=string} "User Not Found"
// @Failure 400 {object} ErrorResponse{error=string} "Invalid Data"
// @Router /users/{id} [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, ErrorResponse{
			Error: "Invalid data provided",
		})
		return
	}

	err = ctrl.userUseCase.Update(id, &user)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(404, ErrorResponse{
				Error: "User not found",
			})
			return
		}
		c.JSON(500, ErrorResponse{
			Error: "Failed to update user",
		})
		return
	}

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
// @Success 200 {object} Response{message=string} "Success"
// @Failure 404 {object} ErrorResponse{error=string} "User Not Found"
// @Failure 500 {object} ErrorResponse{error=string} "Delete Failed"
// @Router /users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	err = ctrl.userUseCase.Delete(id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(404, ErrorResponse{
				Error: "User not found",
			})
			return
		}
		c.JSON(500, ErrorResponse{
			Error: "Failed to delete user",
		})
		return
	}

	c.JSON(200, Response{
		Message: "User successfully deleted",
	})
}

// @Summary Add role to user
// @Description Associate a role with a user
// @Tags user-roles
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param roleId path string true "Role ID"
// @Success 200 {object} Response{message=string} "Success"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 404 {object} ErrorResponse "User or Role Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/{userId}/roles/{roleId} [post]
func (ctrl *UserController) AddRoleToUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	roleID := c.Param("roleId")

	err = ctrl.userUseCase.AddRoleToUser(userID, roleID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(404, ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(500, ErrorResponse{Error: "Failed to associate role with user"})
		return
	}

	c.JSON(200, Response{
		Message: "Role successfully added to user",
	})
}

// @Summary Remove role from user
// @Description Remove a role association from a user
// @Tags user-roles
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param roleId path string true "Role ID"
// @Success 200 {object} Response{message=string} "Success"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 404 {object} ErrorResponse "User or Role Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/{userId}/roles/{roleId} [delete]
func (ctrl *UserController) RemoveRoleFromUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	roleID := c.Param("roleId")

	err = ctrl.userUseCase.RemoveRoleFromUser(userID, roleID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(404, ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(500, ErrorResponse{Error: "Failed to remove role from user"})
		return
	}

	c.JSON(200, Response{
		Message: "Role successfully removed to user",
	})
}

// @title Role Management API
// @version 1.0
// @description API for managing roles
// @host localhost:8081
// @BasePath /

package controllers

import (
	"database/sql"
	models "login-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary Get all roles
// @Description Get a list of all roles
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]models.Role} "Success"
// @Failure 500 {object} ErrorResponse{error=string} "Internal Server Error"
// @Router /roles [get]
func GetRoles(c *gin.Context) {
	var roles []models.Role
	result := db.Find(&roles)
	if result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Failed to retrieve roles",
		})
		return
	}

	c.JSON(200, Response{
		Data: roles,
	})
}

// @Summary Get role by ID
// @Description Get a single role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} Response{data=models.Role} "Success"
// @Router /roles/{id} [get]
func GetRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	db.Where("id = @id", sql.Named("id", id)).Find(&role)
	c.JSON(200, Response{
		Data: role,
	})
}

// @Summary Create new role
// @Description Create a new role
// @Tags roles
// @Accept json
// @Produce json
// @Param role body models.Role true "Role information"
// @Success 200 {object} Response{data=models.Role} "Success"
// @Failure 400 {object} ErrorResponse{errors=map[string]string} "Validation Error"
// @Failure 409 {object} ErrorResponse{error=string} "Role Already Exists"
// @Failure 500 {object} ErrorResponse{error=string} "Internal Server Error"
func PostRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			switch fieldName {
			case "Name":
				errorMessages["name"] = "Role name is required"
			case "Description":
				errorMessages["description"] = "Role description is required"
			}
		}
		c.JSON(400, ErrorResponse{
			Errors: errorMessages,
		})
		return
	}

	var existingRole models.Role
	if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err == nil {
		c.JSON(409, ErrorResponse{
			Error: "Role already exists",
		})
		return
	}

	if result := db.Create(&role); result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Failed to create role",
		})
		return
	}

	c.JSON(200, Response{
		Data: role,
	})
}

// @Summary Update role
// @Description Update an existing role's information
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body models.Role true "Updated role information"
// @Success 200 {object} Response{data=models.Role} "Success"
// @Failure 404 {object} ErrorResponse{error=string} "Role Not Found"
// @Failure 400 {object} ErrorResponse{error=string} "Invalid Data"
// @Failure 500 {object} ErrorResponse{error=string} "Update Failed"
func PutRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	if err := db.Where("id = ?", id).First(&role).Error; err != nil {
		c.JSON(404, ErrorResponse{
			Error: "Role not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, ErrorResponse{
			Error: "Invalid data provided",
		})
		return
	}
	if result := db.Save(&role); result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Failed to update role",
		})
		return
	}

	c.JSON(200, Response{
		Data: role,
	})
}

// @Summary Delete role
// @Description Delete a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} Response{message=string} "Success"
// @Failure 404 {object} ErrorResponse{error=string} "Role Not Found"
// @Failure 400 {object} ErrorResponse{error=string} "Role In Use"
// @Failure 500 {object} ErrorResponse{error=string} "Delete Failed"
// @Router /roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	if err := db.Where("id = ?", id).First(&role).Error; err != nil {
		c.JSON(404, ErrorResponse{
			Error: "Role not found",
		})
		return
	}
	var count int64
	db.Table("user_roles").Where("role_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(400, ErrorResponse{
			Error: "Cannot delete role as it is assigned to users",
		})
		return
	}
	if result := db.Delete(&role); result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Failed to delete role",
		})
		return
	}

	c.JSON(200, Response{
		Message: "Role successfully deleted",
	})
}

// @title Role Management API
// @version 1.0
// @description API for managing roles
// @host localhost:8081
// @BasePath /

package controllers

import (
	"login-api/internal/usecases"
	"login-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RoleController struct {
	roleUseCase *usecases.RoleUseCase
}

func NewRoleController(roleUseCase *usecases.RoleUseCase) *RoleController {
	return &RoleController{roleUseCase: roleUseCase}
}

// @Summary Get all roles
// @Description Get a list of all roles
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]domains.Role} "Success"
// @Failure 500 {object} ErrorResponse{error=string} "Internal Server Error"
// @Router /roles [get]
func (ctrl *RoleController) GetRoles(c *gin.Context) {
	roles, err := ctrl.roleUseCase.GetAll()
	if err != nil {
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
// @Success 200 {object} Response{data=domains.Role} "Success"
// @Router /roles/{id} [get]
func (ctrl *RoleController) GetRole(c *gin.Context) {
	id := c.Param("id")
	role, err := ctrl.roleUseCase.GetByID(id)
	if err != nil {
		c.JSON(404, ErrorResponse{
			Error: "Role not found",
		})
		return
	}

	c.JSON(200, Response{
		Data: role,
	})
}

// @Summary Create new role
// @Description Create a new role
// @Tags roles
// @Accept json
// @Produce json
// @Param role body domains.Role true "Role information"
// @Success 200 {object} Response{data=domains.Role} "Success"
// @Failure 400 {object} ErrorResponse{errors=map[string]string} "Validation Error"
// @Failure 409 {object} ErrorResponse{error=string} "Role Already Exists"
// @Failure 500 {object} ErrorResponse{error=string} "Internal Server Error"
func (ctrl *RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				switch fieldErr.Field() {
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
	}

	err := ctrl.roleUseCase.Create(&role)
	if err != nil {
		if err.Error() == "role already exists" {
			c.JSON(409, ErrorResponse{
				Error: "Role already exists",
			})
			return
		}
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
// @Param role body domains.Role true "Updated role information"
// @Success 200 {object} Response{data=domains.Role} "Success"
// @Failure 404 {object} ErrorResponse{error=string} "Role Not Found"
// @Failure 400 {object} ErrorResponse{error=string} "Invalid Data"
// @Failure 500 {object} ErrorResponse{error=string} "Update Failed"
func (ctrl *RoleController) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, ErrorResponse{
			Error: "Invalid data provided",
		})
		return
	}

	err := ctrl.roleUseCase.Update(id, &role)
	if err != nil {
		if err.Error() == "role not found" {
			c.JSON(404, ErrorResponse{
				Error: "Role not found",
			})
			return
		}
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
func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.roleUseCase.Delete(id)
	if err != nil {
		switch err.Error() {
		case "role not found":
			c.JSON(404, ErrorResponse{
				Error: "Role not found",
			})
		case "role is assigned to users":
			c.JSON(400, ErrorResponse{
				Error: "Cannot delete role as it is assigned to users",
			})
		default:
			c.JSON(500, ErrorResponse{
				Error: "Failed to delete role",
			})
		}
		return
	}

	c.JSON(200, Response{
		Message: "Role successfully deleted",
	})
}

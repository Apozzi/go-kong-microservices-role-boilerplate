// @title Item Management API
// @version 1.0
// @description API for item management
// @host localhost:8081
// @BasePath /

package controllers

import (
	"login-api/internal/usecases"
	"login-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ItemController struct {
	itemUseCase *usecases.ItemUseCase
}

func NewItemController(itemUseCase *usecases.ItemUseCase) *ItemController {
	return &ItemController{itemUseCase: itemUseCase}
}

// @Summary List all items
// @Description Get a list of all registered items
// @Tags items
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]domains.Item} "Success"
// @Failure 500 {object} ErrorResponse{error=string} "Internal Server Error"
// @Router /items [get]
func (ctrl *ItemController) GetItems(c *gin.Context) {
	items, err := ctrl.itemUseCase.GetAll()
	if err != nil {
		c.JSON(500, ErrorResponse{
			Error: "Failed to retrieve items",
		})
		return
	}

	c.JSON(200, Response{
		Data: items,
	})
}

// @Summary Get item by ID
// @Description Get a specific item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} Response{data=domains.Item} "Success"
// @Failure 404 {object} ErrorResponse{error=string} "Item not found"
// @Router /items/{id} [get]
func (ctrl *ItemController) GetItem(c *gin.Context) {
	id := c.Param("id")
	item, err := ctrl.itemUseCase.GetByID(id)
	if err != nil {
		c.JSON(404, ErrorResponse{
			Error: "Item not found",
		})
		return
	}

	c.JSON(200, Response{
		Data: item,
	})
}

// @Summary Create new item
// @Description Create a new item with the provided data
// @Tags items
// @Accept json
// @Produce json
// @Param item body domains.Item true "Item information"
// @Success 201 {object} Response{data=domains.Item} "Item created successfully"
// @Failure 400 {object} ErrorResponse{errors=map[string]string} "Validation error"
// @Failure 500 {object} ErrorResponse{error=string} "Failed to create item"
// @Router /items [post]
func (ctrl *ItemController) CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				switch fieldErr.Field() {
				case "Descricao":
					errorMessages["descricao"] = "Descrição é obrigatória"
				case "Valor":
					errorMessages["valor"] = "Valor é obrigatório"
				}
			}
			c.JSON(400, ErrorResponse{
				Errors: errorMessages,
			})
			return
		}

		c.JSON(400, ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	err := ctrl.itemUseCase.Create(&item)
	if err != nil {
		if err.Error() == "valor deve ser maior que zero" {
			c.JSON(400, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(500, ErrorResponse{
			Error: "Failed to create item",
		})
		return
	}

	c.JSON(201, Response{
		Data: item,
	})
}

// @Summary Update item
// @Description Update an existing item's information
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body domains.Item true "Updated item information"
// @Success 200 {object} Response{data=domains.Item} "Item updated successfully"
// @Failure 400 {object} ErrorResponse{error=string} "Invalid data provided"
// @Failure 404 {object} ErrorResponse{error=string} "Item not found"
// @Failure 500 {object} ErrorResponse{error=string} "Failed to update item"
// @Router /items/{id} [put]
func (ctrl *ItemController) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, ErrorResponse{
			Error: "Invalid data provided",
		})
		return
	}

	err := ctrl.itemUseCase.Update(id, &item)
	if err != nil {
		switch err.Error() {
		case "item not found":
			c.JSON(404, ErrorResponse{
				Error: "Item não encontrado",
			})
		case "valor deve ser maior que zero":
			c.JSON(400, ErrorResponse{
				Error: "Valor deve ser maior que zero",
			})
		default:
			c.JSON(500, ErrorResponse{
				Error: "Failed to update item",
			})
		}
		return
	}

	c.JSON(200, Response{
		Data: item,
	})
}

// @Summary Delete item
// @Description Remove an item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} Response{message=string} "Item deleted successfully"
// @Failure 404 {object} ErrorResponse{error=string} "Item not found"
// @Failure 500 {object} ErrorResponse{error=string} "Failed to delete item"
// @Router /items/{id} [delete]
func (ctrl *ItemController) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.itemUseCase.Delete(id)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(404, ErrorResponse{
				Error: "Item não encontrado",
			})
			return
		}

		c.JSON(500, ErrorResponse{
			Error: "Failed to delete item",
		})
		return
	}

	c.JSON(200, Response{
		Message: "Item deletado com sucesso",
	})
}

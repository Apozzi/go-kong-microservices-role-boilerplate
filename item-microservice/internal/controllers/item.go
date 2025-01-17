// @title Item Management API
// @version 1.0
// @description API for item management
// @host localhost:8081
// @BasePath /

package controllers

import (
	models "login-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary List all items
// @Description Get a list of all registered items
// @Tags items
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]models.Item} "Success"
// @Failure 500 {object} ErrorResponse{error=string} "Internal Server Error"
// @Router /items [get]
func GetItems(c *gin.Context) {
	var items []models.Item
	result := db.Find(&items)
	if result.Error != nil {
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
// @Success 200 {object} Response{data=models.Item} "Success"
// @Failure 404 {object} ErrorResponse{error=string} "Item not found"
// @Router /items/{id} [get]
func GetItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	if err := db.Where("id = ?", id).First(&item).Error; err != nil {
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
// @Param item body models.Item true "Item information"
// @Success 201 {object} Response{data=models.Item} "Item created successfully"
// @Failure 400 {object} ErrorResponse{errors=map[string]string} "Validation error"
// @Failure 500 {object} ErrorResponse{error=string} "Failed to create item"
// @Router /items [post]
func PostItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			switch fieldName {
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
	item.CriadoEm = time.Now()
	item.AtualizadoEm = time.Now()
	if item.Valor <= 0 {
		c.JSON(400, ErrorResponse{
			Error: "Valor deve ser maior que zero",
		})
		return
	}

	if result := db.Create(&item); result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Falha ao criar item",
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
// @Param item body models.Item true "Updated item information"
// @Success 200 {object} Response{data=models.Item} "Item updated successfully"
// @Failure 400 {object} ErrorResponse{error=string} "Invalid data provided"
// @Failure 404 {object} ErrorResponse{error=string} "Item not found"
// @Failure 500 {object} ErrorResponse{error=string} "Failed to update item"
// @Router /items/{id} [put]
func PutItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	if err := db.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(404, ErrorResponse{
			Error: "Item não encontrado",
		})
		return
	}
	oldItem := item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, ErrorResponse{
			Error: "Dados inválidos fornecidos",
		})
		return
	}

	if item.Valor <= 0 {
		c.JSON(400, ErrorResponse{
			Error: "Valor deve ser maior que zero",
		})
		return
	}

	item.AtualizadoEm = time.Now()
	item.CriadoEm = oldItem.CriadoEm
	item.ID = oldItem.ID

	if result := db.Save(&item); result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Falha ao atualizar item",
		})
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
func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	if err := db.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(404, ErrorResponse{
			Error: "Item não encontrado",
		})
		return
	}

	if result := db.Delete(&item); result.Error != nil {
		c.JSON(500, ErrorResponse{
			Error: "Falha ao deletar item",
		})
		return
	}

	c.JSON(200, Response{
		Message: "Item deletado com sucesso",
	})
}

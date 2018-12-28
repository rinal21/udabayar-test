package handler

import (
	"fmt"
	"net/http"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CategoryHandler struct {
	DB *gorm.DB
}

func (h *CategoryHandler) AllCategory(c *gin.Context) {
	var category []model.Category

	if err := h.DB.Find(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("Category is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   category,
	})
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	var category model.Category

	id := c.Param("id")

	if err := h.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   category,
	})
}

func NewCategoryHandler(database *gorm.DB) *CategoryHandler {
	return &CategoryHandler{
		DB: database,
	}
}

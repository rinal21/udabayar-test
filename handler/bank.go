package handler

import (
	"fmt"
	"net/http"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type BankHandler struct {
	DB *gorm.DB
}

func (h *BankHandler) AllBank(c *gin.Context) {
	var bank []model.Bank

	if err := h.DB.Preload("PaymentMethod").Find(&bank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("Bank is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   bank,
	})
}

func (h *BankHandler) GetBank(c *gin.Context) {
	var bank model.Bank

	id := c.Param("id")

	if err := h.DB.Preload("PaymentMethod").First(&bank, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   bank,
	})
}

func NewBankHandler(database *gorm.DB) *BankHandler {
	return &BankHandler{
		DB: database,
	}
}

package handler

import (
	"net/http"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CourierHandler struct {
	DB *gorm.DB
}

func (h *CourierHandler) AllCourier(c *gin.Context) {
	var couriers []model.Courier

	if check := h.DB.Find(&couriers).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Kurir tidak ditemukan.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   couriers,
	})
}

// func (h *CourierHandler) GetCourier(c *gin.Context) {
// 	var bank model.Bank

// 	id := c.Param("id")

// 	if err := h.DB.First(&bank, id).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  0,
// 			"message": fmt.Sprint(err.Error()),
// 		})

// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status": 1,
// 		"data":   bank,
// 	})
// }

func NewCourierHandler(database *gorm.DB) *CourierHandler {
	return &CourierHandler{
		DB: database,
	}
}

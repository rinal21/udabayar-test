package handler

import (
	"net/http"
	"udabayar-go-api-di/model"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

type ShopCourierHandler struct {
	DB *gorm.DB
}

func (h *ShopCourierHandler) CreateShopCourier(c *gin.Context) {
	var newShopCourier model.ShopCourier

	if err := c.ShouldBind(&newShopCourier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if err := h.DB.FirstOrCreate(&newShopCourier).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": "Berhasil.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil memasukkan review kedalam produk ini.",
	})
}

func NewShopCourierHandler(database *gorm.DB) *ShopCourierHandler {
	return &ShopCourierHandler{
		DB: database,
	}
}

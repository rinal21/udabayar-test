package handler

import (
	"fmt"
	"net/http"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type ShopHandler struct {
	DB *gorm.DB
}

func (h *ShopHandler) AllShop(c *gin.Context) {
	var shop []model.Shop

	UserID, _ := c.Get("UserID")

	if err := h.DB.Where("owner_id = ?", UserID.(uint)).Find(&shop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Toko tidak ditemukan.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   shop,
	})
}

func (h *ShopHandler) GetShop(c *gin.Context) {
	var shop model.Shop

	id := c.Param("id")

	if err := h.DB.Preload("ShopCouriers.Courier").First(&shop, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Toko tidak ditemukan.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   shop,
	})
}

func (h *ShopHandler) CreateShop(c *gin.Context) {
	var shop model.Shop

	if err := c.ShouldBind(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	UserID, _ := c.Get("UserID")
	shop.OwnerID = UserID.(uint)

	shop.SlugName = slug.Make(shop.Name)

	if err := h.DB.Omit("deleted_at").Create(&shop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal membuat toko, silakan coba kembali.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil membuat toko.",
	})
}

func (h *ShopHandler) UpdateShop(c *gin.Context) {
	var shop model.Shop
	var newShop model.Shop

	c.ShouldBind(&newShop)

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if check := h.DB.Where("owner_id = ?", UserID).First(&shop, id).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Toko tidak ditemukan.",
		})

		return
	}

	if err := h.DB.Omit("deleted_at").Model(&shop).Updates(&newShop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal mengubah data toko.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil mengubah data toko.",
	})
}

func (h *ShopHandler) GetMainAddress(c *gin.Context) {
	var shop model.Shop
	var address model.Address

	id := c.Param("id")

	if err := h.DB.First(&shop, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	if err := h.DB.Model(&shop).Related(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   address,
	})
}

func (h *ShopHandler) DeleteShop(c *gin.Context) {
	var shop model.Shop

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if err := h.DB.Where("owner_id = ?", UserID.(uint)).First(&shop, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	if err := h.DB.Delete(&shop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("Shop is not delete: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Shop deleted.",
	})
}

func NewShopHandler(database *gorm.DB) *ShopHandler {
	return &ShopHandler{
		DB: database,
	}
}

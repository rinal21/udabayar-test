package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type StorefrontHandler struct {
	DB *gorm.DB
}

func (h *StorefrontHandler) AllStorefront(c *gin.Context) {
	var storefront []model.Storefront

	UserID, _ := c.Get("UserID")

	if err := h.DB.Where("owner_id = ?", UserID.(uint)).Find(&storefront).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("Storefront is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   storefront,
	})
}

func (h *StorefrontHandler) GetStorefront(c *gin.Context) {
	var storefront model.Storefront

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if err := h.DB.Where("owner_id = ?", UserID.(uint)).First(&storefront, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   storefront,
	})
}

func (h *StorefrontHandler) GetProductsByStorefront(c *gin.Context) {
	var productsStorefront []model.ProductStorefront

	storefront_id := c.Param("id")

	if err := h.DB.Preload("Storefront").Preload("Product").Find(&productsStorefront, "storefront_id = ?", storefront_id).RecordNotFound(); err == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Tidak ada produk didalam etalase ini.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   productsStorefront,
	})
}

func (h *StorefrontHandler) SetProductStorefront(c *gin.Context) {
	var productStorefront model.ProductStorefront
	var newProductsStorefront model.ProductStorefront

	if err := c.Bind(&newProductsStorefront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	storefront_id, _ := strconv.Atoi(c.Param("id"))
	productStorefront.StorefrontID = storefront_id

	if err := h.DB.Preload("Storefront").Where("product_id = ?", newProductsStorefront.ProductID).Where("storefront_id = ?", newProductsStorefront.StorefrontID).First(&productStorefront).RecordNotFound(); err != true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Produk ini sudah terdaftar didalam etalase ini.",
		})

		return
	}

	UserID, _ := c.Get("UserID")

	if productStorefront.Storefront.OwnerID == UserID {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Etalase ini bukan milik kamu.",
		})

		return
	}

	if err := h.DB.Create(&newProductsStorefront).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal memasukan produk kedalam etalase.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil memasukkan produk kedalam etalase.",
	})
}

func (h *StorefrontHandler) CreateStorefront(c *gin.Context) {
	var storefront model.Storefront

	if err := c.ShouldBind(&storefront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})
	}

	UserID, _ := c.Get("UserID")
	storefront.OwnerID = UserID.(uint)

	if err := h.DB.Omit("deleted_at").Create(&storefront).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("Storefront is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Storefront created.",
	})
}

func (h *StorefrontHandler) UpdateStorefront(c *gin.Context) {
	var storefront model.Storefront
	var newStorefront model.Storefront

	if err := c.Bind(&newStorefront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})

		return
	}

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if check := h.DB.Where("owner_id = ?", UserID.(uint)).First(&storefront, id).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Etalase tidak ditemukan.",
		})

		return
	}

	if err := h.DB.Omit("deleted_at").Model(&storefront).Updates(&newStorefront).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal mengubah data etalase.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Etalase berhasil diubah.",
	})
}

func (h *StorefrontHandler) DeleteStorefront(c *gin.Context) {
	var storefront model.Storefront

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if err := h.DB.Where("owner_id = ?", UserID.(uint)).First(&storefront, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	if err := h.DB.Delete(&storefront).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("Storefront is not delete: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Storefront deleted.",
	})
}

func NewStorefrontHandler(database *gorm.DB) *StorefrontHandler {
	return &StorefrontHandler{
		DB: database,
	}
}

package handler

import (
	"fmt"
	"net/http"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ProductVariantHandler struct {
	DB *gorm.DB
}

func (h *ProductVariantHandler) AllProductVariant(c *gin.Context) {
	var productVariant []model.ProductVariant

	id := c.Param("id")

	if err := h.DB.Where("product_id = ?", id).Find(&productVariant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("ProductVariant is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   productVariant,
	})
}

func (h *ProductVariantHandler) GetProductVariant(c *gin.Context) {
	var productVariant model.ProductVariant

	variant_id := c.Param("variant_id")
	product_id := c.Param("id")

	if err := h.DB.Where("product_id = ?", product_id).First(&productVariant, variant_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   productVariant,
	})
}

func (h *ProductVariantHandler) CreateProductVariant(c *gin.Context) {
	var productVariant model.ProductVariant

	if err := c.ShouldBind(&productVariant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})
	}

	if err := h.DB.Omit("deleted_at").Create(&productVariant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("ProductVariant is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "ProductVariant created.",
	})
}

func (h *ProductVariantHandler) UpdateProductVariant(c *gin.Context) {
	var productVariant model.ProductVariant

	id := c.Param("id")

	if err := h.DB.First(&productVariant, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	updateProductVariant := productVariant

	if err := c.Bind(&updateProductVariant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})

		return
	}

	if err := h.DB.Omit("deleted_at").Model(&productVariant).Updates(&updateProductVariant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("ProductVariant is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "ProductVariant updated.",
	})
}

func (h *ProductVariantHandler) DeleteProductVariant(c *gin.Context) {
	var productVariant model.ProductVariant

	product_id := c.Param("id")
	variant_id := c.Param("variant_id")

	if err := h.DB.Where("product_id = ?", product_id).First(&productVariant, variant_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	if err := h.DB.Delete(&productVariant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("ProductVariant is not delete: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "ProductVariant deleted.",
	})
}

func NewProductVariantHandler(database *gorm.DB) *ProductVariantHandler {
	return &ProductVariantHandler{
		DB: database,
	}
}

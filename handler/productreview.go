package handler

import (
	"net/http"
	"strconv"
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/model"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

type ProductReviewHandler struct {
	DB       *gorm.DB
	TxConfig *config.TransactionConfig
}

func (h *ProductReviewHandler) CreateProductReview(c *gin.Context) {
	var newProductReview model.ProductReview
	var productReview model.ProductReview
	var transaction model.Transaction

	id := c.Param("id")
	product_id, _ := strconv.Atoi(id)

	if err := c.ShouldBind(&newProductReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if check := h.DB.Preload("Order.Product").First(&transaction, newProductReview.TransactionID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Transaksi tidak ditemukan.",
		})

		return
	}

	if transaction.Status != h.TxConfig.Status.Received {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Transaksi tidak dalam status diterima.",
		})

		return
	}

	if product_id != transaction.Order.ProductID {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Tidak ada produk didalam transaksi ini.",
		})

		return
	}

	if check := h.DB.First(&productReview, "transaction_id = ?", newProductReview.TransactionID).RecordNotFound(); check != true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Transaksi ini sudah terdapat review.",
		})

		return
	}

	newProductReview.ProductID = product_id

	if err := h.DB.Create(&newProductReview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal memasukkan review kedalam produk ini, silakan coba kembali.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil memasukkan review kedalam produk ini.",
	})
}

func NewProductReviewHandler(database *gorm.DB, txConfig *config.TransactionConfig) *ProductReviewHandler {
	return &ProductReviewHandler{
		DB:       database,
		TxConfig: txConfig,
	}
}

package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/helper"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type ProductHandler struct {
	DB     *gorm.DB
	Upload *helper.UploadHelper
	Config *config.Config
}

func (h *ProductHandler) AllProduct(c *gin.Context) {
	var products []model.Product

	shopID := c.Param("shop_id")

	if shopID == "" {
		if check := h.DB.Preload("ProductReview").Preload("ProductImage").Preload("ProductDetail.ProductVariant").Preload("Shop").Find(&products).RecordNotFound(); check == true {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Produk tidak ada.",
			})

			return
		}
	} else {
		if check := h.DB.Preload("ProductReview").Preload("ProductImage").Preload("ProductDetail.ProductVariant").Preload("Shop").Find(&products, "shop_id = ?", shopID).RecordNotFound(); check == true {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Produk tidak ada.",
			})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   products,
	})
}

type ProductResponse struct {
	model.Product
	UserBank []UserBankResponse    `json:"userbank" binding:"-"`
	Rating   float32               `json:"rating" binding:"-"`
	Reviewer int                   `json:"reviewer" binding:"-"`
	Sold     []ProductSoldResponse `json:"sold" binding:"-"`
}

type UserBankResponse struct {
	BankID     int    `json:"bank_id"`
	UserBankID uint   `json:"userbank_id"`
	BankCode   string `json:"bank_code"`
	BankName   string `json:"bank_name"`
}

type ProductSoldResponse struct {
	ProductDetailID uint `json:"product_detail_id"`
	Sold            int  `json:"sold"`
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	var product model.Product
	var userBanks []model.UserBank

	var productResponse ProductResponse

	id := c.Param("id")
	shop_slug := c.Param("shop")
	product_slug := c.Param("product")

	if id != "" {
		if check := h.DB.Preload("Shop").First(&product, id).RecordNotFound(); check == true {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Produk tidak ditemukan.",
			})

			return
		}
	} else {
		if check := h.DB.Preload("Orders").Preload("ProductReview").Preload("ProductImage").Preload("ProductDetail.ProductVariant").Preload("Shop").Find(&product, "slug_name = ?", product_slug).RecordNotFound(); check == true {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Produk tidak ditemukan.",
			})

			return
		}

		if product.Shop.SlugName != shop_slug {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Toko tidak mempunyai produk tersebut.",
			})

			return
		}
	}

	// if err := h.DB.Model(&product).Related(&shop).Related(&productVariant).Related(&productImages).Related(&productReview).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":  0,
	// 		"message": err.Error(),
	// 	})

	// 	return
	// }

	if check := h.DB.Preload("Bank").Find(&userBanks, "owner_id = ?", product.Shop.OwnerID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Tidak ada userbank",
		})

		return
	}

	var bankID []int

	for _, element := range userBanks {
		bankID = append(bankID, element.BankID)
	}

	var UserBankResponseData []UserBankResponse

	for _, element := range userBanks {
		UserBankResponseData = append(UserBankResponseData, UserBankResponse{
			BankID:     element.BankID,
			UserBankID: element.ID,
			BankCode:   element.Bank.Code,
			BankName:   element.Bank.Name,
		})
	}

	go h.DB.Model(&product).Update("view_count", gorm.Expr("view_count + ?", 1))

	var ratingTotal float32
	ratingPerson := len(product.ProductReview)

	for _, element := range product.ProductReview {
		ratingTotal = ratingTotal + float32(element.Rating)
	}

	rating := ratingTotal / float32(ratingPerson)

	if math.IsNaN(float64(rating)) {
		rating = 0
	}

	var solds []ProductSoldResponse

	for _, element := range product.ProductDetail {
		var orders []model.Order
		var sold int

		if err := h.DB.Find(&orders, "product_detail_id = ?", element.ID).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"data":   err.Error(),
			})

			return
		}

		for _, order := range orders {
			var orderTransaction model.Transaction

			if err := h.DB.First(&orderTransaction, order.TransactionID).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": 0,
					"data":   err.Error(),
				})

				return
			}

			if orderTransaction.Status == 4 {
				sold = sold + order.Qty
			}
		}

		solds = append(solds, ProductSoldResponse{
			ProductDetailID: element.ID,
			Sold:            sold,
		})
	}

	productResponse.Product = product

	productResponse.ViewCount = product.ViewCount + 1
	productResponse.Rating = rating
	productResponse.Sold = solds
	productResponse.Reviewer = ratingPerson
	productResponse.UserBank = UserBankResponseData

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   productResponse,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product model.Product

	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	product.SlugName = slug.Make(product.Name)

	if err := h.DB.Omit("deleted_at").Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Product created.",
	})
}

func (h *ProductHandler) UploadProductImage(c *gin.Context) {
	var product model.Product
	product_id := c.Param("id")

	if err := h.DB.First(&product, product_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	formFile, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	inNow := time.Now().Unix()
	format := fmt.Sprint(product_id, formFile.Filename, inNow)

	hash := md5.New()
	hash.Write([]byte(format))
	hashed := hex.EncodeToString(hash.Sum(nil))
	newFileName := fmt.Sprintf("%s.jpg", hashed)

	path := fmt.Sprintf("/images/products/%s", newFileName)

	go h.Upload.UploadImage(formFile, path)

	imageUrl := fmt.Sprintf("https://s3.amazonaws.com/%s%s", h.Config.S3Bucket, path)

	ProductID, err := strconv.Atoi(product_id)

	newProductImage := &model.ProductImage{
		ProductID: ProductID,
		ImageUrl:  imageUrl,
	}

	if err := h.DB.Omit("deleted_at").Create(&newProductImage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Upload Product Image success.",
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var product model.Product

	id := c.Param("id")

	if err := h.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	updateProduct := product

	if err := c.Bind(&updateProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})

		return
	}

	if err := h.DB.Omit("deleted_at", "shop").Model(&product).Updates(&updateProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("Product is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Product updated.",
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	var product model.Product

	id := c.Param("id")

	if err := h.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if err := h.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("Product is not delete: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Product deleted.",
	})
}

func NewProductHandler(database *gorm.DB, upload *helper.UploadHelper, config *config.Config) *ProductHandler {
	return &ProductHandler{
		DB:     database,
		Upload: upload,
		Config: config,
	}
}

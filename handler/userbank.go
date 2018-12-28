package handler

import (
	"fmt"
	"net/http"
	"udabayar-go-api-di/helper"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserBankHandler struct {
	DB     *gorm.DB
	Bcrypt *helper.Bcrypt
}

type UserBankShow struct {
	model.UserBank
	BankCode string `json:"bank_code"`
}

func (h *UserBankHandler) AllUserBank(c *gin.Context) {
	var userbank []model.UserBank

	userID, _ := c.Get("UserID")

	if userID == nil {
		if err := h.DB.Preload("Bank").Find(&userbank).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": err.Error(),
			})

			return
		}
	} else {
		if err := h.DB.Preload("Bank").Find(&userbank, "owner_id = ?", userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": err.Error(),
			})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   userbank,
	})
}

func (h *UserBankHandler) GetUserBank(c *gin.Context) {
	var userBank model.UserBank

	id := c.Param("id")
	userID, _ := c.Get("UserID")

	if check := h.DB.Preload("Bank").Where("owner_id = ?", userID).First(&userBank, id).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Akun bank tidak ditemukan.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   userBank,
	})
}

type UserBankByShopIDResponse struct {
	ID       uint   `json:"userbank_id"`
	BankID   uint   `json:"bank_id"`
	BankName string `json:"bank_name"`
	BankCode string `json:"bank_code"`
}

func (h *UserBankHandler) GetUserBankByShopID(c *gin.Context) {
	var shop model.Shop
	var userBank []model.UserBank
	var response []UserBankByShopIDResponse

	shopID := c.Param("shop_id")

	if check := h.DB.First(&shop, shopID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Toko tidak ditemukan.",
		})

		return
	}

	if check := h.DB.Preload("Bank").Find(&userBank, "owner_id = ?", shop.OwnerID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "User tidak mempunyai bank",
		})

		return
	}

	for _, element := range userBank {
		response = append(response, UserBankByShopIDResponse{
			ID:       element.ID,
			BankID:   element.Bank.ID,
			BankCode: element.Bank.Code,
			BankName: element.Bank.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   userBank,
	})
}

func (h *UserBankHandler) CreateUserBank(c *gin.Context) {
	var userBank model.UserBank

	if err := c.ShouldBind(&userBank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	// hash, _ := h.Bcrypt.HashPassword(userBank.Password)
	// userBank.Password = hash

	UserID, _ := c.Get("UserID")
	userBank.OwnerID = UserID.(uint)

	if check := h.DB.Find(&userBank, "account_number = ?", userBank.AccountNumber).RecordNotFound(); check != true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Nomor Rekening sudah tersedia.",
		})

		return
	}

	if err := h.DB.Omit("deleted_at").Create(&userBank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal menyimpan akun bank.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil menyimpan akun bank.",
	})
}

func (h *UserBankHandler) UpdateUserBank(c *gin.Context) {
	var userBank model.UserBank

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if err := h.DB.Where("owner_id = ?", UserID.(uint)).First(&userBank, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Akun bank tidak ditemukan.",
		})

		return
	}

	newUserBank := userBank

	if err := c.Bind(&newUserBank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	// if check := newUserBank.Password; check != "" {
	// 	hash, _ := h.Bcrypt.HashPassword(newUserBank.Password)
	// 	newUserBank.Password = hash
	// }

	if check := newUserBank.AccountNumber; check != "" {
		if check := h.DB.First(&userBank, "account_number = ?", newUserBank.AccountNumber).RecordNotFound(); check != true {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Nomor Rekening sudah tersedia.",
			})

			return
		}
	}

	if err := h.DB.Omit("deleted_at", "bank").Model(&userBank).Updates(&newUserBank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("UserBank is not created: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "UserBank updated.",
	})
}

func (h *UserBankHandler) DeleteUserBank(c *gin.Context) {
	var userBank model.UserBank

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if err := h.DB.Where("owner_id = ?", UserID.(uint)).First(&userBank, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	if err := h.DB.Delete(&userBank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint("UserBank is not delete: ", err.Error()),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "UserBank deleted.",
	})
}

type WebhookUserBankInput struct {
	UserBankID int `form:"userbank_id" json:"userbank_id" binding:"required"`
	Balance    int `form:"balance" json:"balance" binding:"required"`
}

func (h *UserBankHandler) WebhookUserBank(c *gin.Context) {
	var userBank model.UserBank
	var input WebhookUserBankInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if check := h.DB.First(&userBank, input.UserBankID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Akun bank tidak ditemukan.",
		})

		return
	}

	if err := h.DB.Omit("deleted_at", "bank").Model(&userBank).Update("balance", input.Balance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal mengubah saldo akun bank.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil mengubah saldo akun bank.",
	})
}

func NewUserBankHandler(database *gorm.DB, bcrypt *helper.Bcrypt) *UserBankHandler {
	return &UserBankHandler{
		DB:     database,
		Bcrypt: bcrypt,
	}
}

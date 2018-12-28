package handler

import (
	"net/http"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AddressHandler struct {
	DB *gorm.DB
}

func (h *AddressHandler) AllAddresses(c *gin.Context) {
	var addresses []model.Address

	UserID, _ := c.Get("UserID")

	if check := h.DB.Where("owner_id = ?", UserID).Find(&addresses).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Tidak ada alamat yang terdaftar.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   addresses,
	})
}

func (h *AddressHandler) GetAddress(c *gin.Context) {
	var address model.Address

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if check := h.DB.Where("owner_id = ?", UserID).First(&address, id).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Alamat tidak ditemukan.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   address,
	})
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var address model.Address

	if err := c.ShouldBind(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	UserID, _ := c.Get("UserID")
	address.OwnerID = UserID.(uint)

	if err := h.DB.Omit("deleted_at").Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal menambahkan alamat, silakan coba kembali.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil menambahkan alamat.",
	})
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	var address model.Address
	var newAddress model.Address

	c.ShouldBind(&newAddress)

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if check := h.DB.Where("owner_id = ?", UserID.(uint)).First(&address, id).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Alamat tidak ditemukan.",
		})

		return
	}

	if err := h.DB.Omit("deleted_at").Model(&address).Updates(&newAddress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal mengubah data alamat, silakan coba kembali.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil mengubah data alamat.",
	})
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	var address model.Address

	id := c.Param("id")
	UserID, _ := c.Get("UserID")

	if check := h.DB.Where("owner_id = ?", UserID.(uint)).First(&address, id).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Alamat tidak ditemukan.",
		})

		return
	}

	if err := h.DB.Delete(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal menghapus alamat, silakan coba kembali.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil menghapus alamat.",
	})
}

func NewAddressHandler(database *gorm.DB) *AddressHandler {
	return &AddressHandler{
		DB: database,
	}
}

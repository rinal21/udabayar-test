package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/helper"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
)

type ProfileHandler struct {
	Config *config.Config
	DB     *gorm.DB
	Upload *helper.UploadHelper
}

type UserProfile struct {
	Email string `form:"email" json:"email" binding:"email"`
	Name  string `form:"name" json:"name" binding:"-"`
	Phone string `form:"phone" json:"phone" binding:"-"`
	Image string `form:"image" json:"image" binding:"-"`
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	var user model.User

	UserID, _ := c.Get("UserID")

	if check := h.DB.Preload("Profile").First(&user, UserID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "User tidak ditemukan.",
		})

		return
	}

	profile := &UserProfile{
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Profile.Phone,
		Image: user.Profile.Image,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   profile,
	})

	return
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var user model.User
	var profile UserProfile

	if err := c.ShouldBind(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	UserID, _ := c.Get("UserID")

	if check := h.DB.Preload("Profile").First(&user, UserID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "User tidak ditemukan.",
		})

		return
	}

	if err := h.DB.Model(&user).Updates(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil mengubah data profil.",
	})

	return
}

type ProfileImage struct {
	Image string `form:"image" json:"image"`
}

func (h *ProfileHandler) UpdateProfileImage(c *gin.Context) {
	var profile model.Profile
	var profileImage ProfileImage

	UserID, _ := c.Get("UserID")

	if check := h.DB.First(&profile, "user_id = ?", UserID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "User tidak ditemukan.",
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
	format := fmt.Sprint(UserID, formFile.Filename, inNow)

	hash := md5.New()
	hash.Write([]byte(format))
	hashed := hex.EncodeToString(hash.Sum(nil))
	newFileName := fmt.Sprintf("%s.jpg", hashed)

	path := fmt.Sprintf("/images/users/%s", newFileName)

	go h.Upload.UploadImage(formFile, path)

	imageUrl := fmt.Sprintf("https://s3.amazonaws.com/%s%s", h.Config.S3Bucket, path)
	profileImage.Image = imageUrl

	if err := h.DB.Model(&profile).Updates(&profileImage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil mengubah foto profil.",
	})

	return
}

func NewProfileHandler(config *config.Config, database *gorm.DB, upload *helper.UploadHelper) *ProfileHandler {
	return &ProfileHandler{
		Config: config,
		DB:     database,
		Upload: upload,
	}
}

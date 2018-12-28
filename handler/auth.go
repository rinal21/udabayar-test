package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/helper"
	"udabayar-go-api-di/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AuthHandler struct {
	Config *config.Config
	DB     *gorm.DB
	JWT    *helper.JWT
	Bcrypt *helper.Bcrypt
	Email  *helper.Email
}

func NewAuthHandler(
	config *config.Config,
	database *gorm.DB,
	jwt *helper.JWT,
	bcrypt *helper.Bcrypt,
	email *helper.Email,
) *AuthHandler {
	return &AuthHandler{
		Config: config,
		DB:     database,
		JWT:    jwt,
		Bcrypt: bcrypt,
		Email:  email,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var auth model.User

	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if check := h.DB.Where("email = ?", auth.Email).First(&auth).RecordNotFound(); check != true {
		c.JSON(http.StatusConflict, gin.H{
			"status":  0,
			"message": "Email sudah digunakan.",
		})

		return
	}

	hash, _ := h.Bcrypt.HashPassword(auth.Password)
	auth.Password = hash

	if err := h.DB.Omit("deleted_at").Create(&auth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal melakukan pendaftaran, silakan coba kembali.",
		})

		return
	}

	token, _ := h.JWT.SignJWT(auth.ID, "emailconfirmation")
	link := fmt.Sprint(h.Config.BaseUrl, "activation?token=", token)

	templateData := &helper.RegistrationTemplateData{
		Fullname:       auth.Name,
		ActivationLink: link,
		ExpiredAt:      time.Now().Local().Add(time.Hour * 24),
	}

	go h.Email.SendEmailRegistration(templateData, auth.Email)

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil melakukan pendaftaran.",
		"data": gin.H{
			"token": token,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var auth model.User

	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})

		return
	}

	password := auth.Password

	if check := h.DB.Find(&auth, "email = ?", auth.Email).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Email or Password not valid.",
		})

		return
	}

	if check := h.Bcrypt.CheckPasswordHash(password, auth.Password); check != true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Email or Password not match.",
		})

		return
	}

	token, _ := h.JWT.SignJWT(auth.ID, "access")

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data": gin.H{
			"token": token,
		},
	})
}

func (h *AuthHandler) CheckActivatedAccount(c *gin.Context) {
	var auth model.User
	UserID, _ := c.Get("UserID")

	if err := h.DB.First(&auth, UserID).RecordNotFound(); err == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Akun tidak tersedia.",
		})

		return
	}

	if auth.Status == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Harap aktivasi email terlebih dahulu sebelum menggunakan fitur ini.",
		})

		c.Abort()
		return
	}

	c.Next()
}

type ValidateTokenForm struct {
	Token string `form:"token" json:"token" binding:"required"`
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var tokenForm ValidateTokenForm

	if err := c.ShouldBind(&tokenForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})

		return
	}

	claims := &helper.Claims{}
	parseJWT, err := h.JWT.ParseJWT(tokenForm.Token, claims)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if !parseJWT.Valid {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  0,
			"message": "Token is not valid.",
		})

		return
	}

	tokenScope := claims.Scope
	scopeSplit := strings.Split(tokenScope, ":")

	if scopeSplit[1] != "access" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  0,
			"message": "Token scope not valid for access api",
		})

		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Token is valid.",
	})
}

func (h *AuthHandler) Authorizer(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")

	if tokenHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  0,
			"message": "Unauthorized.",
		})

		c.Abort()
		return
	}

	tokenSplit := strings.Split(tokenHeader, " ")

	if len(tokenSplit) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  0,
			"message": "Invalid/Malformed authentication token.",
		})

		c.Abort()
		return
	}

	token := tokenSplit[1]
	claims := &helper.Claims{}

	parseJWT, err := h.JWT.ParseJWT(token, claims)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		c.Abort()
		return
	}

	if !parseJWT.Valid {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  0,
			"message": "Token is not valid.",
		})

		c.Abort()
		return
	}

	c.Set("UserID", claims.UserId)

	tokenScope := claims.Scope
	scopeSplit := strings.Split(tokenScope, ":")

	if scopeSplit[1] != "access" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  0,
			"message": "Token scope not valid for access api",
		})

		c.Abort()
		return
	}

	c.Next()
}

func (h *AuthHandler) ApiKey(c *gin.Context) {
	tokenHeader := c.GetHeader("x-api-key")

	if tokenHeader != h.Config.APIKey {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  0,
			"message": "Invalid API Key.",
		})

		c.Abort()
		return
	}

	c.Next()
}

type ChangePasswordForm struct {
	CurrentPassword string `form:"current_password" json:"current_password" binding:"required"`
	NewPassword     string `form:"new_password" json:"new_password" binding:"required"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var auth model.User
	var newAuth model.User
	var changePassword ChangePasswordForm

	if err := c.ShouldBind(&changePassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})

		return
	}

	UserID, _ := c.Get("UserID")

	if err := h.DB.First(&auth, UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": fmt.Sprint(err.Error()),
		})

		return
	}

	check := h.Bcrypt.CheckPasswordHash(changePassword.CurrentPassword, auth.Password)

	if check != true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Current Password not valid.",
		})

		return
	}

	hash, _ := h.Bcrypt.HashPassword(changePassword.NewPassword)
	newAuth.Password = hash

	if err := h.DB.Omit("deleted_at").Model(&auth).Updates(&newAuth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Update Password failed.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Update Password success.",
	})
}

func (h *AuthHandler) ActivateEmail(c *gin.Context) {
	var auth model.User

	token := c.Query("token")
	claims := &helper.Claims{}

	check, err := h.JWT.ParseJWT(token, claims)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if !check.Valid {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Token is not valid.",
		})

		return
	}

	tokenScope := claims.Scope
	scopeSplit := strings.Split(tokenScope, ":")

	if scopeSplit[1] != "emailconfirmation" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  0,
			"message": "Token scope not valid for email confirmation api",
		})

		return
	}

	if err := h.DB.First(&auth, claims.UserId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if err := h.DB.Model(&auth).Update("status", 1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Email confirmation failed, please try again.",
		})

		return
	}

	templateData := helper.ActivatedAccountTemplateData{
		Fullname:  auth.Name,
		LoginLink: "http://dev.udabayar.com/",
	}

	h.Email.SendEmailActivatedAccount(&templateData, auth.Email)

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Email confirmation success.",
	})
}

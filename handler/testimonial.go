package handler

import (
	"net/http"
	"strconv"
	"udabayar-go-api-di/model"

	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type TestimonialHandler struct {
	DB *gorm.DB
}

type TestimonialResponse struct {
	Name         string `form:"name" json:"name"`
	ImageProfile string `form:"image_profile" json:"image_profile"`
	ShopName     string `form:"shop_name" json:"shop_name"`
	Testimonial  string `form:"testimonial" json:"testimonial"`
}

func (h *TestimonialHandler) AllTestimonials(c *gin.Context) {
	var testimonials []model.Testimonial
	var response []TestimonialResponse

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	order := c.DefaultQuery("order", "asc")

	db := h.DB.Preload("User.Profile").Preload("Shop")

	pagination.Pagging(&pagination.Param{
		DB:      db,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id " + order},
	}, &testimonials)

	for _, testimonial := range testimonials {
		response = append(response, TestimonialResponse{
			Name:         testimonial.User.Name,
			ImageProfile: testimonial.User.Profile.Image,
			ShopName:     testimonial.Shop.Name,
			Testimonial:  testimonial.Testimonial,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   response,
	})
}

func (h *TestimonialHandler) GetTestimonial(c *gin.Context) {
	var testimonial model.Testimonial
	var response TestimonialResponse

	id := c.Param("id")

	if check := h.DB.Preload("User.Profile").Preload("Shop").First(&testimonial, id).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  1,
			"message": "Testimonial tidak ditemukan.",
		})
	}

	response = TestimonialResponse{
		Name:         testimonial.User.Name,
		ImageProfile: testimonial.User.Profile.Image,
		ShopName:     testimonial.Shop.Name,
		Testimonial:  testimonial.Testimonial,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   response,
	})
}

func (h *TestimonialHandler) CreateTestimonial(c *gin.Context) {
	var testimonial model.Testimonial

	if err := c.ShouldBind(&testimonial); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})
	}

	userId, _ := c.Get("UserID")
	testimonial.UserID = userId.(uint)

	if err := h.DB.FirstOrCreate(&testimonial).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  1,
			"message": "Gagal membuat testimonial.",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil memasukkan testimonial",
	})
}

func NewTestimonialHandler(database *gorm.DB) *TestimonialHandler {
	return &TestimonialHandler{
		DB: database,
	}
}

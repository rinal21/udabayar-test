package model

import (
	"github.com/jinzhu/gorm"
)

type Testimonial struct {
	gorm.Model
	User        User   `form:"user" json:"user" binding:"-"`
	UserID      uint   `form:"user_id" json:"user_id"`
	Shop        Shop   `form:"-" json:"-" binding:"-"`
	ShopID      uint   `form:"shop_id" json:"shop_id" binding:"required"`
	Testimonial string `form:"testimonial" json:"testimonial" binding:"required"`
}

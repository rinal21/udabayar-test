package model

import (
	"github.com/jinzhu/gorm"
)

type ProductVariant struct {
	gorm.Model
	Sku  string `form:"sku" json:"sku" binding:"required"`
	Name string `form:"name" json:"name"`
}

package model

import (
	"github.com/jinzhu/gorm"
)

type Storefront struct {
	gorm.Model
	OwnerID uint   `form:"owner_id" json:"owner_id" binding:"required`
	Name    string `form:"name" json:"name" binding:"required"`
}

type ProductStorefront struct {
	gorm.Model
	Storefront   Storefront `form:"storefront" json:"storefront" binding:"-"`
	StorefrontID int        `form:"storefront_id" json:"storefront_id" binding:"-"`
	Product      Product    `form:"product" json:"product" binding:"-"`
	ProductID    int        `form:"product_id" json:"product_id" binding:"required"`
}

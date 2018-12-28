package model

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Shop          Shop            `json:"shop" binding:"-"`
	ShopID        int             `form:"shop_id" json:"shop_id" binding:"required"`
	CategoryID    int             `form:"category_id" json:"category_id" binding:"required"`
	Name          string          `form:"name" json:"name" binding:"required"`
	SlugName      string          `form:"slug" json:"slug"`
	Description   string          `form:"description" json:"description" binding:"required" sql:"type:longtext;"`
	NonPyshical   int             `form:"non_pyshical" json:"non_pyshical"`
	Condition     int             `form:"condition" json:"condition"`
	MainSku       string          `form:"main_sku" json:"main_sku" binding:"required"`
	Weight        int             `form:"weight" json:"weight" binding:"required"`
	ViewCount     int             `json:"view_count"`
	IsActive      int             `form:"is_active" json:"is_active" binding:"required"`
	ProductDetail []ProductDetail `form:"product_detail" json:"product_detail"`
	ProductImage  []ProductImage  `json:"product_images" binding:"-"`
	ProductReview []ProductReview `json:"product_review" binding:"-"`
	Orders        []Order         `json:"orders" binding:"-"`
}

type ProductDetail struct {
	gorm.Model
	Product          Product        `form:"product" json:"product"`
	ProductID        int            `form:"product_id" json:"product_id"`
	ProductVariant   ProductVariant `form:"product_variant" json:"product_variant" binding:"-"`
	ProductVariantID int            `form:"product_variant_id" json:"product_variant_id"`
	Qty              int            `form:"qty" json:"qty" binding:"required"`
	Price            int            `form:"price" json:"price" binding:"required"`
}

type ProductImage struct {
	gorm.Model
	ProductID int    `form:"product_id" json:"product_id"`
	ImageUrl  string `form:"image_url" json:"image_url"`
}

type ProductReview struct {
	gorm.Model
	ProductID     int    `form:"product_id" json:"product_id"`
	TransactionID int    `form:"transaction_id" json:"transaction_id" binding:"required"`
	Rating        int    `form:"rating" json:"rating" binding:"required"`
	Review        string `form:"review" json:"review"`
}

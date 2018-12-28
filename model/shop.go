package model

import (
	"github.com/jinzhu/gorm"
)

type Shop struct {
	gorm.Model
	OwnerID           uint          `form:"owner_id" json:"owner_id" binding:"required`
	Name              string        `form:"name" json:"name" binding:"required"`
	SlugName          string        `form:"slug_name" json:"slug_name"`
	Description       string        `form:"description" json:"description" binding:"required,min=30"`
	ProfilePicture    string        `form:"profile_picture" json:"profile_picture" binding:"required,url"`
	BackgroundPicture string        `form:"background_picture" json:"background_picture" binding:"required,url"`
	VideoUrl          string        `form:"video_url" json:"video_url" binding:"required,url"`
	AddressID         uint          `form:"main_address_id" json:"main_address_id" binding:"required"`
	ShopCouriers      []ShopCourier `form:"shop_couriers" json:"shop_couriers" binding:"-"`
}

type ShopCourier struct {
	gorm.Model
	ShopID    int     `form:"shop_id" json:"shop_id"`
	Courier   Courier `form:"courier" json:"courier"`
	CourierID int     `form:"courier_id" json:"courier_id"`
}

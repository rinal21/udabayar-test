package model

import (
	"github.com/jinzhu/gorm"
)

type Address struct {
	gorm.Model
	Name           string `form:"name" json:"name" binding:"required"`
	Phone          string `form:"phone" json:"phone" binding:"required"`
	Address        string `form:"address" json:"address" binding:"required"`
	Province       string `form:"province" json:"province" binding:"required"`
	CityOrDistrict string `form:"city_or_district" json:"city_or_district" binding:"required"`
	SubDistrict    string `form:"sub_district" json:"sub_district" binding:"required"`
	OriginCode     int    `form:"origin_code" json:"origin_code" binding:"required"`
	PostalCode     int    `form:"postal_code" json:"postal_code" binding:"required"`
	OwnerID        uint   `form:"owner_id" json:"owner_id" binding:"required`
}

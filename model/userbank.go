package model

import (
	"github.com/jinzhu/gorm"
)

type UserBank struct {
	gorm.Model
	OwnerID        uint   `form:"owner_id" json:"owner_id"`
	Bank           Bank   `json:"bank" binding:"-"`
	BankID         int    `form:"bank_id" json:"bank_id" binding:"required"`
	AccountName    string `form:"account_name" json:"account_name" binding:"required"`
	AccountNumber  string `form:"account_number" json:"account_number" binding:"required"`
	Username       string `form:"username" json:"username" binding:"required"`
	Password       string `form:"password" json:"password" binding:"required"`
	Balance        int    `form:"balance" json:"balance"`
	Branch         string `form:"branch" json:"branch" binding:"required"`
	CityOrDistrict string `form:"city_or_district" json:"city_or_district" binding:"required"`
}

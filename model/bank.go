package model

import (
	"github.com/jinzhu/gorm"
)

type Bank struct {
	gorm.Model
	Name          string          `form:"name" json:"name"`
	Code          string          `form:"code" json:"code"`
	Logo          string          `form:"logo" json:"logo"`
	PaymentMethod []PaymentMethod `form:"payment_method" json:"payment_method"`
}

type PaymentMethod struct {
	gorm.Model
	BankID  int    `form:"bank_id" json:"bank_id"`
	Name    string `form:"name" json:"name"`
	Context string `form:"context" json:"context" sql:"type:longtext;"`
}

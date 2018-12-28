package model

import (
	"github.com/jinzhu/gorm"
)

type Courier struct {
	gorm.Model
	Name string `form:"" json:"name" binding:"required"`
	Code string `form:"" json:"code" binding:"required"`
}

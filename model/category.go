package model

import (
	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	Name string `form:"name" json:"name" binding:"required"`
}

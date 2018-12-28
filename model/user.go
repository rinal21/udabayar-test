package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email    string  `form:"email" json:"email" binding:"required,email"`
	Password string  `form:"password" json:"password" binding:"required,min=6"`
	Name     string  `form:"name" json:"name" binding:"-"`
	Status   int     `form:"status" json:"status"`
	Profile  Profile `form:"profile" json:"profile" binding:"-"`
}

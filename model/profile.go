package model

import (
	"github.com/jinzhu/gorm"
)

type Profile struct {
	gorm.Model
	UserID int    `form:"user_id" json:"user_id" binding:"-"`
	Phone  string `form:"phone" json:"phone" binding:"-"`
	Image  string `form:"image" json:"image" binding:"-"`
}

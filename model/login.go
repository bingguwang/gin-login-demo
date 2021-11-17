package model

import (
	"github.com/jinzhu/gorm"
)

type LoginModel struct {
	gorm.Model
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

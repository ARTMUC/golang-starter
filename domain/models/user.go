package models

import (
	"golang-starter/core/basemodel"
)

type User struct {
	//gorm.Model
	basemodel.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"-"`
}

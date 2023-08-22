package models

import (
	"github.com/golang-starter/core/basemodel"
	"github.com/google/uuid"
)

type Post struct {
	basemodel.Model
	Title       string     `gorm:"type:varchar(255)" json:"title,omitempty"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	CategoryID  *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Price       *float64   `json:"price,omitempty"`
}

type Category struct {
	basemodel.Model
	Name   string     `gorm:"type:varchar(255)" json:"name,omitempty"`
	Posts  []*Post    `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
	UserID *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	User   *User      `gorm:"foreignKey:UserID" json:"category,omitempty"`
}

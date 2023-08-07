package models

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	Title       string    `gorm:"type:varchar(255)" json:"title,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CategoryID  uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Price       *float64  `json:"price,omitempty"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
}

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	Name      string    `gorm:"type:varchar(255)" json:"name,omitempty"`
	Posts     []*Post   `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	//UserID    interface{} `gorm:"-" json:"user_id,omitempty"` // "-" means the field is ignored by GORM
	UserID uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	User   *User     `gorm:"foreignKey:UserID" json:"category,omitempty"`
}

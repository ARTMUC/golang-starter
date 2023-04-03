package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	//gorm.Model
	ID        uuid.UUID `json:"id,omitempty" gorm:"type:uuid; default:uuid_generate_v4()"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Password  string    `gorm:"size:255;not null;" json:"-"`
}

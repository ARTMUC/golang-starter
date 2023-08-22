package post

import "github.com/google/uuid"

type CreateDto struct {
	Title       string     `json:"title,omitempty" binding:"required"`
	Description string     `json:"description,omitempty" binding:"required"`
	CategoryID  *uuid.UUID `json:"category,omitempty" binding:"required,user-category"`
}

type UpdateDto struct {
	Title       string     `json:"title,omitempty" binding:"required"`
	Description string     `json:"description,omitempty" binding:"required"`
	CategoryID  *uuid.UUID `json:"category,omitempty" binding:"required,user-category"`
}

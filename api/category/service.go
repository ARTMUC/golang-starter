package category

import (
	"golang-starter/core/baserepo"
	"golang-starter/core/crud"
	"golang-starter/domain/models"
	"golang-starter/domain/repo"
)

type Service[T any] struct {
	crud.Service[T]
	repo baserepo.Dao[T]
}

func NewService[T models.Category](repository repo.CategoryRepo[T]) crud.Service[T] {
	return &Service[T]{
		repo:    repository,
		Service: crud.NewService[T](repository),
	}
}

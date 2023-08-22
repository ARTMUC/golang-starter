package category

import (
	"github.com/golang-starter/core/baserepo"
	"github.com/golang-starter/core/crud"
	"github.com/golang-starter/domain/models"
	"github.com/golang-starter/domain/repo"
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

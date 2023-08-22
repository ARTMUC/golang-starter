package post

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

func NewService[T models.Post](repository repo.PostRepo[T]) crud.Service[T] {
	return &Service[T]{
		repo:    repository,
		Service: crud.NewService[T](repository),
	}
}

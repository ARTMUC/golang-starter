package post

import (
	"github.com/golang-starter/crud"
	"github.com/golang-starter/domain/baserepo"
	"github.com/golang-starter/domain/repo"
)

type Service[T any] struct {
	crud.Service[T]
	repo baserepo.Dao[T]
}

func NewService[T model](repository repo.PostRepo[T]) *Service[T] {
	return &Service[T]{
		repo:    repository,
		Service: crud.NewService[T](repository),
	}
}

//var service = initService[model](repository)

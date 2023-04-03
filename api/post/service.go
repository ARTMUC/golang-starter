package post

import "github.com/golang-starter/crud"

type Service[T any] struct {
	crud.Service[T]
	repo crud.Dao[T]
}

func initService[T model](repository crud.Dao[T]) *Service[T] {
	return &Service[T]{
		repo:    repository,
		Service: crud.NewService[T](repository),
	}
}

var service = initService[model](repository)

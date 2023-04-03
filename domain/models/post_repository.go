package models

import (
	"github.com/golang-starter/crud"
	"github.com/golang-starter/db"
)

type PostRepo[T any] interface {
	crud.Dao[T]
}

type postRepo[T any] struct {
	crud.Repository[T]
}

func initPostRepository[T any]() PostRepo[T] {
	return &postRepo[T]{
		Repository: crud.Repository[T]{
			DB:    db.DB,
			Model: Post{},
		},
	}
}

var PostRepository = initPostRepository[Post]()

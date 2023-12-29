package repo

import (
	"golang-starter/core/baserepo"
	"golang-starter/core/db"
	"golang-starter/domain/models"
)

type PostRepo[T any] interface {
	baserepo.Dao[T]
}

type postRepo[T any] struct {
	baserepo.Repository[T]
}

func NewPostRepository[T any]() PostRepo[T] {
	return &postRepo[T]{
		Repository: baserepo.Repository[T]{
			DB:    db.DB,
			Model: models.Post{},
		},
	}
}

//var PostRepository = initPostRepository[Post]()

package models

import (
	"github.com/golang-starter/crud"
	"github.com/golang-starter/db"
)

type CategoryRepo[T any] interface {
	crud.Dao[T]
}

type categoryRepo[T Category] struct {
	crud.Repository[T]
}

func initCategoryRepository[T Category]() CategoryRepo[T] {
	return &categoryRepo[T]{
		Repository: crud.Repository[T]{
			DB:    db.DB,
			Model: T{},
		},
	}
}

var CategoryRepository = initCategoryRepository[Category]()

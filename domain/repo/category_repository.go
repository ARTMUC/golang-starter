package repo

import (
	"golang-starter/core/baserepo"
	"golang-starter/core/db"
	"golang-starter/domain/models"
)

type CategoryRepo[T any] interface {
	baserepo.Dao[T]
}

type categoryRepo[T models.Category] struct {
	baserepo.Repository[T]
}

func NewCategoryRepository[T models.Category]() CategoryRepo[T] {
	return &categoryRepo[T]{
		Repository: baserepo.Repository[T]{
			DB:    db.DB,
			Model: T{},
		},
	}
}

//var CategoryRepository = initCategoryRepository[Category]()

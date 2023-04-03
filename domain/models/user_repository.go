package models

import (
	"github.com/golang-starter/crud"
	"github.com/golang-starter/db"
)

type UserRepo[T any] interface {
	crud.Dao[T]
}

type userRepo[T any] struct {
	crud.Repository[T]
}

func initUserRepository[T any]() UserRepo[T] {
	return &userRepo[T]{
		Repository: crud.Repository[T]{
			DB:    db.DB,
			Model: User{},
		},
	}
}

var UserRepository = initUserRepository[User]()

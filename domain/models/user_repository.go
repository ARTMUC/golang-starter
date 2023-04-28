package models

import (
	"database/sql"
	"github.com/golang-starter/crud"
	"github.com/golang-starter/db"
)

type UserRepo[T any] interface {
	crud.Dao[T]
	FindOneByName(name string) (user *User, err error)
}

type userRepo[T any] struct {
	crud.Repository[T]
}

func (u *userRepo[T]) FindOneByName(name string) (user *User, err error) {
	err = u.DB.Where("user.name = ?", name).First(user).Error
	if err != nil {
		return nil, sql.ErrNoRows
	}
	return user, nil
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

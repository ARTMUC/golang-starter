package repo

import (
	"database/sql"
	"github.com/golang-starter/db"
	"github.com/golang-starter/domain/baserepo"
	"github.com/golang-starter/domain/models"
)

type UserRepo[T any] interface {
	baserepo.Dao[T]
	FindOneByName(name string) (user *models.User, err error)
}

type userRepo[T any] struct {
	baserepo.Repository[T]
}

func (u *userRepo[T]) FindOneByName(name string) (user *models.User, err error) {
	err = u.DB.Where("user.name = ?", name).First(user).Error
	if err != nil {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func NewUserRepository[T any]() UserRepo[T] {
	return &userRepo[T]{
		Repository: baserepo.Repository[T]{
			DB:    db.DB,
			Model: models.User{},
		},
	}
}

//var UserRepository = initUserRepository[User]()

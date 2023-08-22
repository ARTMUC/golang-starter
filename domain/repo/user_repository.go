package repo

import (
	"github.com/golang-starter/core/baserepo"
	"github.com/golang-starter/core/db"
	"github.com/golang-starter/domain/models"
)

type UserRepo[T any] interface {
	baserepo.Dao[T]
	FindOneByName(name string) (user *models.User, err error)
}

type userRepo[T any] struct {
	baserepo.Repository[T]
}

func (u *userRepo[T]) FindOneByName(name string) (*models.User, error) {
	user := &models.User{}
	if err := u.DB.Where("username = ?", name).First(user).Error; err != nil {
		return nil, err
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

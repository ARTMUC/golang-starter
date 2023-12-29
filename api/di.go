package api

import (
	"golang-starter/api/auth"
	"golang-starter/api/category"
	"golang-starter/api/post"
	"golang-starter/domain/models"
)

func GetProviders() []interface{} {
	return []interface{}{
		post.NewService[models.Post],
		post.NewController[models.Post],
		category.NewService[models.Category],
		category.NewController[models.Category],
		auth.NewController[models.User],
	}
}

var Controllers = []any{
	(*post.Controller[models.Post])(nil),
	(*category.Controller[models.Category])(nil),
	(*auth.Controller[models.User])(nil),
}

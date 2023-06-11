package api

import (
	"github.com/golang-starter/api/auth"
	"github.com/golang-starter/api/post"
	"github.com/golang-starter/domain/models"
)

func GetProviders() []interface{} {
	return []interface{}{
		post.NewService[models.Post],
		post.NewController[models.Post],
		auth.NewController[models.User],
	}
}

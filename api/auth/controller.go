package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/core/router"
	"github.com/golang-starter/domain/models"
	"github.com/golang-starter/domain/repo"
	"net/http"
)

type Controller[T models.User] struct {
	userRepository repo.UserRepo[T]
}

func NewController[T models.User](userRepository repo.UserRepo[T]) *Controller[T] {
	return &Controller[T]{userRepository}
}

func (c *Controller[T]) GetMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

func (c *Controller[T]) MainPath() string {
	return "auth"
}

func (c *Controller[T]) GetRoutes() []router.Handler {
	return []router.Handler{
		{
			Method:  http.MethodPost,
			Path:    "signin",
			Handler: func(ctx *gin.Context) { router.WrapResult(c.signin(ctx))(ctx) },
		},
		{
			Method:  http.MethodPost,
			Path:    "register",
			Handler: func(ctx *gin.Context) { router.WrapResult(c.register(ctx))(ctx) },
		},
	}
}

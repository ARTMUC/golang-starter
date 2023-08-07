package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/domain/models"
	"github.com/golang-starter/domain/repo"
	"github.com/golang-starter/router"
	"net/http"
)

type Controller[T models.User] struct {
	userRepository repo.UserRepo[T]
	router         *router.Routes
}

func NewController[T models.User](userRepository repo.UserRepo[T], router *router.Routes) *Controller[T] {
	controller := &Controller[T]{userRepository, router}
	router.AddController(controller)
	return controller
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
			Path:    "register",
			Handler: c.register,
		},
		{
			Method:  http.MethodPost,
			Path:    "signin",
			Handler: c.signin,
		},
	}
}

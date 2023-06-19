package auth

import (
	"github.com/gin-gonic/gin"
	sw "github.com/go-swagno/swagno"
	"github.com/golang-starter/domain/models"
	"github.com/golang-starter/domain/repo"
	"github.com/golang-starter/routes"
	"net/http"
)

type Controller[T models.User] struct {
	userRepository repo.UserRepo[T]
	router         *routes.Routes
}

func NewController[T models.User](userRepository repo.UserRepo[T], router *routes.Routes) *Controller[T] {
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

func (c *Controller[T]) GetRoutes() []routes.Handler {
	return []routes.Handler{
		{
			Docs:   sw.Endpoint{Body: &RegisterInput{}},
			Method: http.MethodPost,
			Path:   "register",
			Handler: func(ctx *gin.Context) {
				routes.WrapResult(c.register(ctx))
			},
		},
		{
			Docs:   sw.Endpoint{Body: &LoginInput{}},
			Method: http.MethodPost,
			Path:   "signin",
			Handler: func(ctx *gin.Context) {
				routes.WrapResult(c.signin(ctx))
			},
		},
	}
}

package post

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/crud"
	"github.com/golang-starter/middlewares"
	"github.com/golang-starter/routes"
)

type Controller[T ResponseDto] struct {
	*crud.Controller[T]
	router *routes.Routes
}

func NewController[T ResponseDto](router *routes.Routes, service crud.Service[T]) *Controller[T] {
	controller := &Controller[T]{
		router:     router,
		Controller: crud.NewController[T](crudConfig[T](), service),
	}
	controller.router.AddController(controller)

	return controller
}

func (c *Controller[T]) GetMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{middlewares.JwtAuthMiddleware()}
}

func (c *Controller[T]) MainPath() string {
	return crudName
}

func (c *Controller[T]) GetRoutes() []routes.Handler {
	return c.Endpoints
}

package post

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/crud"
	"github.com/golang-starter/middlewares"
	"github.com/golang-starter/routes"
)

type Controller[T model] struct {
	*crud.Controller[T]
	router *routes.Routes
}

func NewController[T model](router *routes.Routes) *Controller[T] {
	controller := &Controller[T]{
		router:     router,
		Controller: crud.NewController[T](crudConfig[T]()),
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

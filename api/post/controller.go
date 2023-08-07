package post

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/crud"
	"github.com/golang-starter/domain/models"
	"github.com/golang-starter/middlewares"
	"github.com/golang-starter/pkg/jwt"
	"github.com/golang-starter/router"
)

type ResponseDto = models.Post

var crudName = "post"

type Controller[T ResponseDto] struct {
	*crud.Controller[T]
	router *router.Routes
}

func NewController[T ResponseDto](router *router.Routes, service crud.Service[T]) *Controller[T] {
	controller := &Controller[T]{
		router: router,
		Controller: crud.NewController[T](
			&crud.Config[T]{
				ReadConstraint: &crud.ReadConstraint{
					Joins:  []string{"Author"},
					Field:  "author.id",
					Getter: jwt.MustExtractTokenID,
				},
				CreateDto:           &CreateDto{},
				UpdateDto:           &UpdateDto{},
				ResponseDto:         &ResponseDto{},
				DefaultCrudHandlers: crud.AllDefaultActions,
				Hooks:               &crud.Hooks[T]{},
			},
			service,
		),
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

func (c *Controller[T]) GetRoutes() []router.Handler {
	return c.Endpoints
}

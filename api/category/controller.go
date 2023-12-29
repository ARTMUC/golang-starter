package category

import (
	"github.com/gin-gonic/gin"
	crud2 "golang-starter/core/crud"
	"golang-starter/core/router"
	"golang-starter/domain/models"
	"golang-starter/middlewares"
)

type ResponseDto = models.Category

type Controller[T models.Category] struct {
	*crud2.Controller[T]
}

func NewController[T models.Category](service crud2.Service[T]) *Controller[T] {
	return &Controller[T]{
		Controller: crud2.NewController[T](
			&crud2.Config[T]{
				//ReadConstraint: &crud2.ReadConstraint{
				//	Joins:  []string{"Author"},
				//	Field:  "author.id",
				//	Getter: jwt.MustExtractTokenID,
				//},
				CreateDto:           &CreateDto{},
				UpdateDto:           &UpdateDto{},
				ResponseDto:         &ResponseDto{},
				DefaultCrudHandlers: crud2.AllDefaultActions,
				Hooks:               &crud2.Hooks[T]{},
			},
			service,
		),
	}
}

func (c *Controller[T]) GetMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{middlewares.JwtAuthMiddleware()}
}

func (c *Controller[T]) MainPath() string {
	return "category"
}

func (c *Controller[T]) GetRoutes() []router.Handler {
	return c.Endpoints
}

package post

import (
	"github.com/gin-gonic/gin"
	"golang-starter/core/crud"
	"golang-starter/core/router"
	"golang-starter/domain/models"
	"golang-starter/middlewares"
	"golang-starter/pkg/jwt"
)

type ResponseDto = models.Post

type Controller[T models.Post] struct {
	*crud.Controller[T]
}

func NewController[T models.Post](service crud.Service[T]) *Controller[T] {
	return &Controller[T]{
		Controller: crud.NewController[T](
			&crud.Config[T]{
				ReadConstraint: &crud.ReadConstraint{
					Joins: []string{
						"Inner Join categories ON posts.category_id = categories.id",
						"Inner Join users ON users.id = categories.user_id"},
					Field:  "users.id",
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
}

func (c *Controller[T]) GetMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{middlewares.JwtAuthMiddleware()}
}

func (c *Controller[T]) MainPath() string {
	return "post"
}

func (c *Controller[T]) GetRoutes() []router.Handler {
	return c.Endpoints
}

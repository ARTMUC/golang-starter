package post

import (
	"github.com/gin-gonic/gin"
	sw "github.com/go-swagno/swagno"
	"github.com/golang-starter/crud"
	"github.com/golang-starter/middlewares"
	"github.com/golang-starter/routes"
	"net/http"
)

type Controller[T model] struct {
	*crud.Controller[T]
	service crud.Service[T]
	router  *routes.Routes
}

func NewController[T model](service *Service[T], router *routes.Routes) *Controller[T] {
	controller := &Controller[T]{
		service:    service,
		router:     router,
		Controller: crud.NewController[T](service, defaultCrudConfig[T]()),
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
	return []routes.Handler{
		{
			Docs:   sw.Endpoint{Params: sw.Params(sw.IntParam("id", true, "")), Return: &model{}},
			Method: http.MethodGet,
			Path:   ":id",
			Handler: func(ctx *gin.Context) {
				c.FindOne(ctx, nil, nil)
			},
		},
		{
			Docs: sw.Endpoint{Return: crud.PaginationResponse[model]{},
				Params: sw.Params(
					sw.StrQuery("s", false, "{'$and': [ {'title': { '$cont':'cul' } } ]}"),
					sw.StrQuery("fields", false, "fields to select eg: name,age"),
					sw.IntQuery("page", false, "page of pagination"),
					sw.IntQuery("limit", false, "limit of pagination"),
					sw.StrQuery("join", false, "join relations eg: category, parent"),
					sw.StrQuery("filter", false, "filters eg: name||$eq||ad price||$gte||200"),
					sw.StrQuery("sort", false, "filters eg: created_at,desc title,asc"),
				),
			},
			Method: http.MethodGet,
			Path:   "",
			Handler: func(ctx *gin.Context) {
				c.FindAll(ctx, nil, nil)
			},
		},
		{
			Docs:   sw.Endpoint{Body: &CreateDto{}, Return: &model{}},
			Method: http.MethodPost,
			Path:   "",
			Handler: func(ctx *gin.Context) {
				c.Create(ctx, nil, nil)
			},
		},
		{
			Docs:   sw.Endpoint{Body: &UpdateDto{}, Params: sw.Params(sw.IntParam("id", true, "")), Return: &model{}},
			Method: http.MethodPatch,
			Path:   ":id",
			Handler: func(ctx *gin.Context) {
				c.Update(ctx, nil, nil)
			},
		},
		{
			Docs:   sw.Endpoint{Params: sw.Params(sw.IntParam("id", true, ""))},
			Method: http.MethodDelete,
			Path:   ":id",
			Handler: func(ctx *gin.Context) {
				c.Delete(ctx, nil, nil)
			},
		},
	}
}

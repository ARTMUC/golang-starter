package routes

import (
	"github.com/gin-gonic/gin"
	sw "github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-gin/swagger"
	"github.com/golang-starter/container"
	"github.com/golang-starter/di"
	"net/http"
)

type Controller interface {
	GetRoutes() []Handler
	MainPath() string
	GetMiddlewares() []gin.HandlerFunc
}

type Param struct {
	Name        string
	Required    bool
	Description string
}

type Docs struct {
	Params      []string
	QueryParams []string
	Body        interface{}
}

type Handler struct {
	Method  string
	Handler func(*gin.Context)
	Path    string
	Docs    sw.Endpoint
}

type Routes struct {
	controllers []interface{}
}

func NewRoutes() *Routes {
	return &Routes{controllers: []interface{}{}}
}

func (r *Routes) AddController(c interface{}) {
	r.controllers = append(r.controllers, c)
}

func (r *Routes) RegisterRoutes(e *gin.Engine) {
	controllers := di.GetMany[Controller](container.Container, r.controllers)

	var docs []sw.Endpoint
	for _, controller := range controllers {
		group := e.Group(controller.MainPath())
		routes := controller.GetRoutes()
		for _, route := range routes {
			e.Handle(route.Method, route.Path, route.Handler)

			doc := route.Docs
			doc.Path = "/" + controller.MainPath() + route.Path
			doc.Method = route.Method
			docs = append(docs, doc)
		}

		// @TODO: add option to add middlewares to each route separatelly
		for _, middleware := range controller.GetMiddlewares() {
			group.Use(middleware)
		}

		swaggerDocs := sw.CreateNewSwagger("Swagger API", "1.0")
		sw.AddEndpoints(docs)
		e.GET("/docs/*any", swagger.SwaggerHandler(swaggerDocs.GenerateDocs()))
	}
}

func WrapResult(result interface{}, err error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err != nil {
			c.Error(err)
		}

		successStatusCode := http.StatusOK
		if c.Request.Method == http.MethodPost {
			successStatusCode = http.StatusCreated
		}

		c.JSON(successStatusCode, result)
	}
}

package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/container"
	"github.com/golang-starter/core/di"
	"github.com/golang-starter/pkg/httperr"
	"net/http"
	"reflect"
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

type Handler struct {
	Method  string
	Handler func(*gin.Context)
	Path    string
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

	for _, controller := range controllers {
		group := e.Group(controller.MainPath())
		for _, route := range controller.GetRoutes() {
			handlers := controller.GetMiddlewares()
			handlers = append(handlers, route.Handler)
			group.Handle(
				route.Method,
				route.Path,
				handlers...,
			)

		}

		// @TODO: add option to add middlewares to each route separatelly
		for _, middleware := range controller.GetMiddlewares() {
			group.Use(middleware)
		}

	}
}

func WrapResult(result any, err error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err != nil {
			fmt.Println(reflect.TypeOf(err))
			fmt.Println(err.Error())
			switch e := err.(type) {
			case httperr.ErrCustomError:
				c.AbortWithStatusJSON(e.StatusCode, map[string]string{"message": e.Message})
				return
			default:
				c.Error(err)
				return
			}
		}

		successStatusCode := http.StatusOK
		if c.Request.Method == http.MethodPost {
			successStatusCode = http.StatusCreated
		}

		c.JSON(successStatusCode, result)
	}
}

package post

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/crud"
)

//var CrudController = initController(service)

type Controller[T model] struct {
	service crud.Service[T]
	*crud.Controller[T]
	constraint *crud.ReadConstraint
}

func (c *Controller[T]) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("", c.actionList)
	routerGroup.GET(":id", c.actionGet)
	routerGroup.POST("", c.actionCreate)
	routerGroup.DELETE(":id", c.actionDelete)
	routerGroup.PATCH(":id", c.actionUpdate)
}

// actionUpdate godoc
// @Success	200	{string}	string	"ok"
// @Tags		CrudName
// @param		id		path	string		true	"uuid of item"
// @param		item	body	UpdateDto	true	"update body"
// @Router		CrudName/{id} [put]
func (c *Controller[T]) actionUpdate(ctx *gin.Context) {
	c.Update(ctx, nil, nil)
}

// actionDelete godoc
// @Success	200	{string}	string	"ok"
// @Tags		CrudName
// @param		id	path	string	true	"uuid of item"
// @Router		CrudName/{id} [delete]
func (c *Controller[T]) actionDelete(ctx *gin.Context) {
	c.Delete(ctx, nil, nil)
}

// actionCreate godoc
// @Success	201	{object}	model
// @Tags		CrudName
// @param		{object}	body	CreateDto	true	"item to create"
// @Router		CrudName [post]
func (c *Controller[T]) actionCreate(ctx *gin.Context) {
	c.Create(ctx, nil, nil)
}

// actionGet godoc
// @Success	200	{object}	model
// @Tags		CrudName
// @param		id	path	string	true	"uuid of item"
// @Router		CrudName/{id} [get]
func (c *Controller[T]) actionGet(ctx *gin.Context) {
	c.FindOne(ctx, nil, nil)
}

// actionList godoc
// @Success	200	{array}	model
// @Tags		CrudName
// @param		s		query	string		false	"{'$and': [ {'title': { '$cont':'cul' } } ]}"
// @param		fields	query	string		false	"fields to select eg: name,age"
// @param		page	query	int			false	"page of pagination"
// @param		limit	query	int			false	"limit of pagination"
// @param		join	query	string		false	"join relations eg: category, parent"
// @param		filter	query	[]string	false	"filters eg: name||$eq||ad price||$gte||200"
// @param		sort	query	[]string	false	"filters eg: created_at,desc title,asc"
// @Router		CrudName [get]
func (c *Controller[T]) actionList(ctx *gin.Context) {
	c.FindAll(ctx, nil, nil)
}

func NewController[T model](service *Service[T]) *Controller[T] {
	return &Controller[T]{
		service:    service,
		Controller: crud.NewController[T](service, getConfig[T]()),
	}
}

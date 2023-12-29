package crud

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-starter/core/router"
	"golang-starter/pkg/httperr"
	"github.com/jinzhu/copier"
	"math"
	"net/http"
)

// @TODO: change GetAllRequest to pointer!!!
type Hooks[T any] struct {
	BeforeSave   func(data *T) error
	AfterSave    func(data *T) error
	BeforeUpdate func(data *T) error
	AfterUpdate  func(data *T) error
	BeforeGet    func(queryParams GetAllRequest) error
	AfterGet     func(data *T) error
	BeforeList   func(queryParams GetAllRequest) error
	AfterList    func(data *PaginationResponse[T]) error
	BeforeDelete func(data *T) error
	AfterDelete  func(data *T) error
}

type Controller[T any] struct {
	service             Service[T]
	constraint          *ReadConstraint
	createDto           Dto
	updateDto           Dto
	responseDto         Dto
	defaultCrudHandlers []Action
	Endpoints           []router.Handler
	hooks               *Hooks[T]
}

type ReadConstraint struct {
	Field  string
	Getter func(*gin.Context) string
	Joins  []string
}

type Config[T any] struct {
	ReadConstraint      *ReadConstraint
	CreateDto           Dto
	UpdateDto           Dto
	ResponseDto         Dto
	DefaultCrudHandlers []Action
	Hooks               *Hooks[T]
}

type Dto interface {
}

var AllDefaultActions = []Action{ActionCreate, ActionGet, ActionList, ActionUpdate, ActionDelete}

type Action = string

const (
	ActionCreate Action = "create"
	ActionGet    Action = "get"
	ActionList   Action = "list"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

func NewController[T any](config *Config[T], service Service[T]) *Controller[T] {
	controller := &Controller[T]{
		service:             service,
		constraint:          config.ReadConstraint,
		createDto:           config.CreateDto,
		updateDto:           config.UpdateDto,
		responseDto:         config.ResponseDto,
		defaultCrudHandlers: config.DefaultCrudHandlers,
		hooks:               config.Hooks,
	}
	controller.addDefaultEndpoints()
	return controller
}

func (c *Controller[T]) FindAll(ctx *gin.Context) (any, error) {
	var queryParams GetAllRequest

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	// @TODO throw this to struct -> validate for minimum
	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}
	if queryParams.Page < 1 {
		queryParams.Page = 1
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(&queryParams, ctx)
	}

	if f := c.hooks.BeforeList; f != nil {
		if err := f(queryParams); err != nil {
			return nil, err
		}
	}

	result, totalRows, err := c.service.Find(queryParams)
	if err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	data := &PaginationResponse[T]{
		Data:       result,
		Total:      totalRows,
		TotalPages: int64(math.Ceil(float64(totalRows) / float64(queryParams.Limit))),
	}

	if f := c.hooks.AfterList; f != nil {
		if err := f(data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (c *Controller[T]) FindOne(ctx *gin.Context) (any, error) {
	var queryParams GetAllRequest
	var pathParams ById
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}
	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(&queryParams, ctx)
	}

	queryParams.Filter = append(queryParams.Filter, fmt.Sprintf("id||eq||%s", pathParams.ID))

	if c.hooks.BeforeGet != nil {
		if err := c.hooks.BeforeGet(queryParams); err != nil {
			return nil, err
		}
	}

	result, err := c.service.FindOne(queryParams)
	if err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	if c.hooks.AfterGet != nil {
		if err := c.hooks.AfterGet(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Controller[T]) Create(ctx *gin.Context) (any, error) {
	dto := c.createDto
	var item = new(T)
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}
	if err := copier.CopyWithOption(item, dto, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if f := c.hooks.BeforeSave; f != nil {
		if err := f(item); err != nil {
			return nil, err
		}
	}

	if err := c.service.Create(item); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if f := c.hooks.AfterSave; f != nil {
		if err := f(item); err != nil {
			return nil, err
		}
	}

	return item, nil
}

func (c *Controller[T]) Update(ctx *gin.Context) (any, error) {
	var queryParams GetAllRequest
	dto := &c.updateDto
	var item *T
	var pathParams ById
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}
	if err := copier.CopyWithOption(item, dto, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}
	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(&queryParams, ctx)
	}
	queryParams.Filter = append(queryParams.Filter, fmt.Sprintf("id||eq||%s", pathParams.ID))

	result, err := c.service.FindOne(queryParams)
	if err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	if f := c.hooks.BeforeUpdate; f != nil {
		if err := f(item); err != nil {
			return nil, err
		}
	}

	if err = c.service.Update(result, item); err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	if f := c.hooks.AfterUpdate; f != nil {
		if err := f(item); err != nil {
			return nil, err
		}
	}

	return item, nil
}

func (c *Controller[T]) Delete(ctx *gin.Context) (any, error) {
	var queryParams GetAllRequest
	var pathParams ById
	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(&queryParams, ctx)
	}

	queryParams.Filter = append(queryParams.Filter, fmt.Sprintf("id||eq||%s", pathParams.ID))

	result, err := c.service.FindOne(queryParams)
	if err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	if f := c.hooks.BeforeDelete; f != nil {
		if err := f(result); err != nil {
			return nil, err
		}
	}

	if err = c.service.Delete(result); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if f := c.hooks.AfterDelete; f != nil {
		if err := f(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Controller[T]) joinConstraint(queryParams *GetAllRequest, ctx *gin.Context) {
	queryParams.Joins = c.constraint.Joins
	queryParams.Filter = append(queryParams.Filter, fmt.Sprintf("%s||eq||%s", c.constraint.Field, c.constraint.Getter(ctx)))
}

func (c *Controller[T]) addDefaultEndpoints() {
	defaultActions := map[Action]router.Handler{
		ActionGet: {
			Method:  http.MethodGet,
			Path:    ":id",
			Handler: func(ctx *gin.Context) { router.WrapResult(c.FindOne(ctx))(ctx) },
		},
		ActionList: {
			Method:  http.MethodGet,
			Path:    "",
			Handler: func(ctx *gin.Context) { router.WrapResult(c.FindAll(ctx))(ctx) },
		},
		ActionCreate: {
			Method:  http.MethodPost,
			Path:    "",
			Handler: func(ctx *gin.Context) { router.WrapResult(c.Create(ctx))(ctx) },
		},
		ActionUpdate: {
			Method:  http.MethodPatch,
			Path:    ":id",
			Handler: func(ctx *gin.Context) { router.WrapResult(c.Update(ctx))(ctx) },
		},
		ActionDelete: {
			Method:  http.MethodDelete,
			Path:    ":id",
			Handler: func(ctx *gin.Context) { router.WrapResult(c.Delete(ctx))(ctx) },
		},
	}

	for _, handler := range c.defaultCrudHandlers {
		if endpoint, ok := defaultActions[handler]; ok {
			c.Endpoints = append(c.Endpoints, endpoint)
		}
	}
}

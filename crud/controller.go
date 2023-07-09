package crud

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/pkg/httperr"
	"github.com/golang-starter/routes"
	"github.com/jinzhu/copier"
	"math"
	"net/http"
	"strings"
)

type Controller[T any] struct {
	service             Service[T]
	constraint          *ReadConstraint
	createDto           Dto
	updateDto           Dto
	responseDto         Dto
	defaultCrudHandlers []Action
	Endpoints           []routes.Handler
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
}

type Dto interface {
}

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
	}
	controller.addDefaultEndpoints()
	return controller
}

func (c *Controller[T]) FindAll(ctx *gin.Context, beforeFn func(queryParams GetAllRequest) error, afterFn func(data *PaginationResponse[T]) error) (*PaginationResponse[T], error) {
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
		c.joinConstraint(queryParams, ctx)
	}

	if beforeFn != nil {
		if err := beforeFn(queryParams); err != nil {
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

	if afterFn != nil {
		if err := afterFn(data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (c *Controller[T]) FindOne(ctx *gin.Context, beforeFn func(queryParams GetAllRequest) error, afterFn func(data *T) error) (*T, error) {
	var queryParams GetAllRequest
	var pathParams ById
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}
	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(queryParams, ctx)
	}

	queryParams.Filter = append(queryParams.Filter, fmt.Sprintf("id||eq||%s", pathParams.ID))

	if beforeFn != nil {
		if err := beforeFn(queryParams); err != nil {
			return nil, err
		}
	}

	result, err := c.service.FindOne(queryParams)
	if err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	if afterFn != nil {
		if err := afterFn(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Controller[T]) Create(ctx *gin.Context, beforeFn func(item *T) error, afterFn func(data *T) error) (*T, error) {
	dto := c.createDto
	var item *T
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}
	if err := copier.CopyWithOption(item, dto, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if beforeFn != nil {
		if err := beforeFn(item); err != nil {
			return nil, err
		}
	}

	if err := c.service.Create(item); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if afterFn != nil {
		if err := afterFn(item); err != nil {
			return nil, err
		}
	}

	return item, nil
}

func (c *Controller[T]) Update(ctx *gin.Context, beforeFn func(item *T) error, afterFn func(data *T) error) (*T, error) {
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
		c.joinConstraint(queryParams, ctx)
	}
	queryParams.Filter = append(queryParams.Filter, fmt.Sprintf("id||eq||%s", pathParams.ID))

	result, err := c.service.FindOne(queryParams)
	if err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	if beforeFn != nil {
		if err := beforeFn(item); err != nil {
			return nil, err
		}
	}

	if err = c.service.Update(result, item); err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	if afterFn != nil {
		if err := afterFn(item); err != nil {
			return nil, err
		}
	}

	return item, nil
}

func (c *Controller[T]) Delete(ctx *gin.Context, beforeFn func(item *T) error, afterFn func(data *T) error) (*T, error) {
	var queryParams GetAllRequest
	var pathParams ById
	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(queryParams, ctx)
	}

	queryParams.Filter = append(queryParams.Filter, fmt.Sprintf("id||eq||%s", pathParams.ID))

	result, err := c.service.FindOne(queryParams)
	if err != nil {
		return nil, httperr.NewNotFoundError("not found", err)
	}

	if beforeFn != nil {
		if err := beforeFn(result); err != nil {
			return nil, err
		}
	}

	if err = c.service.Delete(result); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if afterFn != nil {
		if err := afterFn(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Controller[T]) joinConstraint(queryParams GetAllRequest, ctx *gin.Context) {
	queryParams.C = map[string]interface{}{c.constraint.Field: c.constraint.Getter(ctx)}

	var fullJoins []string
	joinsArray := strings.Split(queryParams.Join, ",")
	joinsMap := make(map[string]bool)

	for _, s := range joinsArray {
		joinsMap[s] = true
	}

	if len(c.constraint.Joins) > 0 {
		for _, join := range c.constraint.Joins {
			if !joinsMap[strings.ToLower(join)] {
				fullJoins = append(fullJoins, join)
			}
		}
	}

	queryParams.Join = strings.Join(fullJoins, ",")
}

func (c *Controller[T]) addDefaultEndpoints() {
	//createDto := c.docsCreateDto
	//updateDto := c.docsUpdateDto
	//responseDto := c.docsResponseDto
	defaultActions := map[Action]routes.Handler{
		ActionGet: {
			//Docs:   sw.Endpoint{Params: sw.Params(sw.IntParam("id", true, "")), Return: responseDto},
			Method: http.MethodGet,
			Path:   ":id",
			Handler: func(ctx *gin.Context) {
				routes.WrapResult(c.FindOne(ctx, nil, nil))
			},
		},
		ActionList: {
			//Docs: sw.Endpoint{Return: PaginationResponse[T]{}, Params: sw.Params(
			//	sw.StrQuery("s", false, "{'$and': [ {'title': { '$cont':'cul' } } ]}"),
			//	sw.StrQuery("fields", false, "fields to select eg: name,age"),
			//	sw.IntQuery("page", false, "page of pagination"),
			//	sw.IntQuery("limit", false, "limit of pagination"),
			//	sw.StrQuery("join", false, "join relations eg: category, parent"), // @TODO we should restrict joins
			//	sw.StrQuery("filter", false, "filters eg: name||$eq||ad price||$gte||200"),
			//	sw.StrQuery("sort", false, "filters eg: created_at,desc title,asc"),
			//)},
			Method: http.MethodGet,
			Path:   "",
			Handler: func(ctx *gin.Context) {
				routes.WrapResult(c.FindAll(ctx, nil, nil))
			},
		},
		ActionCreate: {
			//Docs:   sw.Endpoint{Body: createDto, Return: responseDto},
			Method: http.MethodPost,
			Path:   "",
			Handler: func(ctx *gin.Context) {
				routes.WrapResult(c.Create(ctx, nil, nil))
			},
		},
		ActionUpdate: {
			//Docs:   sw.Endpoint{Body: updateDto, Params: sw.Params(sw.IntParam("id", true, "")), Return: responseDto},
			Method: http.MethodPatch,
			Path:   ":id",
			Handler: func(ctx *gin.Context) {
				routes.WrapResult(c.Update(ctx, nil, nil))
			},
		},
		ActionDelete: {
			//Docs:   sw.Endpoint{Params: sw.Params(sw.IntParam("id", true, ""))},
			Method: http.MethodDelete,
			Path:   ":id",
			Handler: func(ctx *gin.Context) {
				routes.WrapResult(c.Delete(ctx, nil, nil))
			},
		},
	}

	for _, handler := range c.defaultCrudHandlers {
		if endpoint, ok := defaultActions[handler]; ok {
			c.Endpoints = append(c.Endpoints, endpoint)
		}
	}
}

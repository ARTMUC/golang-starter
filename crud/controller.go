package crud

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/common"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"math"
	"net/http"
	"strings"
)

type Controller[T any] struct {
	service    Service[T]
	constraint *ReadConstraint
	createDto  interface{}
	updateDto  interface{}
}

type ReadConstraint struct {
	Field  string
	Getter func(*gin.Context) string
	Joins  []string
}

type Config[T any] struct {
	ReadConstraint *ReadConstraint
	CreateDto      Dto
	UpdateDto      Dto
}

type Dto interface {
}

func (c *Controller[T]) FindAll(ctx *gin.Context, before func(api GetAllRequest) error, after func(data *PaginationResponse[T]) error) {
	var api GetAllRequest
	if api.Limit == 0 {
		api.Limit = 20
	}
	if api.Page < 1 {
		api.Page = 1
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(api, ctx)
	}

	if err := ctx.ShouldBindQuery(&api); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var result []*T
	var totalRows int64

	if before != nil {
		err := before(api)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	err := c.service.Find(api, &result, &totalRows)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	data := &PaginationResponse[T]{
		Data:       result,
		Total:      totalRows,
		TotalPages: int64(math.Ceil(float64(totalRows) / float64(api.Limit))),
	}

	if after != nil {
		err := after(data)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(200, data)
}

func (c *Controller[T]) FindOne(ctx *gin.Context, before func(api GetAllRequest) error, after func(data *T) error) {
	var api GetAllRequest
	var item common.ById
	if err := ctx.ShouldBindQuery(&api); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if err := ctx.ShouldBindUri(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(api, ctx)
	}

	api.Filter = append(api.Filter, fmt.Sprintf("id||eq||%s", item.ID))

	var result *T

	if before != nil {
		err := before(api)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	err := c.service.FindOne(api, result)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if after != nil {
		err := after(result)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(200, result)
}

func (c *Controller[T]) Create(ctx *gin.Context, before func(item *T) error, after func(data *T) error) {
	dto := c.createDto
	var item *T
	if err := ctx.ShouldBind(dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := copier.CopyWithOption(item, dto, copier.Option{
		IgnoreEmpty: false,
		DeepCopy:    true,
		Converters:  nil,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if before != nil {
		err := before(item)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	err = c.service.Create(item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if after != nil {
		err := after(item)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": item})
}

func (c *Controller[T]) Update(ctx *gin.Context, before func(item *T) error, after func(data *T) error) {
	var api GetAllRequest
	dto := c.updateDto
	var item *T
	var byId common.ById
	if err := ctx.ShouldBind(dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := copier.CopyWithOption(item, dto, copier.Option{
		IgnoreEmpty: false,
		DeepCopy:    true,
		Converters:  nil,
	})
	if err := ctx.ShouldBindUri(&byId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id, err := uuid.Parse(byId.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := ctx.ShouldBindUri(&byId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(api, ctx)
	}

	api.Filter = append(api.Filter, fmt.Sprintf("id||eq||%s", id))

	var result *T

	err = c.service.FindOne(api, result)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if before != nil {
		err := before(item)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	err = c.service.Update(result, item)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if after != nil {
		err := after(item)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *Controller[T]) Delete(ctx *gin.Context, before func(item *T) error, after func(data *T) error) {
	var api GetAllRequest
	var item common.ById
	if err := ctx.ShouldBindUri(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if c.constraint.Getter != nil && len(c.constraint.Field) > 0 {
		c.joinConstraint(api, ctx)
	}

	id, err := uuid.ParseBytes([]byte(item.ID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	api.Filter = append(api.Filter, fmt.Sprintf("id||eq||%s", id))

	var result *T

	err = c.service.FindOne(api, result)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if before != nil {
		err := before(result)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	err = c.service.Delete(result)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if after != nil {
		err := after(result)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (c *Controller[T]) joinConstraint(api GetAllRequest, ctx *gin.Context) {
	api.C = map[string]interface{}{c.constraint.Field: c.constraint.Getter(ctx)}

	var fj []string
	joinsArray := strings.Split(api.Join, ",")
	if len(c.constraint.Joins) > 0 {
		for _, join := range joinsArray {
			for _, cstrJoin := range c.constraint.Joins {
				if strings.ToLower(join) != strings.ToLower(cstrJoin) {
					fj = append(fj, join)
				}
			}
		}
	}
	for _, join := range c.constraint.Joins {
		fj = append(fj, join)
	}
	api.Join = strings.Join(fj, ",")
}

func NewController[T any](service Service[T], config *Config[T]) *Controller[T] {
	return &Controller[T]{
		service:    service,
		constraint: config.ReadConstraint,
		createDto:  config.CreateDto,
		updateDto:  config.UpdateDto,
	}
}

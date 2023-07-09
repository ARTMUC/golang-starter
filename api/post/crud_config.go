package post

import (
	"github.com/golang-starter/crud"
	models "github.com/golang-starter/domain/models"
	"github.com/golang-starter/pkg/jwt"
)

type ResponseDto = models.Post

var crudName = "post"

func crudConfig[T any]() *crud.Config[T] {
	return &crud.Config[T]{
		ReadConstraint: &crud.ReadConstraint{
			Joins:  []string{"Author"},
			Field:  "author.id",
			Getter: jwt.MustExtractTokenID,
		},
		CreateDto:   &CreateDto{},
		UpdateDto:   &UpdateDto{},
		ResponseDto: &ResponseDto{},
		DefaultCrudHandlers: []crud.Action{
			crud.ActionCreate,
			crud.ActionGet,
			crud.ActionList,
			crud.ActionUpdate,
			crud.ActionDelete,
		},
	}
}

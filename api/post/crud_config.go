package post

import (
	"github.com/golang-starter/crud"
	models "github.com/golang-starter/domain/models"
	"github.com/golang-starter/pkg/jwt"
)

type model = models.Post

var crudName = "post"

func defaultCrudConfig[T any]() *crud.Config[T] {
	return &crud.Config[T]{
		ReadConstraint: &crud.ReadConstraint{
			Joins:  []string{"Author"},
			Field:  "author.id",
			Getter: jwt.MustExtractTokenID,
		},
		CreateDto: &CreateDto{},
		UpdateDto: &UpdateDto{},
	}
}

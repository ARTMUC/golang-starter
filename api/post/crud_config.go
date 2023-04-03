package post

import (
	"github.com/golang-starter/common"
	"github.com/golang-starter/crud"
	models "github.com/golang-starter/domain/models"
)

type model = models.Post

var repository = models.PostRepository

func getConfig[T any]() *crud.Config[T] {
	return &crud.Config[T]{
		ReadConstraint: &crud.ReadConstraint{
			Joins:  []string{"Author"},
			Field:  "author.id",
			Getter: common.MustExtractTokenID,
		},
		CreateDto: &CreateDto{},
		UpdateDto: &UpdateDto{},
	}
}

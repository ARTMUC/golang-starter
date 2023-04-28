package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-starter/domain/models"
	"github.com/google/uuid"
)

var (
	UserCategoryValidatorTag = "user-category"
)

func RegisterBindingsValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidationCtx(UserCategoryValidatorTag, func(ctx context.Context, fl validator.FieldLevel) bool {
			categoryID := fl.Field().String()
			c, ok := fl.Parent().Interface().(gin.Context)
			if !ok {
				return false
			}
			userID := MustExtractTokenID(&c)
			return models.CheckCategoryUser(uuid.MustParse(categoryID), uuid.MustParse(userID))
		})
	}
}

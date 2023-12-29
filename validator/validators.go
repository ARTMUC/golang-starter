package validator

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang-starter/container"
	"golang-starter/core/basemodel"
	"golang-starter/core/di"
	"golang-starter/domain/models"
	"golang-starter/domain/repo"
	"golang-starter/pkg/jwt"
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
			userID := jwt.MustExtractTokenID(&c)
			return CheckCategoryUser(uuid.MustParse(categoryID), uuid.MustParse(userID))
		})
	}
}

func CheckCategoryUser(categoryID uuid.UUID, userID uuid.UUID) bool {
	var categoryRepository repo.CategoryRepo[models.Category]
	di.GetProvider(container.Container, &categoryRepository)
	category := &models.Category{}

	err := categoryRepository.FindOne(&models.Category{Model: basemodel.Model{ID: categoryID}, UserID: &userID}, category)
	if err != nil {
		return false
	}
	return true
}

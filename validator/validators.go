package validator

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-starter/container"
	"github.com/golang-starter/di"
	"github.com/golang-starter/domain/models"
	"github.com/golang-starter/domain/repo"
	"github.com/golang-starter/pkg/jwt"
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

	var category *models.Category
	err := categoryRepository.FindOne(&models.Category{ID: categoryID, UserID: userID}, category)
	if err != nil {
		return false
	}
	return true
}

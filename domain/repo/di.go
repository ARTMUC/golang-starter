package repo

import "golang-starter/domain/models"

func GetProviders() []interface{} {
	return []interface{}{
		NewUserRepository[models.User],
		NewCategoryRepository[models.Category],
		NewPostRepository[models.Post],
	}
}

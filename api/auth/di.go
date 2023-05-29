package auth

import "github.com/golang-starter/domain/models"

func GetProviders() []interface{} {
	return []interface{}{
		NewController[models.User],
	}
}

package post

func GetProviders() []interface{} {
	return []interface{}{
		NewService[model],
		NewController[model],
	}
}

package router

func GetProviders() []interface{} {
	return []interface{}{
		NewRoutes,
	}
}

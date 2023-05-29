package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/api/auth"
	"github.com/golang-starter/api/post"
	"github.com/golang-starter/container"
	"github.com/golang-starter/di"
	"github.com/golang-starter/domain/models"
	"github.com/golang-starter/domain/repo"
	"github.com/golang-starter/middlewares"
	"github.com/golang-starter/validator"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	server := gin.New()
	server.Use(gin.Recovery())

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	server.Use(cors.New(config))

	if os.Getenv("GIN_MODE") == "debug" {
		server.Use(gin.Logger())
	}

	validator.RegisterBindingsValidators()

	//gin.SetMode(os.Getenv("GIN_MODE"))
	//if err := db.Open(os.Getenv("DB_URL")); err != nil {
	//	log.Fatal(err)
	//}
	log.Println("server started")

	server.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "ok 5"})
	})

	// migrations
	//db.AddUUIDExtension()

	//if err := db.DB.AutoMigrate(
	//	models.Category{},
	//	models.Post{},
	//); err != nil {
	//	log.Fatal(err)
	//}

	authGroup := server.Group("auth")
	postsGroup := server.Group(post.CrudName)

	postsGroup.Use(middlewares.JwtAuthMiddleware())

	var allDependencies = [][]interface{}{
		repo.GetProviders(),
		post.GetProviders(),
		auth.GetProviders(),
	}

	for _, dependencies := range allDependencies {
		for _, dependency := range dependencies {
			container.Container.Provide(dependency)
		}
	}
	
	di.MustGet(container.Container, &auth.Controller[models.User]{}).RegisterRoutes(authGroup)
	di.MustGet(container.Container, &post.Controller[models.Post]{}).RegisterRoutes(postsGroup)

	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := server.Run(":8081")
	if err != nil {
		log.Printf("error while starting server %+v", err)
	}
}

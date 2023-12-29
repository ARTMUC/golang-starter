package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang-starter/api"
	"golang-starter/container"
	"golang-starter/core/db"
	"golang-starter/core/di"
	router2 "golang-starter/core/router"
	"golang-starter/docs/docsgenerator"
	"golang-starter/domain/models"
	"golang-starter/domain/repo"
	"golang-starter/middlewares"
	"golang-starter/validator"
	"log"
	"os"
)

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(middlewares.LoggerMiddleware())
	server.Use(middlewares.ErrorHandler())

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	server.Use(cors.New(config))

	if os.Getenv("GIN_MODE") == "debug" {
		server.Use(gin.Logger())
	}

	validator.RegisterBindingsValidators()

	gin.SetMode(os.Getenv("GIN_MODE"))
	//if err := db.Open(os.Getenv("DB_URL")); err != nil {
	//	log.Fatal(err)
	//}

	if err := db.OpenTestDB(); err != nil {
		log.Fatal(err)
	}

	log.Println("server started")

	//migrations
	//db.AddUUIDExtension()
	if err := db.DB.AutoMigrate(
		models.User{},
		models.Category{},
		models.Post{},
	); err != nil {
		log.Fatal(err)
	}

	var allDependencies = [][]interface{}{
		repo.GetProviders(),
		api.GetProviders(),
	}

	for _, dependencies := range allDependencies {
		for _, dependency := range dependencies {
			container.Container.Provide(dependency)
		}
	}

	apiRouter := router2.NewRoutes()
	for _, controller := range api.Controllers {
		apiRouter.AddController(controller)
	}
	apiRouter.RegisterRoutes(server)

	docsgenerator.GenerateDocs(di.GetMany[any](container.Container, api.Controllers))

	//server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := server.Run(":8081")
	if err != nil {
		log.Printf("error while starting server %+v", err)
	}
}

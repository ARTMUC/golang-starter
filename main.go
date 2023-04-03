package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-starter/api/auth"
	"github.com/golang-starter/api/post"
	"github.com/golang-starter/common"
	"github.com/golang-starter/domain/models"
	"github.com/golang-starter/middlewares"
	"github.com/google/uuid"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidationCtx("user-category", func(ctx context.Context, fl validator.FieldLevel) bool {
			categoryID := fl.Field().String()
			c, ok := fl.Parent().Interface().(gin.Context)
			if !ok {
				return false
			}
			userID := common.MustExtractTokenID(&c)
			return models.CheckCategoryUser(uuid.MustParse(categoryID), uuid.MustParse(userID))
		})
	}

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
	postsGroup := server.Group("post")

	postsGroup.Use(middlewares.JwtAuthMiddleware())

	auth.AuthController.RegisterRoutes(authGroup)
	post.CrudController.RegisterRoutes(postsGroup)

	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := server.Run(":8081")
	if err != nil {
		log.Printf("error while starting server %+v", err)
	}
}

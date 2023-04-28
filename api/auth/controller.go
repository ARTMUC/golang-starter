package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/common"
	"github.com/golang-starter/domain/models"
	"golang.org/x/crypto/bcrypt"
	"html"
	"net/http"
	"strings"
)

type Controller[T models.User] struct {
	userRepository models.UserRepo[T]
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *Controller[T]) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("register", c.register)
	routerGroup.POST("signin", c.signin)
}

func (c *Controller[T]) register(ctx *gin.Context) {
	var input RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &T{
		Username: string(hashedPassword),
		Password: html.EscapeString(strings.TrimSpace(input.Username)),
	}

	err = c.userRepository.Create(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *Controller[T]) signin(ctx *gin.Context) {
	var input LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userRepository.FindOneByName(input.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(input.Password), []byte(user.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	token, err := common.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func initController[T models.User](userRepository models.UserRepo[T]) *Controller[T] {
	return &Controller[T]{userRepository}
}

var AuthController = initController[models.User](models.UserRepository)

package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

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

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

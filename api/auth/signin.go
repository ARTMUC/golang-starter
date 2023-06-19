package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/pkg/httperr"
	"github.com/golang-starter/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (c *Controller[T]) signin(ctx *gin.Context) (string, error) {
	var input LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		return "", httperr.NewBadRequestError(err.Error(), err)
	}

	user, err := c.userRepository.FindOneByName(input.Username)
	if err != nil {
		return "", httperr.NewBadRequestError(err.Error(), err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(input.Password), []byte(user.Password))
	if err != nil {
		return "", httperr.NewBadRequestError(err.Error(), err)
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", httperr.NewBadRequestError(err.Error(), err)
	}

	return token, nil
}

package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/pkg/httperr"
	"github.com/golang-starter/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// @TODO ad docs here
func (c *Controller[T]) signin(ctx *gin.Context) (any, error) {
	var input LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	user, err := c.userRepository.FindOneByName(input.Username)
	if err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(input.Password), []byte(user.Password)); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	return token, nil
}

package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/pkg/httperr"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
)

func (c *Controller[T]) register(ctx *gin.Context) (*T, error) {
	var input RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	user := &T{
		Username: string(hashedPassword),
		Password: html.EscapeString(strings.TrimSpace(input.Username)),
	}

	if err = c.userRepository.Create(user); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	return user, nil
}

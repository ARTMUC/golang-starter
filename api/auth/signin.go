package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-starter/pkg/httperr"
	"golang-starter/pkg/jwt"
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

	bytes, _ := json.Marshal(user)
	fmt.Println(string(bytes))

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, httperr.NewBadRequestError(err.Error(), err)
	}

	return token, nil
}

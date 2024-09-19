package middleware

import (
	"github.com/FlyKarlik/auth-service/internal/errs"
	"github.com/FlyKarlik/auth-service/pkg/validator"
	"github.com/gin-gonic/gin"
)

const (
	userKeyParam = "id"
)

func GetUserID(c *gin.Context) (string, error) {

	userId := c.Param(userKeyParam)
	if len(userId) == 0 {
		return "", errs.ErrUserIDNotFound
	}

	if err := validator.IsValidStringUUID(userId); err != nil {
		return "", errs.ErrInvalidUserID
	}

	return userId, nil
}

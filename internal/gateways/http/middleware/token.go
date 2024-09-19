package middleware

import (
	"context"

	"github.com/FlyKarlik/auth-service/internal/errs"
	"github.com/gin-gonic/gin"
)

const (
	access  = "access"
	refresh = "refresh"
)

func BindTokens(c *gin.Context, ctx context.Context) (context.Context, error) {
	accessToken := c.GetHeader(access)
	if len(accessToken) == 0 {
		return nil, errs.ErrIncorrectRerfreshOperation
	}

	refreshToken := c.GetHeader(refresh)
	if len(refreshToken) == 0 {
		return nil, errs.ErrIncorrectRerfreshOperation
	}

	ctx = context.WithValue(ctx, access, accessToken)
	ctx = context.WithValue(ctx, refresh, refreshToken)

	return ctx, nil
}

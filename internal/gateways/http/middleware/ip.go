package middleware

import (
	"github.com/FlyKarlik/auth-service/internal/errs"
	"github.com/gin-gonic/gin"
)

const (
	clientIPKey = "client-ip"
)

func BindClientIP(c *gin.Context) {
	clientIP := c.RemoteIP()
	c.Set(clientIPKey, clientIP)
	c.Next()
}

func GetClientIP(c *gin.Context) (string, error) {
	clientIp, ok := c.Get(clientIPKey)
	if !ok {
		return "", errs.ErrClientIPNotFound
	}

	strClientIP, ok := clientIp.(string)
	if !ok {
		return "", errs.ErrInvalidClientIP
	}

	return strClientIP, nil
}

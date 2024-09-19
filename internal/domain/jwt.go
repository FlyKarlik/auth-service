package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Token Expirations
var (
	AccessToken        = time.Minute * 15
	RefreshAccessToken = time.Hour * 2
)

type JWT struct {
	ID string
	jwt.StandardClaims
	ClientIP string
	UserID   string
	Variety  string
}

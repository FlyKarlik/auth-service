package auth

import (
	"context"
	"fmt"

	"github.com/FlyKarlik/auth-service/internal/domain"
	"github.com/FlyKarlik/auth-service/internal/errs"
	"github.com/FlyKarlik/auth-service/internal/tokens"
	"github.com/FlyKarlik/auth-service/pkg/codes"
	"github.com/FlyKarlik/auth-service/pkg/jwt"
)

const (
	access  = "access"
	refresh = "refresh"
)

type Auth struct {
	token tokens.IAuthTokens
}

func New(token tokens.IAuthTokens) *Auth {
	return &Auth{token: token}
}

func (a *Auth) Authentication(ctx context.Context, userId string, clientIP string) (*domain.Token, error) {

	var token *domain.Token

	variety, err := jwt.GenerateVariety()
	if err != nil {
		return nil, errs.New(codes.ErrorInternal, err.Error())
	}

	token, err = a.token.CreateAuthTokens(ctx, userId, clientIP, variety)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *Auth) Refresh(ctx context.Context, clientIP string) (*domain.Token, error) {
	accessToken, refreshToken := getAuthTokensFromContext(ctx)

	variety, err := jwt.GenerateVariety()
	if err != nil {
		return nil, errs.New(codes.ErrorInternal, err.Error())
	}

	if err := a.token.ValidateAuthTokens(ctx, accessToken, refreshToken); err != nil {
		return nil, err
	}

	accessClaims, err := a.token.ParseToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	refreshClaims, err := a.token.ParseToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	//TODO: Checking if the user exists

	if accessClaims.ClientIP != clientIP && refreshClaims.ClientIP != clientIP {
		accessClaims.ClientIP = clientIP
		refreshClaims.ClientIP = clientIP
		//TODO: Send an event about an ip change to mail via a broker or grpc
		fmt.Printf("client ip had changed")
	}

	accessClaims.Variety = variety
	refreshClaims.Variety = variety

	return a.token.RefreshAuthTokens(ctx, accessClaims, refreshClaims)

}

func getAuthTokensFromContext(ctx context.Context) (string, string) {
	accessToken := ctx.Value(access)
	strAccessToken, _ := accessToken.(string)

	refreshToken := ctx.Value(refresh)
	strRefreshToken, _ := refreshToken.(string)

	return strAccessToken, strRefreshToken
}

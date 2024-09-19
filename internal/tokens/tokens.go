package tokens

import (
	"context"

	"github.com/FlyKarlik/auth-service/internal/config"
	"github.com/FlyKarlik/auth-service/internal/domain"
	"github.com/FlyKarlik/auth-service/internal/repository"
	authtoken "github.com/FlyKarlik/auth-service/internal/tokens/auth-tokens"
	"github.com/FlyKarlik/auth-service/pkg/logger"
)

type IAuthTokens interface {
	CreateAuthTokens(ctx context.Context, userId string, clientIP string, variety string) (*domain.Token, error)
	ValidateAuthTokens(ctx context.Context, accessToken, refreshToken string) error
	RefreshAuthTokens(ctx context.Context, accessClaims *domain.JWT, refreshClaims *domain.JWT) (*domain.Token, error)
	ParseToken(ctx context.Context, token string) (*domain.JWT, error)
}

type Tokens struct {
	IAuthTokens
}

func New(cfg *config.Config, repo *repository.Repository, log *logger.Logger) *Tokens {
	return &Tokens{
		IAuthTokens: authtoken.New(cfg, repo.IDataRefreshTokenRepository, log),
	}
}

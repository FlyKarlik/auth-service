package usecase

import (
	"context"

	"github.com/FlyKarlik/auth-service/internal/domain"
	"github.com/FlyKarlik/auth-service/internal/tokens"
	"github.com/FlyKarlik/auth-service/internal/usecase/auth"
)

type IAuthUsecase interface {
	Authentication(ctx context.Context, userId string, clientIP string) (*domain.Token, error)
	Refresh(ctx context.Context, clientIP string) (*domain.Token, error)
}

type Usecase struct {
	IAuthUsecase
}

func New(tokens *tokens.Tokens) *Usecase {
	return &Usecase{
		IAuthUsecase: auth.New(tokens),
	}
}

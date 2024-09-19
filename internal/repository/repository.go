package repository

import (
	"context"

	"github.com/FlyKarlik/auth-service/internal/domain"
	"github.com/FlyKarlik/auth-service/internal/repository/postgres"
	"github.com/FlyKarlik/auth-service/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type IDataRefreshTokenRepository interface {
	SaveRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error
	GetRefreshToken(ctx context.Context, id string) (*domain.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error
}

type Repository struct {
	IDataRefreshTokenRepository
}

func New(db *sqlx.DB, log *logger.Logger) *Repository {
	return &Repository{
		IDataRefreshTokenRepository: postgres.New(db, log),
	}
}

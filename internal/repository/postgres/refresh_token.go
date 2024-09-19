package postgres

import (
	"context"
	"time"

	"github.com/FlyKarlik/auth-service/internal/domain"
	"github.com/FlyKarlik/auth-service/internal/errs"
	"github.com/FlyKarlik/auth-service/internal/repository/postgres/queries"
	"github.com/FlyKarlik/auth-service/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenPostgres struct {
	db  *sqlx.DB
	log *logger.Logger
}

func New(db *sqlx.DB, log *logger.Logger) *RefreshTokenPostgres {
	return &RefreshTokenPostgres{db: db, log: log}
}

func (r *RefreshTokenPostgres) SaveRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		queries.CreateRefreshTokenQuery,
		refreshToken.ID, refreshToken.UserID,
		refreshToken.RefreshHash,
		refreshToken.UpdatedAt)
	if err != nil {
		r.log.Errorf("[refreshTokenPostgres.SaveRefreshToken] r.db.ExecContext error: %s", err)
		return errs.ErrDatabaseExecContext
	}

	return nil
}

func (r *RefreshTokenPostgres) GetRefreshToken(ctx context.Context, id string) (*domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var refreshTokenObj domain.RefreshToken

	if err := r.db.GetContext(ctx, &refreshTokenObj, queries.GetRefreshTokenQuery, id); err != nil {
		r.log.Errorf("[refreshTokenPostgres.GetRefreshToken] r.db.GetContext error: %s", err)
		return nil, errs.ErrDatabaseGetContext
	}

	return &refreshTokenObj, nil
}

func (r *RefreshTokenPostgres) UpdateRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := r.db.ExecContext(ctx, queries.UpdateRefreshTokenQuery, refreshToken.RefreshHash, refreshToken.UpdatedAt, refreshToken.ID)
	if err != nil {
		r.log.Errorf("[refreshTokenPostgres.UpdateRefreshToken] r.db.ExecContext error: %s", err)
		return errs.ErrDatabaseExecContext
	}

	return nil
}

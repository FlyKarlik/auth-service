package domain

import (
	"time"
)

type Token struct {
	AccessToken           string `json:"access_token"`
	AccessTokenIssuedAt   int    `json:"access_token_issued_at"`
	AccessTokenExpiresAt  int    `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenIssuedAt  int    `json:"refresh_token_issued_at"`
	RefreshTokenExpiresAt int    `json:"refresh_token_expires_at"`
}

type RefreshToken struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	RefreshHash string    `db:"refresh_hash"`
	UpdatedAt   time.Time `db:"updated_at"`
}

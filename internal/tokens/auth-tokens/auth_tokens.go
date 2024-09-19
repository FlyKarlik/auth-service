package authtoken

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/FlyKarlik/auth-service/internal/config"
	"github.com/FlyKarlik/auth-service/internal/domain"
	"github.com/FlyKarlik/auth-service/internal/errs"
	"github.com/FlyKarlik/auth-service/internal/repository"
	"github.com/FlyKarlik/auth-service/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokens struct {
	cfg  *config.Config
	repo repository.IDataRefreshTokenRepository
	log  *logger.Logger
}

func New(cfg *config.Config, repo repository.IDataRefreshTokenRepository, log *logger.Logger) *AuthTokens {
	return &AuthTokens{
		cfg:  cfg,
		repo: repo,
		log:  log,
	}
}

func (a *AuthTokens) CreateAuthTokens(ctx context.Context, userId string, clientIP string, variety string) (*domain.Token, error) {

	accessClaims := domain.JWT{
		ID:       uuid.NewString(),
		ClientIP: clientIP,
		UserID:   userId,
		Variety:  variety,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(domain.AccessToken).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshClaims := domain.JWT{
		ID:       uuid.NewString(),
		ClientIP: clientIP,
		UserID:   userId,
		Variety:  variety,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(domain.RefreshAccessToken).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims).SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		a.log.Errorf("[AuthTokens.CreateAuthTokens] jwt.NewWithClaims access token error: %s", err)
		return nil, errs.ErrCreateJWT
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims).SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		a.log.Errorf("[AuthTokens.CreateAuthTokens] jwt.NewWithClaims refresh token error: %s", err)
		return nil, errs.ErrCreateJWT
	}

	refreshTokenModel, err := a.beforeSaveRefreshToken(userId, refreshClaims.ID, refreshToken)
	if err != nil {
		return nil, err
	}

	if err := a.repo.SaveRefreshToken(ctx, refreshTokenModel); err != nil {
		return nil, err
	}

	return &domain.Token{
		AccessToken:           accessToken,
		AccessTokenIssuedAt:   int(accessClaims.IssuedAt),
		AccessTokenExpiresAt:  int(accessClaims.ExpiresAt),
		RefreshToken:          refreshToken,
		RefreshTokenIssuedAt:  int(refreshClaims.IssuedAt),
		RefreshTokenExpiresAt: int(refreshClaims.ExpiresAt),
	}, nil
}

func (a *AuthTokens) RefreshAuthTokens(ctx context.Context, accessClaims *domain.JWT, refreshClaims *domain.JWT) (*domain.Token, error) {

	newAccessClaims := domain.JWT{
		ID:       accessClaims.ID,
		ClientIP: accessClaims.ClientIP,
		UserID:   accessClaims.UserID,
		Variety:  accessClaims.Variety,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(domain.AccessToken).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	newRefreshClaims := domain.JWT{
		ID:       refreshClaims.ID,
		ClientIP: refreshClaims.ClientIP,
		UserID:   refreshClaims.UserID,
		Variety:  refreshClaims.Variety,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(domain.RefreshAccessToken).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, newAccessClaims).SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		a.log.Errorf("[authTokens.RefreshAuthTokens] jwt.NewWithClaims access token error: %s", err)
		return nil, errs.ErrCreateJWT
	}

	newRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, newRefreshClaims).SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		a.log.Errorf("[authTokens.RefreshAuthTokens] jwt.NewWithClaims refresh token error: %s", err)
		return nil, errs.ErrCreateJWT
	}

	refreshTokenModel, err := a.beforeSaveRefreshToken(refreshClaims.UserID, refreshClaims.ID, newRefreshToken)
	if err != nil {
		return nil, err
	}

	if err := a.repo.UpdateRefreshToken(ctx, refreshTokenModel); err != nil {
		return nil, err
	}

	return &domain.Token{
		AccessToken:           newAccessToken,
		AccessTokenIssuedAt:   int(newAccessClaims.IssuedAt),
		AccessTokenExpiresAt:  int(newAccessClaims.ExpiresAt),
		RefreshToken:          newRefreshToken,
		RefreshTokenIssuedAt:  int(newRefreshClaims.IssuedAt),
		RefreshTokenExpiresAt: int(newRefreshClaims.ExpiresAt),
	}, nil
}

func (a *AuthTokens) ValidateAuthTokens(ctx context.Context, accessToken, refreshToken string) error {
	accessClaims, err := a.ParseToken(ctx, accessToken)
	if err != nil {
		return err
	}

	refreshClaims, err := a.ParseToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	refreshTokenObj, err := a.repo.GetRefreshToken(ctx, refreshClaims.ID)
	if err != nil {
		return err
	}

	if err := compareHashAndToken(refreshTokenObj.RefreshHash, refreshToken); err != nil {
		a.log.Errorf("[authTokens.ValidateAuthTokens] compareHashAndToken error: %s", err)
		return errs.ErrInvalidToken
	}

	if accessClaims.ClientIP != refreshClaims.ClientIP || accessClaims.UserID != refreshClaims.UserID {
		a.log.Errorf("[authTokens.ValidateAuthTokens] dont compare user id and client ip claims in token")
		return errs.ErrMismatchUserData
	}

	if accessClaims.Variety != refreshClaims.Variety {
		a.log.Errorf("[authTokens.ValidateAuthTokens] dont comapre variety in token")
		return errs.ErrMismatchTokenVariety
	}

	return nil

}

func (a *AuthTokens) ParseToken(ctx context.Context, token string) (*domain.JWT, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &domain.JWT{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.cfg.JWTSecret), nil
	})
	if err != nil {
		a.log.Errorf("[authTokens.ParseToken] jwt.ParseWithClaims error: %s", err)
		return nil, errs.ErrInvalidToken
	}

	if !parsedToken.Valid {
		a.log.Errorf("[authTokens.ParseToken] parsed token valid error")
		return nil, errs.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*domain.JWT)
	if !ok {
		a.log.Errorf("[authTokens.ParseToken] parsedToken.Claims.(*domain.JWT) error: %s", err)
		return nil, errs.ErrInvalidToken
	}

	return claims, nil
}

func (a *AuthTokens) beforeSaveRefreshToken(userId string, refreshTokenId string, refreshToken string) (domain.RefreshToken, error) {

	md5Hash := md5.Sum([]byte(refreshToken))
	hashedRefreshToken := fmt.Sprintf("%x", md5Hash)

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(hashedRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		a.log.Errorf("[AuthTokens.CreateAuthTokens] beforeSaveRefreshToken error: %s", err)
		return domain.RefreshToken{}, errs.ErrCreateJWT
	}

	return domain.RefreshToken{
		ID:          refreshTokenId,
		UserID:      userId,
		RefreshHash: string(refreshTokenHash),
		UpdatedAt:   time.Now(),
	}, nil
}

func compareHashAndToken(hash string, token string) error {
	md5Hash := md5.Sum([]byte(token))
	hashedRefreshToken := fmt.Sprintf("%x", md5Hash)

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(hashedRefreshToken))
}

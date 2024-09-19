package authtoken

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"

	"github.com/FlyKarlik/auth-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthTokens_beforeSaveRefreshToken(t *testing.T) {
	a := AuthTokens{log: logger.NewLogger(logger.LevelDebug)} // Initialize logger if needed

	userId := "test_user_id"
	refreshTokenId := "test_refresh_token_id"
	refreshToken := "test_refresh_token"

	refreshTokenModel, err := a.beforeSaveRefreshToken(userId, refreshTokenId, refreshToken)

	t.Run("ErrorCheck", func(t *testing.T) {
		assert.NoError(t, err, "The function should return a nil error")
	})

	t.Run("IDCheck", func(t *testing.T) {
		assert.Equal(t, refreshTokenId, refreshTokenModel.ID, "The ID should match the provided value")
	})

	t.Run("UserIDCheck", func(t *testing.T) {
		assert.Equal(t, userId, refreshTokenModel.UserID, "The UserID should match the provided value")
	})

	t.Run("UpdatedAtCheck", func(t *testing.T) {
		assert.True(t, time.Since(refreshTokenModel.UpdatedAt) < time.Second, "UpdatedAt should be close to the current time")
	})
}

func Test_compareHashAndToken(t *testing.T) {
	token := "test_refresh_token"

	md5Hash := md5.Sum([]byte(token))
	hashedRefreshToken := fmt.Sprintf("%x", md5Hash)
	hash, err := bcrypt.GenerateFromPassword([]byte(hashedRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Error during hashing: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		err := compareHashAndToken(string(hash), token)
		assert.NoError(t, err, "Comparison should pass successfully")
	})

	t.Run("Failure_WrongToken", func(t *testing.T) {
		err := compareHashAndToken(string(hash), "wrong_token")
		assert.Error(t, err, "Comparison should return an error because tokens do not match")
	})

	t.Run("Failure_WrongHash", func(t *testing.T) {
		wrongHash, err := bcrypt.GenerateFromPassword([]byte("another_token"), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("Error during hashing: %v", err)
		}

		err = compareHashAndToken(string(wrongHash), token)
		assert.Error(t, err, "Comparison should return an error because hashes do not match")
	})
}

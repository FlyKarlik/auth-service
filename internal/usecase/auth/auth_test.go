package auth

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_getAuthTokensFromContext(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, access, "test_access_token")
		ctx = context.WithValue(ctx, refresh, "test_refresh_token")

		accessToken, refreshToken := getAuthTokensFromContext(ctx)

		assert.Equal(t, "test_access_token", accessToken)
		assert.Equal(t, "test_refresh_token", refreshToken)
	})

	t.Run("Negative_TokensMissing", func(t *testing.T) {
		ctx := context.Background()

		accessToken, refreshToken := getAuthTokensFromContext(ctx)

		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
	})

	t.Run("Negative_TokensWrongType", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, access, 123)
		ctx = context.WithValue(ctx, refresh, 456)

		accessToken, refreshToken := getAuthTokensFromContext(ctx)

		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
	})
}

package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func GenerateVariety() (string, error) {
	varietyBytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, varietyBytes); err != nil {
		return "", fmt.Errorf("failed to generate random variety: %w", err)
	}
	return base64.StdEncoding.EncodeToString(varietyBytes), nil
}

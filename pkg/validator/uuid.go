package validator

import (
	"fmt"

	"github.com/google/uuid"
)

func IsValidStringUUID(uuidStr string) error {
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("cannot parse string uuid: %w", err)
	}
	return nil
}

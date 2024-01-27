package usecase

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		password := "password123"
		hashedPassword, err := HashPassword(password)

		assert.NoError(t, err)

		// The hashed password should be a hex-encoded string with a '.' separator
		split := strings.Split(hashedPassword, ".")
		assert.Equal(t, 2, len(split))

		// Both parts of the hashed password should be valid hex strings
		_, err = hex.DecodeString(split[0])
		assert.NoError(t, err)
		_, err = hex.DecodeString(split[1])
		assert.NoError(t, err)
	})

	t.Run("Empty password", func(t *testing.T) {
		password := ""
		hashedPassword, err := HashPassword(password)

		assert.NoError(t, err)

		// The hashed password should still be a valid hex-encoded string with a '.' separator
		split := strings.Split(hashedPassword, ".")
		assert.Equal(t, 2, len(split))

		// Both parts of the hashed password should be valid hex strings
		_, err = hex.DecodeString(split[0])
		assert.NoError(t, err)
		_, err = hex.DecodeString(split[1])
		assert.NoError(t, err)
	})
}

package apperror

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Status(t *testing.T) {
	t.Run("Authorization", func(t *testing.T) {
		err := &Error{Type: Authorization}
		assert.Equal(t, http.StatusUnauthorized, err.Status())
	})

	t.Run("BadRequest", func(t *testing.T) {
		err := &Error{Type: BadRequest}
		assert.Equal(t, http.StatusBadRequest, err.Status())
	})

	t.Run("Conflict", func(t *testing.T) {
		err := &Error{Type: Conflict}
		assert.Equal(t, http.StatusConflict, err.Status())
	})

	t.Run("Internal", func(t *testing.T) {
		err := &Error{Type: Internal}
		assert.Equal(t, http.StatusInternalServerError, err.Status())
	})

	t.Run("NotFound", func(t *testing.T) {
		err := &Error{Type: NotFound}
		assert.Equal(t, http.StatusNotFound, err.Status())
	})

	t.Run("PayloadTooLarge", func(t *testing.T) {
		err := &Error{Type: PayloadTooLarge}
		assert.Equal(t, http.StatusRequestEntityTooLarge, err.Status())
	})

	t.Run("UnsupportedMediaType", func(t *testing.T) {
		err := &Error{Type: UnsupportedMediaType}
		assert.Equal(t, http.StatusUnsupportedMediaType, err.Status())
	})

	t.Run("Default", func(t *testing.T) {
		err := &Error{} // some type that is not defined
		assert.Equal(t, http.StatusInternalServerError, err.Status())
	})
}

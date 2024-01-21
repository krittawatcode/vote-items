package helper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend/domain"
	"github.com/stretchr/testify/assert"
)

func TestBindData(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"email":"bob@bob.com","password":"password123"}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/test", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		var user domain.User
		result := BindData(c, &user)

		assert.True(t, result)
		assert.Equal(t, "bob@bob.com", user.Email)
		assert.Equal(t, "", user.Password) // password should be empty
	})

	t.Run("Invalid request body - non-json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`bob@bob.com`) // non-json
		c.Request = httptest.NewRequest(http.MethodPost, "/test", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		var user domain.User
		result := BindData(c, &user)

		assert.False(t, result)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

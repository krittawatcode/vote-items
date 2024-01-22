package route

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Health Check", func(t *testing.T) {
		r := SetupRouter()
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, os.Getenv("USER_SERVICE_API_URL")+"/health", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User service is running")
	})
}

package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krittawatcode/vote-items/backend/delivery/routes"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	router := routes.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"hello":"it's me"}`, w.Body.String())
}

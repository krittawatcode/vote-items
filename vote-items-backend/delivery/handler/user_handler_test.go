package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend/delivery/handler"
	"github.com/stretchr/testify/assert"
)

func TestMeHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	new(handler.UserHandler).Me(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"hello":"it's me"}`, w.Body.String())
}

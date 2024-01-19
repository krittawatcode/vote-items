package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend/delivery/handler"
	"github.com/krittawatcode/vote-items/backend/domain/mocks"
	"github.com/krittawatcode/vote-items/backend/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMeHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("No user in context", func(t *testing.T) {
		new(handler.UserHandler).Me(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"Internal server error"}`, w.Body.String())
	})

	t.Run("UserUseCase.Get returns error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.AnythingOfType("*gin.Context"), uid).Return(nil, errors.New("error"))

		c.Set("user", &model.User{
			UID: uid,
		})

		h := &handler.UserHandler{
			UserUseCase: mockUserService,
		}
		h.Me(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"Not found: user"}`, w.Body.String())
	})

	//TODO: add not found from user usecase

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid,
			Email: "bob@bob.com",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.AnythingOfType("*gin.Context"), uid).Return(mockUserResp, nil)

		c.Set("user", &model.User{
			UID: uid,
		})

		h := &handler.UserHandler{
			UserUseCase: mockUserService,
		}
		h.Me(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"user":{"UID":"`+uid.String()+`","Email":"bob@bob.com"}}`, w.Body.String())
	})
}

package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend/domain"
	"github.com/krittawatcode/vote-items/backend/domain/appmock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Me(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		uid, _ := uuid.NewRandom()
		mockUser := &domain.User{
			UID: uid,
		}

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("Get", mock.Anything, uid).Return(mockUser, nil)

		h := &UserHandler{
			UserUseCase: mockUserUseCase,
		}

		c.Set("user", mockUser)
		h.Me(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		uid, _ := uuid.NewRandom()
		mockUser := &domain.User{
			UID: uid,
		}

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("Get", mock.Anything, uid).Return(nil, errors.New("user not found"))

		h := &UserHandler{
			UserUseCase: mockUserUseCase,
		}

		c.Set("user", mockUser)
		h.Me(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("User not set in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		h := &UserHandler{
			UserUseCase: new(appmock.MockUserUseCase),
		}

		h.Me(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHandler_SignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"email":"bob@bob.com","password":"password123"}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/signUp", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignUp", mock.Anything, mock.Anything).Return(nil)

		h := &UserHandler{
			UserUseCase: mockUserUseCase,
		}
		h.SignUp(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("UserUseCase.SignUp returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"email":"bob@bob.com","password":"password123"}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/signUp", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignUp", mock.Anything, mock.Anything).Return(errors.New("error"))

		h := &UserHandler{
			UserUseCase: mockUserUseCase,
		}
		h.SignUp(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UserUseCase.SignUp returns user already exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := strings.NewReader(`{"email":"bob@bob.com","password":"password123"}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/signUp", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignUp", mock.AnythingOfType("*gin.Context"), mock.Anything).Return(domain.NewConflict("User Already Exists", "test email"))

		h := &UserHandler{
			UserUseCase: mockUserUseCase,
		}
		h.SignUp(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockUserUseCase.AssertExpectations(t)
	})

	t.Run("Invalid request body - missing password", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"email":"bob@bob.com"}`)              // missing password
		c.Request = httptest.NewRequest(http.MethodPost, "/signUp", reqBody) // can use any path here
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignUp", mock.Anything, mock.Anything).Return(nil, errors.New("invalid request body"))

		h := &UserHandler{
			UserUseCase: mockUserUseCase,
		}
		h.SignUp(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUserUseCase.AssertNotCalled(t, "SignUp")
	})

	t.Run("Invalid request body - missing email", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"email":""}`)                         // missing email
		c.Request = httptest.NewRequest(http.MethodPost, "/signUp", reqBody) // can use any path here
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignUp", mock.Anything, mock.Anything).Return(nil, errors.New("invalid request body"))

		h := &UserHandler{
			UserUseCase: mockUserUseCase,
		}
		h.SignUp(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUserUseCase.AssertNotCalled(t, "SignUp")
	})
}

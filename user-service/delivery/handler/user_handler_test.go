package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/user-service/domain"
	"github.com/krittawatcode/vote-items/user-service/domain/apperror"
	"github.com/krittawatcode/vote-items/user-service/domain/appmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Me(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		uid, _ := uuid.NewRandom()
		mockUser := &domain.User{
			UID: uid,
		}

		// Define a custom context key type to prevent collisions with other packages using context keys
		type contextKey string
		const key contextKey = "user"
		c.Set(string(key), mockUser)
		// Set up a request context and add it to the Gin context
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), key, mockUser))

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("Get", mock.Anything, uid).Return(mockUser, nil)
		mockTokenUseCase := new(appmock.MockTokenUseCase)

		h := &UserHandler{
			Router:       gin.Default(),
			UserUseCase:  mockUserUseCase,
			TokenUseCase: mockTokenUseCase,
		}

		h.Me(c)

		assert.Equal(t, http.StatusOK, w.Code)
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

		mockTokenUseCase := new(appmock.MockTokenUseCase)
		mockTokenUseCase.On("NewPairFromUser", mock.Anything, mock.Anything, "").Return(&domain.TokenPair{}, nil)

		h := &UserHandler{
			UserUseCase:  mockUserUseCase,
			TokenUseCase: mockTokenUseCase,
		}
		h.SignUp(c)

		assert.Equal(t, http.StatusCreated, w.Code)
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
		mockUserUseCase.On("SignUp", mock.Anything, mock.Anything).Return(apperror.NewConflict("User Already Exists", "test email"))

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
	t.Run("TokenUseCase.NewPairFromUser returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"email":"bob@bob.com","password":"password123"}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/signUp", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignUp", mock.Anything, mock.Anything).Return(nil)

		mockTokenUseCase := new(appmock.MockTokenUseCase)
		mockTokenUseCase.On("NewPairFromUser", mock.Anything, mock.Anything, "").Return(nil, errors.New("error"))

		h := &UserHandler{
			UserUseCase:  mockUserUseCase,
			TokenUseCase: mockTokenUseCase,
		}
		h.SignUp(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

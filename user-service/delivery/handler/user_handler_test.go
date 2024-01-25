package handler

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestSignIn(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Bad request data", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)

		reqBody := strings.NewReader(`{"email":"bob@bob.com","password":"short"}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/signIp", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockTokenUseCase := new(appmock.MockTokenUseCase)

		h := &UserHandler{
			UserUseCase:  mockUserUseCase,
			TokenUseCase: mockTokenUseCase,
		}

		h.SignIn(c)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserUseCase.AssertNotCalled(t, "SignIn")
		mockTokenUseCase.AssertNotCalled(t, "NewTokensFromUser")
	})

	t.Run("Error Returned from UserUseCase.SignIn", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)

		email := "bob@bob.com"
		password := "pwdoesnotmatch123"

		mockUSArgs := mock.Arguments{
			mock.Anything,
			&domain.User{Email: email, Password: password},
		}

		// so we can check for a known status code
		mockError := apperror.NewAuthorization("invalid email/password combo")

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignIn", mockUSArgs...).Return(mockError)

		mockTokenUseCase := new(appmock.MockTokenUseCase)

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/signIn", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		h := &UserHandler{
			UserUseCase:  mockUserUseCase,
			TokenUseCase: mockTokenUseCase,
		}

		h.SignIn(c)

		mockUserUseCase.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTokenUseCase.AssertNotCalled(t, "NewTokensFromUser")
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)

		email := "bob@bob.com"
		password := "pwworksgreat123"

		mockUSArgs := mock.Arguments{
			mock.Anything,
			&domain.User{Email: email, Password: password},
		}

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignIn", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.Anything,
			&domain.User{Email: email, Password: password},
			"",
		}

		mockTokenPair := &domain.TokenPair{
			IDToken:      "idToken",
			RefreshToken: "refreshToken",
		}

		mockTokenUseCase := new(appmock.MockTokenUseCase)
		mockTokenUseCase.On("NewPairFromUser", mockTSArgs...).Return(mockTokenPair, nil)

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		c.Request = httptest.NewRequest(http.MethodPost, "/signIn", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		h := &UserHandler{
			UserUseCase:  mockUserUseCase,
			TokenUseCase: mockTokenUseCase,
		}

		h.SignIn(c)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenPair,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserUseCase.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTokenUseCase.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)

		email := "cannotproducetoken@bob.com"
		password := "cannotproducetoken"

		mockUSArgs := mock.Arguments{
			mock.Anything,
			&domain.User{Email: email, Password: password},
		}

		mockUserUseCase := new(appmock.MockUserUseCase)
		mockUserUseCase.On("SignIn", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.Anything,
			&domain.User{Email: email, Password: password},
			"",
		}

		mockError := apperror.NewInternal()

		mockTokenUseCase := new(appmock.MockTokenUseCase)
		mockTokenUseCase.On("NewPairFromUser", mockTSArgs...).Return(nil, mockError)

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		c.Request = httptest.NewRequest(http.MethodPost, "/signIn", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		h := &UserHandler{
			UserUseCase:  mockUserUseCase,
			TokenUseCase: mockTokenUseCase,
		}

		h.SignIn(c)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserUseCase.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTokenUseCase.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})
}

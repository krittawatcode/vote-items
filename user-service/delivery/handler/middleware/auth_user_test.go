package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/user-service/domain"
	"github.com/krittawatcode/vote-items/user-service/domain/apperror"
	"github.com/krittawatcode/vote-items/user-service/domain/appmock"

	"github.com/stretchr/testify/assert"
)

func TestAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTokenUseCase := new(appmock.MockTokenUseCase)

	uid, _ := uuid.NewRandom()
	u := &domain.User{
		UID:   uid,
		Email: "bob@bob.com",
	}

	// Since we mock tokenService, we need not
	// create actual JWTs
	validTokenHeader := "validTokenString"
	invalidTokenHeader := "invalidTokenString"
	invalidTokenErr := apperror.NewAuthorization("Unable to verify user from idToken")

	mockTokenUseCase.On("ValidateIDToken", validTokenHeader).Return(u, nil)
	mockTokenUseCase.On("ValidateIDToken", invalidTokenHeader).Return(nil, invalidTokenErr)

	t.Run("Adds a user to context", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// creates a test context and gin engine
		_, r := gin.CreateTestContext(rr)

		// will be populated with user in a handler
		// if AuthUser middleware is successful
		var contextUser *domain.User

		// see this issue - https://github.com/gin-gonic/gin/issues/323
		// https://github.com/gin-gonic/gin/blob/master/auth_test.go#L91-L126
		// we create a handler to return "user added to context" as this
		// is the only way to test modified context
		r.GET("/me", AuthUser(mockTokenUseCase), func(c *gin.Context) {
			contextKeyVal, _ := c.Get("user")
			contextUser = contextKeyVal.(*domain.User)
		})

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validTokenHeader))
		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, u, contextUser)

		mockTokenUseCase.AssertCalled(t, "ValidateIDToken", validTokenHeader)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// creates a test context and gin engine
		_, r := gin.CreateTestContext(rr)

		r.GET("/me", AuthUser(mockTokenUseCase))

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidTokenHeader))
		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockTokenUseCase.AssertCalled(t, "ValidateIDToken", invalidTokenHeader)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// creates a test context and gin engine
		_, r := gin.CreateTestContext(rr)

		r.GET("/me", AuthUser(mockTokenUseCase))

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockTokenUseCase.AssertNotCalled(t, "ValidateIDToken")
	})
}

package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valide token", func(t *testing.T) {
		authServiceMock := new(AuthServiceMock)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   "1",
			"exp":  time.Now().Add(time.Second * 5).Unix(),
			"role": "moderator",
		})
		token.Valid = true

		authServiceMock.On("ValidateToken", "fake_token").Return(token, nil)

		router := gin.New()
		router.Use(JWTAuthMiddleware(authServiceMock))
		router.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer fake_token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		authServiceMock.AssertCalled(t, "ValidateToken", "fake_token")
	})

	t.Run("authorization header is required", func(t *testing.T) {
		authServiceMock := new(AuthServiceMock)

		router := gin.New()
		router.Use(JWTAuthMiddleware(authServiceMock))
		router.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Authorization header is required")
	})

	t.Run("invalid token format", func(t *testing.T) {
		authServiceMock := new(AuthServiceMock)

		router := gin.New()
		router.Use(JWTAuthMiddleware(authServiceMock))
		router.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "fake_token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid token format")
	})

}

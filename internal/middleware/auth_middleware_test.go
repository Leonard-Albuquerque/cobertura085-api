package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/middleware"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "super_secret_for_middleware_test"
	jwtSvc := service.NewJWTService(secret)

	r := gin.New()
	r.Use(middleware.AuthMiddleware(jwtSvc))
	r.GET("/protected", func(c *gin.Context) {
		userID, _ := middleware.GetUserID(c)
		storeID, _ := middleware.GetStoreID(c)
		email, _ := middleware.GetUserEmail(c)

		c.JSON(http.StatusOK, gin.H{
			"userID":  userID,
			"storeID": storeID,
			"email":   email,
		})
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Status esperado: %d, obteve: %d", http.StatusUnauthorized, resp.Code)
		}
	})

	t.Run("Invalid Header Format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidHeaderFormat")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Status esperado: %d, obteve: %d", http.StatusUnauthorized, resp.Code)
		}
	})

	t.Run("Valid Token Passed", func(t *testing.T) {
		token, err := jwtSvc.GenerateAccessToken("user-1", "store-1", "user@test.com", 15*time.Minute)
		if err != nil {
			t.Fatalf("Erro ao gerar token: %v", err)
		}

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Status esperado: %d, obteve: %d", http.StatusOK, resp.Code)
		}
	})
}

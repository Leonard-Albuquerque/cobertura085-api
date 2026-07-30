package service_test

import (
	"testing"
	"time"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
)

func TestJWTService(t *testing.T) {
	secret := "test_secret_key_12345"
	jwtSvc := service.NewJWTService(secret)

	userID := "user-123"
	storeID := "store-456"
	email := "test@example.com"

	t.Run("Generate and Validate Valid Access Token", func(t *testing.T) {
		tokenStr, err := jwtSvc.GenerateAccessToken(userID, storeID, email, 15*time.Minute)
		if err != nil {
			t.Fatalf("Esperava sucesso ao gerar token, obteve erro: %v", err)
		}

		claims, err := jwtSvc.ValidateAccessToken(tokenStr)
		if err != nil {
			t.Fatalf("Esperava sucesso ao validar token, obteve erro: %v", err)
		}

		if claims.UserID != userID {
			t.Errorf("UserID esperado: %s, obteve: %s", userID, claims.UserID)
		}
		if claims.StoreID != storeID {
			t.Errorf("StoreID esperado: %s, obteve: %s", storeID, claims.StoreID)
		}
		if claims.Email != email {
			t.Errorf("Email esperado: %s, obteve: %s", email, claims.Email)
		}
	})

	t.Run("Validate Expired Access Token", func(t *testing.T) {
		tokenStr, err := jwtSvc.GenerateAccessToken(userID, storeID, email, -1*time.Minute)
		if err != nil {
			t.Fatalf("Erro ao gerar token expirado: %v", err)
		}

		_, err = jwtSvc.ValidateAccessToken(tokenStr)
		if err == nil {
			t.Error("Esperava erro ao validar token expirado, mas foi validado com sucesso")
		}
	})

	t.Run("Password Hashing and Verification", func(t *testing.T) {
		password := "my_secure_password"
		hash, err := jwtSvc.HashPassword(password)
		if err != nil {
			t.Fatalf("Erro ao gerar hash da senha: %v", err)
		}

		if !jwtSvc.CheckPassword(password, hash) {
			t.Error("Esperava sucesso ao checar senha correta")
		}

		if jwtSvc.CheckPassword("wrong_password", hash) {
			t.Error("Esperava falha ao checar senha incorreta")
		}
	})

	t.Run("Token Hashing Consistency", func(t *testing.T) {
		token := "some_random_refresh_token_string"
		hash1 := jwtSvc.HashToken(token)
		hash2 := jwtSvc.HashToken(token)

		if hash1 != hash2 {
			t.Errorf("Hashes do token diferem: %s vs %s", hash1, hash2)
		}
	})
}

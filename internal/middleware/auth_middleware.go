package middleware

import (
	"net/http"
	"strings"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey    = "userID"
	ContextStoreIDKey   = "storeID"
	ContextUserEmailKey = "userEmail"
	ContextClaimsKey    = "userClaims"
)

func respondUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
}

// AuthMiddleware valida o header Authorization (Bearer JWT) e armazena os dados do usuário no contexto.
func AuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondUnauthorized(c, "Cabeçalho Authorization não fornecido")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			respondUnauthorized(c, "Formato do token deve ser 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			respondUnauthorized(c, "Token de acesso inválido ou expirado")
			c.Abort()
			return
		}

		// Armazena as informações no contexto do Gin
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextStoreIDKey, claims.StoreID)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Set(ContextClaimsKey, claims)

		c.Next()
	}
}

// GetUserID extrai o ID do usuário autenticado do contexto Gin.
func GetUserID(c *gin.Context) (string, bool) {
	val, ok := c.Get(ContextUserIDKey)
	if !ok {
		return "", false
	}
	id, ok := val.(string)
	return id, ok
}

// GetStoreID extrai o ID da loja associada do contexto Gin.
func GetStoreID(c *gin.Context) (string, bool) {
	val, ok := c.Get(ContextStoreIDKey)
	if !ok {
		return "", false
	}
	id, ok := val.(string)
	return id, ok
}

// GetUserEmail extrai o e-mail do usuário autenticado do contexto Gin.
func GetUserEmail(c *gin.Context) (string, bool) {
	val, ok := c.Get(ContextUserEmailKey)
	if !ok {
		return "", false
	}
	email, ok := val.(string)
	return email, ok
}

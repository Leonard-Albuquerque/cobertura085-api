package handler

import (
	"errors"
	"net/http"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/middleware"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register gerencia o cadastro de um novo usuário para uma loja.
func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrStoreNameAlreadyExists) || errors.Is(err, service.ErrEmailAlreadyInUse) {
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, user)
}

// Login realiza a autenticação do usuário e retorna Access Token + Refresh Token.
func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	if input.UserAgent == nil {
		ua := c.GetHeader("User-Agent")
		if ua != "" {
			input.UserAgent = &ua
		}
	}
	if input.IP == nil {
		ip := c.ClientIP()
		if ip != "" {
			input.IP = &ip
		}
	}

	authRes, err := h.authService.Login(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrUserInactive) {
			RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, authRes)
}

// RefreshToken renova o Access Token e realiza rotação do Refresh Token.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input service.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	if input.UserAgent == nil {
		ua := c.GetHeader("User-Agent")
		if ua != "" {
			input.UserAgent = &ua
		}
	}
	if input.IP == nil {
		ip := c.ClientIP()
		if ip != "" {
			input.IP = &ip
		}
	}

	authRes, err := h.authService.RefreshToken(c.Request.Context(), input)
	if err != nil {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, authRes)
}

// Logout revoga o Refresh Token fornecido.
func (h *AuthHandler) Logout(c *gin.Context) {
	var input service.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	_ = h.authService.Logout(c.Request.Context(), input.RefreshToken)

	RespondSuccess(c, http.StatusOK, gin.H{"message": "Sessão encerrada com sucesso"})
}

// GetProfile retorna os dados do usuário autenticado no contexto.
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Usuário não autenticado")
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Perfil do usuário não encontrado")
		return
	}

	RespondSuccess(c, http.StatusOK, user)
}

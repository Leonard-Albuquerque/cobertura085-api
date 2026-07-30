package handler

import (
	"github.com/gin-gonic/gin"
)

// Response representa o envelope padrão para respostas HTTP de sucesso.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorDetails armazena detalhes adicionais sobre um erro ocorrido.
type ErrorDetails struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// ErrorResponse representa o envelope padrão para respostas HTTP de erro.
type ErrorResponse struct {
	Success bool         `json:"success"`
	Error   ErrorDetails `json:"error"`
}

// RespondSuccess envia uma resposta JSON de sucesso consistente.
func RespondSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// RespondError envia uma resposta JSON de erro consistente.
func RespondError(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error: ErrorDetails{
			Code:    code,
			Message: message,
		},
	})
}

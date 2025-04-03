package controller

import (
	"message-server/internal/domain"
	"message-server/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase usecases.AuthUseCase
}

func NewAuthHandler(authUseCase *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: *authUseCase}
}

func (s *AuthHandler) Register(c *gin.Context) {
	var request domain.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := s.authUseCase.Register(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (s *AuthHandler) Login(c *gin.Context) {
	var request domain.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := s.authUseCase.Login(&request)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.SetCookie("auth_token", token, 60*60*24*30, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

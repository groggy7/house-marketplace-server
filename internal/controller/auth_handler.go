package controller

import (
	"message-server/internal/controller/auth"
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
		switch err {
		case domain.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrDuplicateUsername:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error(), "field": "username"})
		case domain.ErrDuplicateEmail:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error(), "field": "email"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
		}
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

	user, token, err := s.authUseCase.Login(&request)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		case domain.ErrUserNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("auth_token", token, 24*60*60*30, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"full_name":  user.FullName,
			"email":      user.Email,
			"avatar_key": user.AvatarKey,
		},
	})
}

func (s *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (s *AuthHandler) CheckIsLoggedIn(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username := claims.(*auth.Claims).Username
	user, err := s.authUseCase.GetUserByUsername(username)
	if err != nil {
		c.SetCookie("auth_token", "", -1, "/", "", true, true)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User is logged in", "user": user})
}

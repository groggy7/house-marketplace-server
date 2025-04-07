package controller

import (
	"message-server/internal/controller/auth"
	"message-server/internal/domain"
	"message-server/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase *usecases.UserUseCase
}

func NewUserHandler(userUseCase *usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (s *UserHandler) UpdateUserInfo(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := claims.(*auth.Claims).UserID
	email := claims.(*auth.Claims).Email

	var request domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	request.UserID = userID
	request.Email = email

	err := s.userUseCase.UpdateUserInfo(&request)
	if err != nil {
		switch err {
		case domain.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (s *UserHandler) UpdateUserAvatar(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := claims.(*auth.Claims).UserID
	email := claims.(*auth.Claims).Email

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	request := &domain.UpdateUserRequest{
		UserID: userID,
		Email:  email,
	}

	err = s.userUseCase.UpdateUserAvatar(request, file)
	if err != nil {
		switch err {
		case domain.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully"})
}

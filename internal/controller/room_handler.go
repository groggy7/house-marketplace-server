package controller

import (
	"message-server/internal/controller/auth"
	"message-server/internal/domain"
	"message-server/internal/usecases"
	"message-server/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	roomUseCase usecases.RoomUseCase
}

func InitRoomServer(svc *usecases.RoomUseCase) *ChatHandler {
	return &ChatHandler{
		roomUseCase: *svc,
	}
}

func (s *ChatHandler) CreateRoom(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Claims not found"})
		return
	}

	user := claims.(*auth.Claims)

	var request domain.CreateChatRoomRequest
	request.CustomerID = user.UserID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.PropertyID == "" || request.PropertyOwnerID == "" || request.CustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	roomID, err := s.roomUseCase.CreateRoom(&request)
	if err != nil {
		pkg.Logger.Printf("Failed to create chat room: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room_id": roomID})
}

func (s *ChatHandler) GetRooms(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Claims not found"})
		return
	}

	user := claims.(*auth.Claims)

	rooms, err := s.roomUseCase.GetRooms(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (s *ChatHandler) GetRoomMessages(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	messages, err := s.roomUseCase.GetMessagesForRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat rooms"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

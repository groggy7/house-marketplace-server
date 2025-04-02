package controller

import (
	"message-server/internal/room"
	"message-server/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomServer struct {
	roomService room.RoomService
}

func InitRoomServer(roomService room.RoomService) *RoomServer {
	return &RoomServer{
		roomService: roomService,
	}
}

func (s *RoomServer) CreateRoom(c *gin.Context) {
	var request room.CreateChatRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.PropertyID == "" || request.PropertyOwnerID == "" || request.CustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	roomID, err := s.roomService.CreateRoom(&request)
	if err != nil {
		pkg.Logger.Printf("Failed to create chat room: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room_id": roomID})
}

func (s *RoomServer) GetRooms(c *gin.Context) {
	customerID := c.Param("customer_id")

	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	rooms, err := s.roomService.GetRooms(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (s *RoomServer) GetRoomMessages(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	messages, err := s.roomService.GetMessagesForRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat rooms"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

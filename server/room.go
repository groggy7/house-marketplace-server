package server

import (
	"message-server/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateChatRoomRequest struct {
	PropertyID      string `json:"property_id"`
	PropertyOwnerID string `json:"property_owner_id"`
	CustomerID      string `json:"customer_id"`
}

type GetRoomsRequest struct {
	CustomerID string `json:"customer_id"`
}

type RoomServer struct {
	db *db.ChatDB
}

func InitRoomServer(db *db.ChatDB) *RoomServer {
	return &RoomServer{
		db: db,
	}
}

func (s *RoomServer) CreateRoom(c *gin.Context) {
	var request CreateChatRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.PropertyID == "" || request.PropertyOwnerID == "" || request.CustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	roomID, err := s.db.CreateRoom(request.PropertyID, request.PropertyOwnerID, request.CustomerID)
	if err != nil {
		logger.Printf("Failed to create chat room: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room_id": roomID})
}

func (s *RoomServer) GetRooms(c *gin.Context) {
	var request GetRoomsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.CustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	rooms, err := s.db.GetRooms(request.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (s *RoomServer) GetRoomMessages(c *gin.Context) {
	roomID := c.Query("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	messages, err := s.db.GetMessagesForRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat rooms"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

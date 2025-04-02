package controller

import (
	"fmt"
	"io"

	"message-server/internal/domain"
	"message-server/internal/usecases"
	"message-server/pkg"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type MessageServer struct {
	roomUseCase usecases.RoomUseCase
	authUseCase usecases.AuthUseCase
	upgrader    websocket.Upgrader
	clients     map[string]*websocket.Conn
	mutex       sync.RWMutex
}

func InitMessageServer(svc *usecases.RoomUseCase, authSvc *usecases.AuthUseCase) *MessageServer {
	godotenv.Load()
	frontURL := os.Getenv("FRONTEND_URL")

	if frontURL == "" {
		pkg.Logger.Println("Warning: FRONTEND_URL not set in .env")
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == frontURL || origin == "test_client" {
				return true
			}
			return false
		},
	}

	return &MessageServer{
		roomUseCase: *svc,
		authUseCase: *authSvc,
		upgrader:    upgrader,
		clients:     make(map[string]*websocket.Conn),
	}
}

func (s *MessageServer) StartWebSocketServer(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		pkg.Logger.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	var authMessage domain.AuthMessage
	if err := conn.ReadJSON(&authMessage); err != nil {
		pkg.Logger.Println("Failed to read auth message:", err)
		conn.WriteJSON(domain.MessageResponse{
			Type:  "error",
			Error: "Authentication failed: " + err.Error(),
		})
		conn.Close()
		return
	}

	if !validateAuthMessage(&authMessage, &s.authUseCase) {
		pkg.Logger.Println("Invalid auth message:", authMessage)
		conn.WriteJSON(domain.MessageResponse{
			Type:  "error",
			Error: "Invalid authentication message",
		})
		conn.Close()
		return
	}

	userID := authMessage.UserID

	var oldConn *websocket.Conn
	var exists bool

	s.mutex.Lock()
	oldConn, exists = s.clients[userID]
	s.clients[userID] = conn
	s.mutex.Unlock()

	if exists {
		pkg.Logger.Printf("User %s already has an active connection, closing old one", userID)
		oldConn.WriteJSON(domain.MessageResponse{
			Type:   "disconnect",
			Status: "replaced",
			Error:  "New connection established from another device",
		})
		oldConn.Close()
	}

	conn.WriteJSON(domain.MessageResponse{
		Type:      "auth_success",
		Status:    "connected",
		Timestamp: time.Now().Unix(),
	})
	pkg.Logger.Printf("User %s connected", userID)

	go s.ping(conn, userID)

	s.handleMessages(conn, userID)
}

func (s *MessageServer) ping(conn *websocket.Conn, userID string) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for range ticker.C {
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			s.mutex.Lock()
			if s.clients[userID] == conn {
				delete(s.clients, userID)
			}
			s.mutex.Unlock()
			return
		}
	}
}

func (s *MessageServer) handleMessages(conn *websocket.Conn, senderID string) {
	defer func() {
		s.mutex.Lock()
		if s.clients[senderID] == conn {
			delete(s.clients, senderID)
		}
		s.mutex.Unlock()

		conn.Close()
		pkg.Logger.Printf("User %s disconnected", senderID)
	}()

	for {
		var message domain.ChatMessage
		if err := conn.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				pkg.Logger.Println("WebSocket closed:", err)
				return
			}

			if _, ok := err.(*websocket.CloseError); ok {
				pkg.Logger.Println("Connection closed by client")
				return
			}

			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") ||
				strings.Contains(err.Error(), "unexpected EOF") {
				pkg.Logger.Println("Connection appears to be closed:", err)
				return
			}

			pkg.Logger.Println("Temporary read error:", err)
			continue
		}

		if err := validateChatMessage(&message); err != nil {
			pkg.Logger.Println("Invalid message format:", message)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     err.Error(),
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		if message.SenderID != "" && message.SenderID != senderID {
			pkg.Logger.Printf("Message sender ID mismatch: auth=%s, message=%s", senderID, message.SenderID)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     "Sender ID in message doesn't match authenticated user",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		if message.SenderID == "" {
			message.SenderID = senderID
		}

		exists, err := s.roomUseCase.CheckRoomExists(message.RoomID)
		if err != nil {
			pkg.Logger.Println("Failed to check room existence:", err)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     "Database error when validating room",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		if !exists {
			pkg.Logger.Println("Room does not exist:", message.RoomID)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     "Room does not exist",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		isMember, err := s.roomUseCase.CheckUserInRoom(senderID, message.RoomID)
		if err != nil {
			pkg.Logger.Println("Failed to check room membership:", err)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     "Database error when validating room membership",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		if !isMember {
			pkg.Logger.Printf("User %s is not a member of room %s", senderID, message.RoomID)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     "You are not a member of this room",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		receiverIsMember, err := s.roomUseCase.CheckUserInRoom(message.ReceiverID, message.RoomID)
		if err != nil {
			pkg.Logger.Println("Failed to check receiver room membership:", err)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     "Database error when validating receiver room membership",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		if !receiverIsMember {
			pkg.Logger.Printf("Receiver %s is not a member of room %s", message.ReceiverID, message.RoomID)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     "Receiver is not a member of this room",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		timestamp := time.Now().Unix()
		if err := s.roomUseCase.SaveMessage(message.Text, senderID, message.RoomID); err != nil {
			pkg.Logger.Printf("Error saving message to database: %v", err)
			s.writeJSON(conn, domain.MessageResponse{
				Type:      "error",
				Error:     "Failed to save message",
				Timestamp: timestamp,
			})
			continue
		}

		delivered := s.sendMessage(senderID, message.ReceiverID, message.Text, message.RoomID)

		status := "sent"
		if delivered {
			status = "delivered"
		}

		s.writeJSON(conn, domain.MessageResponse{
			Type:      "status",
			Status:    status,
			Text:      message.Text,
			Timestamp: timestamp,
		})
	}
}

func (s *MessageServer) writeJSON(conn *websocket.Conn, message interface{}) bool {
	if conn == nil {
		return false
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteJSON(message); err != nil {
		pkg.Logger.Printf("Error writing to WebSocket: %v", err)
		return false
	}
	return true
}

func (s *MessageServer) sendMessage(senderID, receiverID, text, roomID string) bool {
	s.mutex.RLock()
	receiverConn, receiverConnected := s.clients[receiverID]
	s.mutex.RUnlock()

	if receiverConnected {
		response := domain.MessageResponse{
			Type:      "message",
			Text:      text,
			SenderID:  senderID,
			RoomID:    roomID,
			Timestamp: time.Now().Unix(),
		}

		s.writeJSON(receiverConn, response)
	}

	pkg.Logger.Printf("Recipient %s not connected or delivery failed, message saved to DB only", receiverID)
	return false
}

func validateAuthMessage(authMessage *domain.AuthMessage, authUseCase *usecases.AuthUseCase) bool {
	if authMessage.Type != "auth" || authMessage.UserID == "" {
		return false
	}

	exists, err := authUseCase.CheckUserExists(authMessage.UserID)
	if err != nil {
		return false
	}

	return exists
}

func validateChatMessage(message *domain.ChatMessage) error {
	if message.Text == "" {
		return fmt.Errorf("message text cannot be empty")
	}

	if len(message.Text) > 5000 {
		return fmt.Errorf("message text is too long (max 5000 characters)")
	}

	if message.ReceiverID == "" {
		return fmt.Errorf("receiver ID cannot be empty")
	}

	if message.RoomID == "" {
		return fmt.Errorf("room ID cannot be empty")
	}

	return nil
}

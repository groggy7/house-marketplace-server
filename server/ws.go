package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var logger = log.New(os.Stdout, "chat server - ", log.Lshortfile)

type MessageServer struct {
	pool *pgxpool.Pool
	upgrader websocket.Upgrader
	senderID string
}

type AuthMessage struct {
    Type string `json:"type"`
    UserID string `json:"user_id"`
}

type Message struct {
	Text string `json:"text"`
	ReceiverID string `json:"receiver_id"`
}

var users = sync.Map{}

func InitMessageServer(pool *pgxpool.Pool) *MessageServer {
	godotenv.Load()
	frontURL := os.Getenv("FRONTEND_URL")

	if frontURL == "" {
        logger.Println("Warning: FRONTEND_URL not set in .env")
    }

	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			logger.Println("Origin:", origin)
			if origin == frontURL || origin == "test_client" {
				return true
			}
			return false
		},
	}

	return &MessageServer{
		pool: pool,
		upgrader: upgrader,
	}
}

func (s *MessageServer) StartWebSocketServer(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Printf("WebSocket upgrade failed: %v", err)
        return
	}

	defer func() {
		if s.senderID != "" {
			users.Delete(s.senderID)
			logger.Println(fmt.Sprintf("User %s disconnected", s.senderID))
		}
		conn.Close()
	}()

	var authMessage AuthMessage
	if err := conn.ReadJSON(&authMessage); err != nil {
		logger.Println("Failed to read auth message:", err)
		return
	}

	if authMessage.Type != "auth" || authMessage.UserID == "" {
		logger.Println("Invalid auth message:", authMessage)
		conn.WriteJSON(map[string]string{"error": "Invalid auth message"})
		return
	}

	users.Store(authMessage.UserID, conn) 
	s.senderID = authMessage.UserID

	s.read(conn)
}

func (s *MessageServer) read(conn *websocket.Conn) {
	for {
		var message Message
		if err := conn.ReadJSON(&message); err != nil {
			logger.Println(err)
			return
		}

		sendMessage(message.ReceiverID, message.Text)
		s.SaveMessage(&message)
	}
}

func (s *MessageServer) SaveMessage(message *Message) {
	query := "INSERT INTO messages (message, sender_id, receiver_id) VALUES ($1, $2, $3)"
	_, err := s.pool.Exec(context.Background(), query, message.Text, s.senderID, message.ReceiverID)
	if err != nil {
		logger.Println(err)
	}
}

func sendMessage(receiverID, message string) {
	conn, ok := users.Load(receiverID)
	if !ok {
		logger.Println(fmt.Sprintf("User %s not found", receiverID))
		return
	}

	wsConn := conn.(*websocket.Conn)
	wsConn.WriteJSON(map[string]string{"message": message})
}
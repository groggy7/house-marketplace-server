package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var logger = log.New(os.Stdout, "chat server - ", log.Lshortfile)

type MessageServer struct {
	pool     *pgxpool.Pool
	upgrader websocket.Upgrader
}

type AuthMessage struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
}

type Message struct {
	Text       string `json:"text"`
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
		pool:     pool,
		upgrader: upgrader,
	}
}

func (s *MessageServer) StartWebSocketServer(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	var authMessage AuthMessage
	if err := conn.ReadJSON(&authMessage); err != nil {
		logger.Println("Failed to read auth message:", err)
		conn.Close()
		return
	}

	if authMessage.Type != "auth" || authMessage.UserID == "" {
		logger.Println("Invalid auth message:", authMessage)
		conn.WriteJSON(map[string]string{"error": "Invalid auth message"})
		conn.Close()
		return
	}

	userID := authMessage.UserID

	users.Store(userID, conn)
	logger.Printf("User %s connected\n", userID)

	defer func() {
		users.Delete(userID)
		conn.Close()
		logger.Printf("User %s disconnected\n", userID)
	}()

	s.handleMessages(conn, userID)
}

func (s *MessageServer) handleMessages(conn *websocket.Conn, senderID string) {
    for {
        var message Message
        if err := conn.ReadJSON(&message); err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                logger.Println("WebSocket closed:", err)
                return
            }
            logger.Println("Temporary read error:", err)
            
            if _, ok := err.(*websocket.CloseError); ok {
                logger.Println("Connection closed by client")
                return
            }
            
            if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") || 
               strings.Contains(err.Error(), "unexpected EOF") {
                logger.Println("Connection appears to be closed:", err)
                return
            }
            
            continue
        }

        s.sendMessage(senderID, message.ReceiverID, message.Text)
        s.SaveMessage(senderID, &message)
    }
}

func (s *MessageServer) SaveMessage(senderID string, message *Message) {
	query := "INSERT INTO messages (message, sender_id, receiver_id) VALUES ($1, $2, $3)"
	_, err := s.pool.Exec(context.Background(), query, message.Text, senderID, message.ReceiverID)
	if err != nil {
		logger.Println(err)
	}
}

func (s *MessageServer) sendMessage(senderID, receiverID, message string) {
	conn, ok := users.Load(receiverID)
	if !ok {
		logger.Printf("User %s not found\n", receiverID)
		return
	}

	wsConn := conn.(*websocket.Conn)
	if err := wsConn.WriteJSON(map[string]string{
		"type":      "message",
		"text":      message,
		"sender_id": senderID,
	}); err != nil {
		logger.Println("Error sending to recipient:", err)
		users.Delete(receiverID)
		return
	}

	client, ok := users.Load(senderID)
	if !ok {
		logger.Printf("User %s not found\n", senderID)
		return
	}

	clientConn := client.(*websocket.Conn)
	if err := clientConn.WriteJSON(map[string]string{"status": "sent", "message": message}); err != nil {
		logger.Println("Error sending confirmation to sender:", err)
		return
	}
}

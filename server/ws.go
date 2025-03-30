package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

var logger = log.New(os.Stdout, "chat server - ", log.Lshortfile)

type MessageServer struct {
	db *pgx.Conn
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

var users = make(map[string]*websocket.Conn)

func InitMessageServer(db *pgx.Conn) *MessageServer {
	return &MessageServer{
		db: db,
		upgrader: upgrader,
	}
}

func (s *MessageServer) StartWebSocketServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Fatalln(err)
	}

	defer func() {
		if s.senderID != "" {
			delete(users, s.senderID)
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

	users[authMessage.UserID] = conn
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

		log.Println(fmt.Sprintf("Received message: %s, from: %s, to: %s", message.Text, s.senderID, message.ReceiverID))
		sendMessage(message.ReceiverID, message.Text)
		s.SaveMessage(s.db, message.Text, message.ReceiverID)
	}
}

func (s *MessageServer) SaveMessage(db *pgx.Conn, message, receiverID string) {
	query := "INSERT INTO messages (message, sender_id, receiver_id) VALUES ($1, $2, $3)"
	_, err := db.Exec(context.Background(), query, message, s.senderID, receiverID)
	if err != nil {
		logger.Println(err)
	}
}

func sendMessage(receiverID, message string) {
	conn, ok := users[receiverID]
	
	if !ok {
		logger.Println(fmt.Sprintf("User %s not found", receiverID))
		return
	}

	if err := conn.WriteJSON(map[string]string{"message": message}); err != nil {
		logger.Println(err)
	}
}
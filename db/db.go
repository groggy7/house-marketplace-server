package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Room struct {
	PropertyID      string `json:"property_id"`
	PropertyOwnerID string `json:"property_owner_id"`
	CustomerID      string `json:"customer_id"`
}

type ChatDB struct {
	conn *pgxpool.Pool
}

func NewChatDB() (*ChatDB, error) {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL not set in .env")
	}

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return &ChatDB{conn: conn}, nil
}

func (db *ChatDB) CreateRoom(propertyID, propertyOwnerID, customerID string) (string, error) {
	query := "INSERT INTO rooms (property_id, property_owner_id, customer_id) VALUES ($1, $2, $3) RETURNING id"
	var roomID string
	err := db.conn.QueryRow(context.Background(), query, propertyID, propertyOwnerID, customerID).Scan(&roomID)
	if err != nil {
		return "", err
	}

	return roomID, nil
}

func (db *ChatDB) CheckRoomExists(roomID string) (bool, error) {
	query := "SELECT id FROM rooms WHERE id = $1"
	var id string
	err := db.conn.QueryRow(context.Background(), query, roomID).Scan(&id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *ChatDB) GetRooms(customerID string) ([]Room, error) {
	query := "SELECT (property_id, property_owner_id, customer_id) FROM rooms WHERE customer_id = $1"
	rows, err := db.conn.Query(context.Background(), query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.PropertyID, &room.PropertyOwnerID, &room.CustomerID); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (db *ChatDB) SaveMessage(text, senderID, roomID string) error {
	query := "INSERT INTO messages (message, sender_id, room_id) VALUES ($1, $2, $3)"
	_, err := db.conn.Exec(context.Background(), query, text, senderID, roomID)
	if err != nil {
		return err
	}

	return nil
}

func (db *ChatDB) CheckUserInRoom(userID, roomID string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM rooms WHERE id = $1 AND (property_owner_id = $2 OR customer_id = $2))"
	var exists bool
	err := db.conn.QueryRow(context.Background(), query, roomID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user room membership: %w", err)
	}

	return exists, nil
}

package repository

import (
	"context"
	"fmt"
	"message-server/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type roomRepository struct {
	pool *pgxpool.Pool
}

func NewRoomRepository(pool *pgxpool.Pool) domain.RoomRepository {
	return &roomRepository{pool: pool}
}

func (db *roomRepository) CreateRoom(
	propertyID,
	propertyOwnerID,
	customerID,
	listingTitle,
	listingImage,
	ownerName,
	customerName string,
) (string, error) {
	query := `
		INSERT INTO rooms (property_id, property_owner_id, customer_id, listing_title, listing_image) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`
	var roomID string
	err := db.pool.QueryRow(context.Background(), query, propertyID,
		propertyOwnerID, customerID, listingTitle, listingImage).Scan(&roomID)
	if err != nil {
		return "", err
	}

	return roomID, nil
}

func (db *roomRepository) CheckRoomExists(roomID string) (bool, error) {
	query := "SELECT id FROM rooms WHERE id = $1"
	var id string
	err := db.pool.QueryRow(context.Background(), query, roomID).Scan(&id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *roomRepository) GetRooms(customerID string) ([]domain.Room, error) {
	query := `
		SELECT id, property_id, property_owner_id, customer_id, listing_title, listing_image
		FROM rooms 
		WHERE customer_id = $1 OR property_owner_id = $1
	`
	rows, err := db.pool.Query(context.Background(), query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var room domain.Room
		if err := rows.Scan(&room.RoomID, &room.PropertyID, &room.PropertyOwnerID,
			&room.CustomerID, &room.Title, &room.Image); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (db *roomRepository) SaveMessage(text, senderID, senderName, roomID string) error {
	query := "INSERT INTO messages (message, sender_id, sender_name, room_id) VALUES ($1, $2, $3, $4)"
	_, err := db.pool.Exec(context.Background(), query, text, senderID, senderName, roomID)
	if err != nil {
		return err
	}

	return nil
}

func (db *roomRepository) CheckUserInRoom(userID, roomID string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM rooms WHERE id = $1 AND (property_owner_id = $2 OR customer_id = $2))"
	var exists bool
	err := db.pool.QueryRow(context.Background(), query, roomID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user room membership: %w", err)
	}

	return exists, nil
}

func (db *roomRepository) GetMessagesForRoom(roomID string) ([]map[string]any, error) {
	query := `
		SELECT id, message, sender_id, sender_name, room_id, created_at 
		FROM messages 
		WHERE room_id = $1 
		ORDER BY created_at ASC
	`

	rows, err := db.pool.Query(context.Background(), query, roomID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving messages: %w", err)
	}
	defer rows.Close()

	messages := []map[string]any{}
	for rows.Next() {
		var id, message, senderID, senderName, roomID string
		var createdAt any

		if err := rows.Scan(&id, &message, &senderID, &senderName, &roomID, &createdAt); err != nil {
			return nil, fmt.Errorf("error scanning message row: %w", err)
		}

		messages = append(messages, map[string]any{
			"id":          id,
			"message":     message,
			"sender_id":   senderID,
			"sender_name": senderName,
			"room_id":     roomID,
			"created_at":  createdAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	return messages, nil
}

package server

import (
	"encoding/json"
	"message-server/db"
	"net/http"
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

func (s *RoomServer) CreateRoom(rw http.ResponseWriter, r *http.Request) {
	var request CreateChatRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	if request.PropertyID == "" || request.PropertyOwnerID == "" || request.CustomerID == "" {
		http.Error(rw, "Missing required fields", http.StatusBadRequest)
		return
	}

	roomID, err := s.db.CreateRoom(request.PropertyID, request.PropertyOwnerID, request.CustomerID)
	if err != nil {
		http.Error(rw, "Failed to create chat room", http.StatusInternalServerError)
		logger.Printf("Failed to create chat room: %v", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"room_id": roomID})
}

func (s *RoomServer) GetRooms(rw http.ResponseWriter, r *http.Request) {
	var request GetRoomsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	if request.CustomerID == "" {
		http.Error(rw, "Missing required fields", http.StatusBadRequest)
		return
	}

	rooms, err := s.db.GetRooms(request.CustomerID)
	if err != nil {
		http.Error(rw, "Failed to get chat rooms", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(rooms)
}

package domain

type Room struct {
	RoomID       string `json:"room_id"`
	PropertyID   string `json:"property_id"`
	OwnerID      string `json:"owner_id"`
	CustomerID   string `json:"customer_id"`
	Title        string `json:"title"`
	Image        string `json:"image"`
	OwnerName    string `json:"owner_name"`
	CustomerName string `json:"customer_name"`
}

type AuthMessage struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
}

type ChatMessage struct {
	Text       string `json:"text"`
	ReceiverID string `json:"receiver_id"`
	SenderID   string `json:"sender_id"`
	RoomID     string `json:"room_id"`
}

type MessageResponse struct {
	Type      string `json:"type"`
	Text      string `json:"text,omitempty"`
	SenderID  string `json:"sender_id,omitempty"`
	RoomID    string `json:"room_id,omitempty"`
	Status    string `json:"status,omitempty"`
	Error     string `json:"error,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type CreateChatRoomRequest struct {
	PropertyID string `json:"property_id" validate:"required"`
	OwnerID    string `json:"owner_id" validate:"required"`
	CustomerID string `json:"-"`
}

type RoomRepository interface {
	CreateRoom(propertyID, ownerID, ownerName, customerID, customerName, title, image string) (string, error)
	CheckRoomExists(roomID string) (bool, error)
	GetRooms(customerID string) ([]Room, error)
	SaveMessage(text, senderID, senderName, roomID string) error
	CheckUserInRoom(userID, roomID string) (bool, error)
	GetMessagesForRoom(roomID string) ([]map[string]any, error)
}

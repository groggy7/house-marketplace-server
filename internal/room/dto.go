package room

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
	PropertyID      string `json:"property_id"`
	PropertyOwnerID string `json:"property_owner_id"`
	CustomerID      string `json:"customer_id"`
}

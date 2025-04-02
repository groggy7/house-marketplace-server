package room

type RoomRepository interface {
	CreateRoom(propertyID, propertyOwnerID, customerID string) (string, error)
	CheckRoomExists(roomID string) (bool, error)
	GetRooms(customerID string) ([]Room, error)
	SaveMessage(text, senderID, roomID string) error
	CheckUserInRoom(userID, roomID string) (bool, error)
	GetMessagesForRoom(roomID string) ([]map[string]any, error)
}

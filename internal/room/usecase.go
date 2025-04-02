package room

type RoomService struct {
	roomRepo RoomRepository
}

func NewRoomService(roomRepo RoomRepository) RoomService {
	return RoomService{roomRepo: roomRepo}
}

func (s *RoomService) CreateRoom(req *CreateChatRoomRequest) (string, error) {
	return s.roomRepo.CreateRoom(req.PropertyID, req.PropertyOwnerID, req.CustomerID)
}

func (s *RoomService) CheckRoomExists(roomID string) (bool, error) {
	return s.roomRepo.CheckRoomExists(roomID)
}

func (s *RoomService) GetRooms(customerID string) ([]Room, error) {
	rooms, err := s.roomRepo.GetRooms(customerID)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *RoomService) SaveMessage(text, senderID, roomID string) error {
	return s.roomRepo.SaveMessage(text, senderID, roomID)
}

func (s *RoomService) CheckUserInRoom(userID, roomID string) (bool, error) {
	return s.roomRepo.CheckUserInRoom(userID, roomID)
}

func (s *RoomService) GetMessagesForRoom(roomID string) ([]map[string]any, error) {
	return s.roomRepo.GetMessagesForRoom(roomID)
}

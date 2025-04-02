package usecases

import "message-server/internal/domain"

type RoomUseCase struct {
	roomRepo domain.RoomRepository
}

func NewRoomUseCase(roomRepo domain.RoomRepository) *RoomUseCase {
	return &RoomUseCase{roomRepo: roomRepo}
}

func (s *RoomUseCase) CreateRoom(req *domain.CreateChatRoomRequest) (string, error) {
	return s.roomRepo.CreateRoom(req.PropertyID, req.PropertyOwnerID, req.CustomerID)
}

func (s *RoomUseCase) CheckRoomExists(roomID string) (bool, error) {
	return s.roomRepo.CheckRoomExists(roomID)
}

func (s *RoomUseCase) GetRooms(customerID string) ([]domain.Room, error) {
	rooms, err := s.roomRepo.GetRooms(customerID)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *RoomUseCase) SaveMessage(text, senderID, roomID string) error {
	return s.roomRepo.SaveMessage(text, senderID, roomID)
}

func (s *RoomUseCase) CheckUserInRoom(userID, roomID string) (bool, error) {
	return s.roomRepo.CheckUserInRoom(userID, roomID)
}

func (s *RoomUseCase) GetMessagesForRoom(roomID string) ([]map[string]any, error) {
	return s.roomRepo.GetMessagesForRoom(roomID)
}

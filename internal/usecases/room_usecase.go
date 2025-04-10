package usecases

import (
	"fmt"
	"message-server/internal/domain"
)

type RoomUseCase struct {
	roomRepo domain.RoomRepository
	authRepo domain.AuthRepository
}

func NewRoomUseCase(roomRepo domain.RoomRepository, authRepo domain.AuthRepository) *RoomUseCase {
	return &RoomUseCase{
		roomRepo: roomRepo,
		authRepo: authRepo,
	}
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
	user, err := s.authRepo.GetUserByID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get sender info: %w", err)
	}

	return s.roomRepo.SaveMessage(text, senderID, user.FullName, roomID)
}

func (s *RoomUseCase) CheckUserInRoom(userID, roomID string) (bool, error) {
	return s.roomRepo.CheckUserInRoom(userID, roomID)
}

func (s *RoomUseCase) GetMessagesForRoom(roomID string) ([]map[string]any, error) {
	return s.roomRepo.GetMessagesForRoom(roomID)
}

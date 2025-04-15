package usecases

import (
	"fmt"
	"message-server/internal/domain"
)

type RoomUseCase struct {
	roomRepo    domain.RoomRepository
	authRepo    domain.AuthRepository
	listingRepo domain.ListingRepository
}

func NewRoomUseCase(
	roomRepo domain.RoomRepository,
	authRepo domain.AuthRepository,
	listingRepo domain.ListingRepository,
) *RoomUseCase {
	return &RoomUseCase{
		roomRepo:    roomRepo,
		authRepo:    authRepo,
		listingRepo: listingRepo,
	}
}

func (s *RoomUseCase) CreateRoom(req *domain.CreateChatRoomRequest) (string, error) {
	listing, err := s.listingRepo.GetListingByID(req.PropertyID)
	if err != nil {
		return "", err
	}

	title := listing.Title
	image := listing.ImageURLs[0]

	owner, err := s.authRepo.GetUserByID(req.OwnerID)
	if err != nil {
		return "", err
	}

	customer, err := s.authRepo.GetUserByID(req.CustomerID)
	if err != nil {
		return "", err
	}

	return s.roomRepo.CreateRoom(req.PropertyID, req.OwnerID, owner.FullName, req.CustomerID, customer.FullName, title, image)
}

func (s *RoomUseCase) CheckRoomExists(roomID string) (bool, error) {
	return s.roomRepo.CheckRoomExists(roomID)
}

func (s *RoomUseCase) GetRooms(customerID string) ([]domain.Room, error) {
	rooms, err := s.roomRepo.GetRooms(customerID)
	if err != nil {
		return nil, err
	}

	if rooms == nil {
		return []domain.Room{}, nil
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

package usecases

import (
	"message-server/internal/domain"
)

type ListingUseCase struct {
	listingRepo domain.ListingRepository
	fileRepo    domain.FileRepository
}

func NewListingUseCase(listingRepo domain.ListingRepository, fileRepo domain.FileRepository) *ListingUseCase {
	return &ListingUseCase{listingRepo: listingRepo, fileRepo: fileRepo}
}

func (s *ListingUseCase) CreateListing(request *domain.CreateListingRequest) (string, error) {
	return s.listingRepo.CreateListing(request)
}

func (s *ListingUseCase) GetListingByID(id string) (*domain.GetListingDetailsResponse, error) {
	return s.listingRepo.GetListingByID(id)
}

func (s *ListingUseCase) GetListings() (*domain.GetListingsResponse, error) {
	return s.listingRepo.GetListings()
}

func (s *ListingUseCase) UpdateListing(listing *domain.Listing) error {
	return s.listingRepo.UpdateListing(listing)
}

func (s *ListingUseCase) DeleteListing(id string) error {
	listing, err := s.listingRepo.GetListingByID(id)
	if err != nil {
		return err
	}

	for _, key := range listing.ImageKeys {
		err := s.fileRepo.DeleteFile(key)
		if err != nil {
			return err
		}
	}

	return s.listingRepo.DeleteListing(id)
}

func (s *ListingUseCase) BookmarkListing(userID, listingID string) error {
	return s.listingRepo.BookmarkListing(userID, listingID)
}

func (s *ListingUseCase) UnbookmarkListing(userID, listingID string) error {
	return s.listingRepo.UnbookmarkListing(userID, listingID)
}

func (s *ListingUseCase) GetBookmarkedListings(userID string) ([]domain.ListingInfo, error) {
	return s.listingRepo.GetBookmarkedListings(userID)
}

package controller

import (
	"message-server/internal/controller/auth"
	"message-server/internal/domain"
	"message-server/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListingHandler struct {
	listingUseCase usecases.ListingUseCase
}

func NewListingHandler(listingUseCase *usecases.ListingUseCase) *ListingHandler {
	return &ListingHandler{listingUseCase: *listingUseCase}
}

func (s *ListingHandler) CreateListing(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := claims.(*auth.Claims).UserID

	var request domain.CreateListingRequest
	request.UserID = userID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	listingID, err := s.listingUseCase.CreateListing(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": listingID})
}

func (s *ListingHandler) GetListingByID(c *gin.Context) {
	id := c.Param("id")

	listing, err := s.listingUseCase.GetListingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, listing)
}

func (s *ListingHandler) GetListings(c *gin.Context) {
	listings, err := s.listingUseCase.GetListings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, listings)
}

func (s *ListingHandler) UpdateListing(c *gin.Context) {
	var request domain.Listing
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.listingUseCase.UpdateListing(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing updated successfully"})
}

func (s *ListingHandler) DeleteListing(c *gin.Context) {
	id := c.Param("id")

	if err := s.listingUseCase.DeleteListing(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}

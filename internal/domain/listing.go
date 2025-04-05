package domain

import "time"

type Listing struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Type               string    `json:"type"`
	Price              int       `json:"price"`
	Location           string    `json:"location"`
	Bathrooms          int       `json:"bathrooms"`
	Bedrooms           int       `json:"bedrooms"`
	CreatedAt          time.Time `json:"created_at"`
	ImageURLs          []string  `json:"image_urls"`
	IsAirConditioned   bool      `json:"is_air_conditioned"`
	IsBalconyAvailable bool      `json:"is_balcony_available"`
	IsDryerAvailable   bool      `json:"is_dryer_available"`
	IsHeated           bool      `json:"is_heated"`
	IsParkingAvailable bool      `json:"is_parking_available"`
	IsPoolAvailable    bool      `json:"is_pool_available"`
	IsWasherAvailable  bool      `json:"is_washer_available"`
	IsWifiAvailable    bool      `json:"is_wifi_available"`
}

type CreateListingRequest struct {
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Type               string   `json:"type"`
	Price              int      `json:"price"`
	Location           string   `json:"location"`
	Bathrooms          int      `json:"bathrooms"`
	Bedrooms           int      `json:"bedrooms"`
	ImageURLs          []string `json:"image_urls"`
	IsAirConditioned   bool     `json:"is_air_conditioned"`
	IsBalconyAvailable bool     `json:"is_balcony_available"`
	IsDryerAvailable   bool     `json:"is_dryer_available"`
	IsHeated           bool     `json:"is_heated"`
	IsParkingAvailable bool     `json:"is_parking_available"`
	IsPoolAvailable    bool     `json:"is_pool_available"`
	IsWasherAvailable  bool     `json:"is_washer_available"`
	IsWifiAvailable    bool     `json:"is_wifi_available"`
	UserID             string   `json:"user_id"`
}

type GetListingsResponse struct {
	Listings []ListingInfo `json:"listings"`
}

type ListingInfo struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Type      string   `json:"type"`
	Price     int      `json:"price"`
	Location  string   `json:"location"`
	Bathrooms int      `json:"bathrooms"`
	Bedrooms  int      `json:"bedrooms"`
	ImageURLs []string `json:"image_urls"`
}

type DeleteListingRequest struct {
	ID string `json:"id"`
}

type ListingRepository interface {
	CreateListing(request *CreateListingRequest) (string, error)
	GetListingByID(id string) (*Listing, error)
	GetListings() (*GetListingsResponse, error)
	UpdateListing(listing *Listing) error
	DeleteListing(id string) error
	BookmarkListing(userID, listingID string) error
	UnbookmarkListing(userID, listingID string) error
}

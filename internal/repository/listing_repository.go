package repository

import (
	"context"
	"message-server/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type listingRepository struct {
	pool *pgxpool.Pool
}

func NewListingRepository(pool *pgxpool.Pool) domain.ListingRepository {
	return &listingRepository{pool: pool}
}

func (r *listingRepository) CreateListing(req *domain.CreateListingRequest) (string, error) {
	query := `
		INSERT INTO listings 
		(id, title, description, type, price, location, bathrooms, 
		bedrooms, image_keys, is_air_conditioned, is_balcony_available,
		is_dryer_available,  is_heated, is_parking_available, 
		is_pool_available, is_washer_available, is_wifi_available, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id
	`

	var id string
	err := r.pool.QueryRow(context.Background(), query, req.ID, req.Title, req.Description, req.Type,
		req.Price, req.Location, req.Bathrooms, req.Bedrooms, req.ImageKeys,
		req.IsAirConditioned, req.IsBalconyAvailable, req.IsDryerAvailable, req.IsHeated,
		req.IsParkingAvailable, req.IsPoolAvailable, req.IsWasherAvailable, req.IsWifiAvailable, req.UserID).Scan(&id)
	return id, err
}

func (r *listingRepository) GetListingByID(id string) (*domain.GetListingDetailsResponse, error) {
	query := `
		SELECT id, title, description, type, price, location, bathrooms, 
		bedrooms, image_keys, is_air_conditioned, is_balcony_available, is_dryer_available,
		is_heated, is_parking_available, is_pool_available, is_washer_available, is_wifi_available, user_id
		FROM listings
		WHERE id = $1
	`

	var listing domain.GetListingDetailsResponse

	err := r.pool.QueryRow(context.Background(), query, id).Scan(&listing.ID, &listing.Title,
		&listing.Description, &listing.Type, &listing.Price, &listing.Location, &listing.Bathrooms,
		&listing.Bedrooms, &listing.ImageKeys, &listing.IsAirConditioned, &listing.IsBalconyAvailable,
		&listing.IsDryerAvailable, &listing.IsHeated, &listing.IsParkingAvailable,
		&listing.IsPoolAvailable, &listing.IsWasherAvailable, &listing.IsWifiAvailable, &listing.UserID)

	return &listing, err
}

func (r *listingRepository) GetListings() (*domain.GetListingsResponse, error) {
	query := `
		SELECT id, title, type, price, location, bathrooms, bedrooms, image_keys
		FROM listings`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listings := []domain.ListingInfo{}
	for rows.Next() {
		var listing domain.ListingInfo
		err := rows.Scan(&listing.ID, &listing.Title, &listing.Type, &listing.Price, &listing.Location, &listing.Bathrooms, &listing.Bedrooms, &listing.ImageKeys)
		if err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}

	return &domain.GetListingsResponse{Listings: listings}, nil
}

func (r *listingRepository) UpdateListing(listing *domain.Listing) error {
	query := `
		UPDATE listings
		SET title = $1, description = $2, type = $3, price = $4, location = $5, bathrooms = $6, bedrooms = $7, image_keys = $8, is_air_conditioned = $9, is_balcony_available = $10, is_dryer_available = $11, is_heated = $12, is_parking_available = $13, is_pool_available = $14, is_washer_available = $15, is_wifi_available = $16
		WHERE id = $1
	`

	_, err := r.pool.Exec(context.Background(), query, listing.Title, listing.Description, listing.Type, listing.Price, listing.Location, listing.Bathrooms, listing.Bedrooms, listing.ImageKeys, listing.IsAirConditioned, listing.IsBalconyAvailable, listing.IsDryerAvailable, listing.IsHeated, listing.IsParkingAvailable, listing.IsPoolAvailable, listing.IsWasherAvailable, listing.IsWifiAvailable)
	return err
}

func (r *listingRepository) DeleteListing(id string) error {
	query := `
		DELETE FROM listings
		WHERE id = $1
	`

	_, err := r.pool.Exec(context.Background(), query, id)
	return err
}

func (r *listingRepository) BookmarkListing(userID, listingID string) error {
	query := `
		INSERT INTO bookmarks (user_id, listing_id)
		VALUES ($1, $2)
	`

	_, err := r.pool.Exec(context.Background(), query, userID, listingID)
	return err
}

func (r *listingRepository) UnbookmarkListing(userID, listingID string) error {
	query := `
		DELETE FROM bookmarks
		WHERE user_id = $1 AND listing_id = $2
	`

	_, err := r.pool.Exec(context.Background(), query, userID, listingID)
	return err
}

func (r *listingRepository) GetBookmarkedListings(userID string) ([]domain.ListingInfo, error) {
	query := `
		SELECT l.id, l.title, l.type, l.price, l.location, l.bathrooms, l.bedrooms, l.image_keys
		FROM listings l
		JOIN bookmarks b ON l.id = b.listing_id
		WHERE b.user_id = $1
	`

	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listings := []domain.ListingInfo{}
	for rows.Next() {
		var listing domain.ListingInfo
		err := rows.Scan(&listing.ID, &listing.Title, &listing.Type, &listing.Price, &listing.Location, &listing.Bathrooms, &listing.Bedrooms, &listing.ImageKeys)
		if err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}

	return listings, nil
}

package room

type Room struct {
	RoomID          string `json:"room_id"`
	PropertyID      string `json:"property_id"`
	PropertyOwnerID string `json:"property_owner_id"`
	CustomerID      string `json:"customer_id"`
}

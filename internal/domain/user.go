package domain

type UpdateUserRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	AvatarURL string `json:"-"`
	UserID    string `json:"-"`
	Email     string `json:"-"`
}

type UserRepository interface {
	UpdateUser(name, avatarURL string, userID string) error
}

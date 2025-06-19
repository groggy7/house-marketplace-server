package domain

type GenerateListingUploadURLRequest struct {
	ListingID   string `json:"listing_id"`
	ContentType string `json:"content_type"`
}

type GenerateAvatarUploadURLRequest struct {
	ID          string `json:"id"`
	ContentType string `json:"content_type"`
}

type URLResponse struct {
	URL string `json:"url"`
}

type DeleteFileRequest struct {
	Key string `json:"key"`
}

type GenerateDownloadURLRequest struct {
	Key string
}

type FileRepository interface {
	GenerateListingUploadURL(string, string, string) (string, error)
	GenerateAvatarUploadURL(string, string) (string, error)
	GenerateDownloadURL(string) (string, error)
	DeleteFile(string) error
}

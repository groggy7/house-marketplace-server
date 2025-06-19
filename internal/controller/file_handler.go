package controller

import (
	"message-server/internal/domain"
	"message-server/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileUseCase usecases.FileUseCase
}

func NewFileHandler(fileUseCase *usecases.FileUseCase) *FileHandler {
	return &FileHandler{fileUseCase: *fileUseCase}
}

func (fh *FileHandler) GenerateListingUploadURL(c *gin.Context) {
	var req *domain.GenerateListingUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	URL, err := fh.fileUseCase.GenerateListingUploadURL(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate upload URL"})
		return
	}

	c.JSON(http.StatusOK, URL)
}

func (fh *FileHandler) GenerateAvatarUploadURL(c *gin.Context) {
	var req *domain.GenerateAvatarUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	URL, err := fh.fileUseCase.GenerateAvatarUploadURL(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate upload URL"})
		return
	}

	c.JSON(http.StatusOK, URL)
}

func (s *FileHandler) GenerateDownloadURL(c *gin.Context) {
	var req domain.GenerateDownloadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := s.fileUseCase.GenerateDownloadURL(req.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *FileHandler) DeleteFile(c *gin.Context) {
	var request domain.DeleteFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := s.fileUseCase.DeleteFile(request.Key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

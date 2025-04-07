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

func (s *FileHandler) UploadListingPicture(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	response, err := s.fileUseCase.UploadListingPicture(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *FileHandler) UploadProfilePicture(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	response, err := s.fileUseCase.UploadListingPicture(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
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

	if err := s.fileUseCase.DeleteFile(request.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

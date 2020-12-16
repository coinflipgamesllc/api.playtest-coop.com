package controller

import (
	"strconv"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/gin-gonic/gin"
)

// FileController handles /files routes
type FileController struct {
	FileService *app.FileService
}

// PresignUploadRequest params for presigning a URL
type PresignUploadRequest struct {
	Name      string `form:"name" binding:"required" example:"my-awesome-file.jpg"`
	Extension string `form:"extension" binding:"required" example:"jpg"`
}

// PresignUploadResponse wrapper for presigned URL
type PresignUploadResponse struct {
	URL string `json:"url" example:"https://assets.playtest-coop.com/..."`
}

// PresignUpload generates a presigned URL for the client to upload directly to S3
// @Summary Generate a presigned URL for the client to upload directly to S3
// @Accept json
// @Produce json
// @Param file body PresignUploadRequest true "File data"
// @Success 200 {object} PresignUploadResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags files
// @Router /files/sign [get]
func (t *FileController) PresignUpload(c *gin.Context) {
	// Validate request
	var req PresignUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	presignedURL, err := t.FileService.PresignUpload(req.Name, req.Extension)
	if err != nil {
		serverErrorResponse(c, "presigned url could not be generated")
		return
	}

	c.JSON(200, PresignUploadResponse{URL: presignedURL})
}

// CreateFileRequest params for storing a record of a file
type CreateFileRequest struct {
	Role     string `json:"role" binding:"required" example:"Image"`
	Caption  string `json:"caption" example:"What a cool image of a game!"`
	Filename string `json:"filename" binding:"required" example:"example-image.png"`
	Object   string `json:"object" binding:"required" example:"asd9fhgaoseucgewio.png"`
	Size     int64  `json:"size" binding:"required" example:"1241231"`
	GameID   uint   `json:"game" example:"123"`
}

// CreateFile saves a record of a file stored in S3
// @Summary Save a record of a file stored in S3
// @Accept json
// @Produce json
// @Param file body CreateFileRequest true "File data"
// @Success 200 {object} AckResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags files
// @Router /files [post]
func (t *FileController) CreateFile(c *gin.Context) {
	// Validate the request
	var req CreateFileRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	userID := userID(c)
	_, err := t.FileService.CreateFile(userID, req.Role, req.Filename, req.Object, req.Size, req.Caption, req.GameID)

	if err != nil {
		requestErrorResponse(c, "failed to save file")
		return
	}

	ackResponse(c)
}

// ListUserFilesResponse wrapper for files belonging to a user
type ListUserFilesResponse struct {
	Files []domain.File `json:"files"`
}

// ListUserFiles lists files belonging to the authenticated user
// @Summary List files belonging to the authenticated user
// @Produce json
// @Success 200 {object} ListUserFilesResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags files
// @Router /files [get]
func (t *FileController) ListUserFiles(c *gin.Context) {
	userID := userID(c)

	files, err := t.FileService.ListUserFiles(userID)
	if err != nil {
		serverErrorResponse(c, "failed to fetch files")
		return
	}

	c.JSON(200, ListUserFilesResponse{Files: files})
}

// DeleteFile removes a file by ID
// @Summary remove a file by ID
// @Produce json
// @Param id path integer true "File ID"
// @Success 200 {object} AckResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags files
// @Router /files/:id [delete]
func (t *FileController) DeleteFile(c *gin.Context) {
	// Pull file by ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, "invalid file id")
		return
	}

	userID := userID(c)

	if err := t.FileService.DeleteFile(uint(id), userID); err != nil {
		serverErrorResponse(c, "failed to delete file")
		return
	}

	ackResponse(c)
}

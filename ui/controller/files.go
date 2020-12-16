package controller

import (
	"strconv"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/gin-gonic/gin"
)

// FileController handles /files routes
type FileController struct {
	FileService *app.FileService
}

// PresignUpload generates a presigned URL for the client to upload directly to S3
// @Summary Generate a presigned URL for the client to upload directly to S3
// @Accept json
// @Produce json
// @Param file body app.PresignUploadRequest true "File data"
// @Success 200 {object} app.PresignUploadResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags files
// @Router /files/sign [get]
func (t *FileController) PresignUpload(c *gin.Context) {
	// Validate request
	var req app.PresignUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	presignedURL, err := t.FileService.PresignUpload(req.Name, req.Extension)
	if err != nil {
		serverErrorResponse(c, "presigned url could not be generated")
		return
	}

	c.JSON(200, app.PresignUploadResponse{URL: presignedURL})
}

// CreateFile saves a record of a file stored in S3
// @Summary Save a record of a file stored in S3
// @Accept json
// @Produce json
// @Param file body app.CreateFileRequest true "File data"
// @Success 200 {object} AckResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags files
// @Router /files [post]
func (t *FileController) CreateFile(c *gin.Context) {
	// Validate the request
	var req app.CreateFileRequest
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

// ListUserFiles lists files belonging to the authenticated user
// @Summary List files belonging to the authenticated user
// @Produce json
// @Success 200 {object} app.ListFilesResponse
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

	c.JSON(200, app.ListFilesResponse{Files: files})
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

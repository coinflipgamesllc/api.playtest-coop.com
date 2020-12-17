package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type (
	// FileService handles file uploads/downloads and indirect interaction with S3
	FileService struct {
		FileRepository domain.FileRepository
		GameRepository domain.GameRepository
		UserRepository domain.UserRepository
		Logger         *zap.Logger
		S3Bucket       string
		S3Client       *minio.Client
	}

	// Request DTOs

	// PresignUploadRequest params for presigning a URL
	PresignUploadRequest struct {
		Name      string `form:"name" binding:"required" example:"my-awesome-file.jpg"`
		Extension string `form:"extension" binding:"required" example:"jpg"`
	}

	// CreateFileRequest params for storing a record of a file
	CreateFileRequest struct {
		Role     string `json:"role" binding:"required" example:"Image"`
		Caption  string `json:"caption" example:"What a cool image of a game!"`
		Filename string `json:"filename" binding:"required" example:"example-image.png"`
		Object   string `json:"object" binding:"required" example:"asd9fhgaoseucgewio.png"`
		Size     int64  `json:"size" binding:"required" example:"1241231"`
		GameID   uint   `json:"game" example:"123"`
	}

	// Response DTOs

	// PresignUploadResponse wrapper for presigned URL
	PresignUploadResponse struct {
		URL string `json:"url" example:"https://assets.playtest-coop.com/..."`
	}

	// ListFilesResponse wrapper for files belonging to a user
	ListFilesResponse struct {
		Files []domain.File `json:"files"`
	}
)

// PresignUpload generates a presigned URL for uploading a file
func (s *FileService) PresignUpload(name, extension string) (string, error) {
	presignedURL, err := s.S3Client.PresignedPutObject(
		context.Background(),
		s.S3Bucket,
		domain.GenerateObjectName(name, extension),
		time.Duration(1000)*time.Minute,
	)

	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// CreateFile stores a file in the database, optionally tied to a game
func (s *FileService) CreateFile(req *CreateFileRequest, userID uint) (*domain.File, error) {
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	// Create the file
	var file *domain.File
	switch req.Role {
	case "Image":
		file, err = domain.NewImage(*user, req.Filename, s.S3Bucket, req.Object, req.Size)
	case "SellSheet":
		file, err = domain.NewSellSheet(*user, req.Filename, s.S3Bucket, req.Object, req.Size)
	case "PrintAndPlay":
		file, err = domain.NewPrintAndPlay(*user, req.Filename, s.S3Bucket, req.Object, req.Size)
	default:
		err = fmt.Errorf("invalid role '%s'", req.Role)
	}

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	// If we included a game, tie it to the game
	if req.GameID != 0 {
		// Make sure the user is allowed to edit this game
		game, err := s.GameRepository.GameOfID(req.GameID)
		if err != nil || game == nil {
			s.Logger.Error(err.Error())
			return nil, err
		}

		if !game.MayBeUpdatedBy(user) {
			s.Logger.Error(err.Error())
			return nil, err
		}

		file.BelongsTo(game)
	}

	// Save
	err = s.FileRepository.Save(file)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return file, nil
}

// ListUserFiles fetches all files belonging to the specified user
func (s *FileService) ListUserFiles(userID uint) ([]domain.File, error) {
	files, err := s.FileRepository.FilesOfUser(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return files, nil
}

// DeleteFile will remove the specified file, if the user is allowed
func (s *FileService) DeleteFile(fileID, userID uint) error {
	file, err := s.FileRepository.FileOfID(fileID)
	if err != nil {
		s.Logger.Error(err.Error())
		return err
	}

	if file == nil {
		return errors.New("file not found")
	}

	// Ensure that the current user is the uploader and deny delete if not
	if file.UploadedByID != userID {
		return errors.New("unauthorized")
	}

	// And delete
	if err := s.FileRepository.Delete(file); err != nil {
		s.Logger.Error(err.Error())
		return err
	}

	return nil
}

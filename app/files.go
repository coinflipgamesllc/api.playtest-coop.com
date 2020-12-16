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

// FileService handles file uploads/downloads and indirect interaction with S3
type FileService struct {
	FileRepository domain.FileRepository
	GameRepository domain.GameRepository
	UserRepository domain.UserRepository
	Logger         *zap.SugaredLogger
	S3Bucket       string
	S3Client       *minio.Client
}

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
func (s *FileService) CreateFile(userID uint, role, filename, object string, size int64, caption string, gameID uint) (*domain.File, error) {
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	// Create the file
	var file *domain.File
	switch role {
	case "Image":
		file, err = domain.NewImage(*user, filename, s.S3Bucket, object, size)
	case "SellSheet":
		file, err = domain.NewSellSheet(*user, filename, s.S3Bucket, object, size)
	case "PrintAndPlay":
		file, err = domain.NewPrintAndPlay(*user, filename, s.S3Bucket, object, size)
	default:
		err = fmt.Errorf("invalid role '%s'", role)
	}

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	// If we included a game, tie it to the game
	if gameID != 0 {
		// Make sure the user is allowed to edit this game
		game, err := s.GameRepository.GameOfID(gameID)
		if err != nil || game == nil {
			s.Logger.Error(err)
			return nil, err
		}

		if !game.MayBeUpdatedBy(user) {
			s.Logger.Error(err)
			return nil, err
		}

		file.BelongsTo(game)
	}

	// Save
	err = s.FileRepository.Save(file)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return file, nil
}

// ListUserFiles fetches all files belonging to the specified user
func (s *FileService) ListUserFiles(userID uint) ([]domain.File, error) {
	files, err := s.FileRepository.FilesOfUser(userID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return files, nil
}

// DeleteFile will remove the specified file, if the user is allowed
func (s *FileService) DeleteFile(fileID, userID uint) error {
	file, err := s.FileRepository.FileOfID(fileID)
	if err != nil {
		s.Logger.Error(err)
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
		s.Logger.Error(err)
		return err
	}

	return nil
}

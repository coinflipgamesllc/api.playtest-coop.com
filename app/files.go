package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/gin-gonic/gin"
)

func (s *Server) handlePresignUpload() gin.HandlerFunc {
	type request struct {
		Name      string `form:"name" binding:"required"`
		Extension string `form:"extension" binding:"required"`
	}

	type response struct {
		URL string `json:"url"`
	}

	return func(c *gin.Context) {
		// Validate request
		var req request
		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(400, serverError(err))
			return
		}

		presignedURL, err := s.s3Client.PresignedPutObject(
			context.Background(),
			s.s3Bucket,
			domain.GenerateObjectName(req.Name, req.Extension),
			time.Duration(1000)*time.Minute,
		)

		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		c.JSON(200, response{URL: presignedURL.String()})
	}
}

func (s *Server) handleCreateFile() gin.HandlerFunc {
	type request struct {
		Role     string `json:"role" binding:"required"`
		Caption  string `json:"caption"`
		Filename string `json:"filename" binding:"required"`
		Object   string `json:"object" binding:"required"`
		Size     int64  `json:"size" binding:"required"`
		GameID   uint   `json:"game"`
	}

	type response struct {
		Message string `json:"message"`
	}

	return func(c *gin.Context) {
		// Validate the request
		var req request
		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(400, serverError(err))
			return
		}

		currentUser := s.user(c)

		// Create the file
		var file *domain.File
		var err error
		switch req.Role {
		case "Image":
			file, err = domain.NewImage(*currentUser, req.Filename, s.s3Bucket, req.Object, req.Size)
		case "SellSheet":
			file, err = domain.NewSellSheet(*currentUser, req.Filename, s.s3Bucket, req.Object, req.Size)
		case "PrintAndPlay":
			file, err = domain.NewPrintAndPlay(*currentUser, req.Filename, s.s3Bucket, req.Object, req.Size)
		default:
			err = fmt.Errorf("invalid role '%s'", req.Role)
		}

		if err != nil {
			c.AbortWithStatusJSON(400, serverError(err))
			return
		}

		// If we included a game, tie it to the game
		if req.GameID != 0 {
			// Make sure the user is allowed to edit this game
			game, err := s.gameRepository.GameOfID(req.GameID)
			if err != nil || game == nil {
				c.AbortWithStatusJSON(400, serverError(errors.New("game invalid")))
				return
			}

			if !game.MayBeUpdatedBy(currentUser) {
				c.AbortWithStatusJSON(401, serverError(errors.New("you may not edit this game")))
				return
			}

			file.BelongsTo(game)
		}

		// Save
		err = s.fileRepository.Save(file)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}
	}
}

func (s *Server) handleListUserFiles() gin.HandlerFunc {
	type response struct {
		Files []domain.File `json:"files"`
	}

	return func(c *gin.Context) {
		currentUser := s.user(c)

		files, err := s.fileRepository.FilesOfUser(currentUser.ID)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		c.JSON(200, response{Files: files})
	}
}

func (s *Server) handleDeleteFile() gin.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}

	return func(c *gin.Context) {
		// Pull file by ID
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		file, err := s.fileRepository.FileOfID(uint(id))
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		if file == nil {
			c.AbortWithStatusJSON(404, serverError(errors.New("not found")))
			return
		}

		// Ensure that the current user is the uploader and deny delete if not
		currentUser := s.user(c)
		if file.UploadedByID != currentUser.ID {
			c.AbortWithStatusJSON(401, serverError(errors.New("unauthorized")))
			return
		}

		// And delete
		if err := s.fileRepository.Delete(file); err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		c.JSON(200, response{Message: "ok"})
	}
}

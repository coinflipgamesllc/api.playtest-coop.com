package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/file"
	"gorm.io/gorm"
)

// File contains all the information for a file in storage. It can be tied to a Game
type File struct {
	ID        uint           `json:"id" gorm:"primarykey" example:"123"`
	CreatedAt time.Time      `json:"created_at" example:"2020-12-11T15:29:49.321629-08:00"`
	UpdatedAt time.Time      `json:"updated_at" example:"2020-12-13T15:42:40.578904-08:00"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UploadedBy   User `json:"-"`
	UploadedByID uint `json:"-"`

	Game   *Game `json:"-"`
	GameID *uint `json:"-"`

	Role    file.Role `json:"role" example:"Image"`
	Caption string    `json:"caption" example:"What a cool image of a game!"`
	OrderBy uint      `json:"order" example:"0"`

	Filename string `json:"filename" example:"example-image.png"`
	Bucket   string `json:"-"`
	Object   string `json:"-"`
	Size     int64  `json:"-"`

	URL string `json:"url" gorm:"-" example:"https://assets.playtest-coop.com/asd9fhgaoseucgewio.png"`
}

// FileRepository defines how to interact with files in database
type FileRepository interface {
	FilesOfUser(userID uint) ([]File, error)
	FileOfID(id uint) (*File, error)
	Save(file *File) error
	Delete(file *File) error
}

// AfterFind hook for decorating the URL field for presentation
func (f *File) AfterFind(tx *gorm.DB) (err error) {
	f.decorateURL()
	return nil
}

// NewImage creates an image file
func NewImage(uploader User, filename, bucket, object string, size int64) (*File, error) {
	extension := file.ExtractExtension(filename)
	if !file.Images.Contains(extension) {
		return nil, file.InvalidExtension{ProvidedValue: extension}
	}

	file := &File{
		UploadedBy:   uploader,
		UploadedByID: uploader.ID,
		Role:         file.Image,
		Filename:     filename,
		Bucket:       bucket,
		Object:       object,
		Size:         size,
	}
	file.decorateURL()

	return file, nil
}

// NewSellSheet creates a new sellsheet pdf file
func NewSellSheet(uploader User, filename, bucket, object string, size int64) (*File, error) {
	extension := file.ExtractExtension(filename)
	if !file.Documents.Contains(extension) {
		return nil, file.InvalidExtension{ProvidedValue: extension}
	}

	file := &File{
		UploadedBy:   uploader,
		UploadedByID: uploader.ID,
		Role:         file.SellSheet,
		Filename:     filename,
		Bucket:       bucket,
		Object:       object,
		Size:         size,
	}
	file.decorateURL()

	return file, nil
}

// NewPrintAndPlay creates a new PnP pdf file
func NewPrintAndPlay(uploader User, filename, bucket, object string, size int64) (*File, error) {
	extension := file.ExtractExtension(filename)
	if !file.Documents.Contains(extension) {
		return nil, file.InvalidExtension{ProvidedValue: extension}
	}

	file := &File{
		UploadedBy:   uploader,
		UploadedByID: uploader.ID,
		Role:         file.PrintAndPlay,
		Filename:     filename,
		Bucket:       bucket,
		Object:       object,
		Size:         size,
	}
	file.decorateURL()

	return file, nil
}

// UpdateCaption will replace the caption for this file
func (f *File) UpdateCaption(newCaption string) {
	if newCaption != "" && f.Caption != newCaption {
		f.Caption = newCaption
	}
}

// UpdateOrder will simply accept the new order provided.
// Calling code is responsible for re-sorting the collection this file appears in.
func (f *File) UpdateOrder(order uint) {
	f.OrderBy = order
}

// AttachGame attaches this file to a specific game
func (f *File) AttachGame(game *Game) {
	if game != nil {
		f.GameID = &game.ID
		f.Game = game
	}
}

// GenerateObjectName will generate a unique, base64-encoded hash given a filename
// The extension will be appended to the end for helping to identify files in storage
func GenerateObjectName(filename, extension string) string {
	// Seed with name
	h := sha256.New()
	h.Write([]byte(filename))

	// Append some random junk
	r := make([]byte, 16)
	rand.Read(r)
	h.Write(r)

	// Encode as base64, but chop off the "=" at the end
	base := strings.TrimSuffix(base64.URLEncoding.EncodeToString(h.Sum(nil)), "=")

	// Append the extension
	return base + "." + extension
}

func (f *File) decorateURL() {
	f.URL = "https://assets.playtest-coop.com/" + f.Object
}

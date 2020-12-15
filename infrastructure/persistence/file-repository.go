package persistence

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"gorm.io/gorm"
)

type FileRepository struct {
	DB *gorm.DB
}

func (r *FileRepository) FilesOfUser(userID uint) ([]domain.File, error) {
	files := []domain.File{}
	result := r.DB.Find(&files, "uploaded_by_id = ?", userID)

	return files, result.Error
}

func (r *FileRepository) FileOfID(id uint) (*domain.File, error) {
	file := &domain.File{}
	result := r.DB.First(file, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return file, nil
}

// Save will upsert a file record
func (r *FileRepository) Save(file *domain.File) error {
	return r.DB.Transaction(func(db *gorm.DB) error {

		var result *gorm.DB
		if file.ID != 0 {
			result = db.Omit("UploadedBy", "Game").Save(file)
		} else {
			result = db.Omit("UploadedBy", "Game").Create(file)
		}

		return result.Error
	})
}

func (r *FileRepository) Delete(file *domain.File) error {
	result := r.DB.Delete(file)

	return result.Error
}

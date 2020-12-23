package persistence

import (
	"strings"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"gorm.io/gorm"
)

// UserRepository for a postgres db
type UserRepository struct {
	DB *gorm.DB
}

// UserOfID retrieves a user by primary key
func (r *UserRepository) UserOfID(id uint) (*domain.User, error) {
	user := &domain.User{}
	result := r.DB.First(user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) UserOfEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	result := r.DB.First(user, "email = ?", email)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) UserOfVerificationID(verificationID string) (*domain.User, error) {
	user := &domain.User{}
	result := r.DB.First(user, "verification_id = ?", verificationID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) UserOfOneTimePassword(otp string) (*domain.User, error) {
	user := &domain.User{}
	result := r.DB.First(user, "one_time_password = ?", otp)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) ListUsers(name string, limit, offset int, sort string) ([]domain.User, int, error) {
	users := []domain.User{}

	// Setup query
	query := r.DB.Model(&domain.User{})

	// Set order
	sortCol := "users.updated_at"
	sortDir := "desc"
	if sort != "" {
		parts := strings.Split(sort, ",")
		sortCol = parts[0]

		sortDir = "asc"
		if len(parts) > 1 {
			sortDir = parts[1]
		}
	}

	query = query.Order(sortCol + " " + sortDir)

	// Apply filters
	if name != "" {
		query = query.Where("users.name % ?", name)
	}

	// And run it
	var total int64
	result := query.
		Count(&total).
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return []domain.User{}, 0, result.Error
	}

	return users, int(total), nil
}

// Save will upsert a user record
func (r *UserRepository) Save(user *domain.User) error {
	var result *gorm.DB
	if user.ID != 0 {
		result = r.DB.Save(user)
	} else {
		result = r.DB.Create(user)
	}

	return result.Error
}

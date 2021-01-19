package persistence

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"gorm.io/gorm"
)

type LoginAttemptRepository struct {
	DB *gorm.DB
}

// Save will insert a login attempt record
func (r *LoginAttemptRepository) Save(la *domain.LoginAttempt) error {
	result := r.DB.Create(la)

	return result.Error
}

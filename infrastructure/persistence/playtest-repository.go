package persistence

import (
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PlaytestRepository struct {
	DB *gorm.DB
}

func (r *PlaytestRepository) PlaytestsOnDate(date time.Time, eventID uint) ([]domain.Playtest, error) {
	playtests := []domain.Playtest{}

	query := r.DB.Model(&domain.Playtest{}).
		Preload("Game").
		Preload("Game.Designers").
		Preload("Event").
		Preload("Players").
		Where("playtests.scheduled_date::date = ?::date", date)

	if eventID != 0 {
		query = query.Where("playtests.event_id = ?", eventID)
	}

	result := query.Find(&playtests)

	if result.Error != nil {
		return []domain.Playtest{}, result.Error
	}

	return playtests, nil
}

func (r *PlaytestRepository) PlaytestOfID(id uint) (*domain.Playtest, error) {
	playtest := &domain.Playtest{}
	result := r.DB.Preload(clause.Associations).First(playtest, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return playtest, nil
}

// Save will upsert an playtest record
func (r *PlaytestRepository) Save(playtest *domain.Playtest) error {
	return r.DB.Transaction(func(db *gorm.DB) error {

		var result *gorm.DB
		if playtest.ID != 0 {
			err := db.Model(playtest).Association("Players").Replace(playtest.Players)
			if err != nil {
				return err
			}

			result = db.Omit(clause.Associations).Save(playtest)
		} else {
			result = db.Omit(clause.Associations).Create(playtest)
		}

		return result.Error
	})
}

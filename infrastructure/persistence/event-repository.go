package persistence

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository struct {
	DB *gorm.DB
}

func (r *EventRepository) ListEvents() ([]domain.Event, error) {
	events := []domain.Event{}

	result := r.DB.Find(&events)

	if result.Error != nil {
		return []domain.Event{}, result.Error
	}

	return events, nil
}

func (r *EventRepository) EventOfID(id uint) (*domain.Event, error) {
	event := &domain.Event{}
	result := r.DB.Preload(clause.Associations).First(event, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return event, nil
}

// Save will upsert an event record
func (r *EventRepository) Save(event *domain.Event) error {
	return r.DB.Transaction(func(db *gorm.DB) error {

		var result *gorm.DB
		if event.ID != 0 {
			err := db.Model(event).Association("Facilitators").Replace(event.Facilitators)
			if err != nil {
				return err
			}

			result = db.Omit(clause.Associations).Save(event)
		} else {
			result = db.Omit(clause.Associations).Create(event)
		}

		return result.Error
	})
}

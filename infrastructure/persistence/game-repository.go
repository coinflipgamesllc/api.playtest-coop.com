package persistence

import (
	"math"
	"strings"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GameRepository struct {
	DB *gorm.DB
}

func (r *GameRepository) ListGames(title, status, designer string, owner uint, playerCount, age, playtime, limit, offset int, sort string) ([]domain.Game, int, error) {
	games := []domain.Game{}

	// Setup query
	query := r.DB.Model(&domain.Game{}).Preload("Designers").Preload("Files", func(db *gorm.DB) *gorm.DB {
		return db.Where("files.role = 'Image'").Order("files.order_by ASC")
	})

	// Set order
	sortCol := "games.updated_at"
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
	if title != "" {
		query = query.Where("games.title % ?", title)
	}

	if owner == 0 {
		if status != "" {
			query = query.Where("games.status = ?", status)
		} else {
			query = query.Where("games.status != 'Archived'")
		}

		if designer != "" {
			designerQuery := r.DB.Select("game_designers.game_id").Table("game_designers").Joins("join users on users.id = game_designers.user_id").Where("users.name % ?", designer)
			query = query.Where("games.id in (?)", designerQuery)
		}
	} else {
		ownerQuery := r.DB.Select("game_designers.game_id").Table("game_designers").Joins("join users on users.id = ?", owner)
		query = query.Where("games.id in (?)", ownerQuery)
	}

	if playerCount != 0 {
		query = query.Where("? BETWEEN games.min_players AND games.max_players", playerCount)
	}

	if age != 0 {
		query = query.Where("games.min_age <= ?", age)
	}

	if playtime != 0 {
		threshold := int64(math.Round(float64(playtime) * 0.20))
		query = query.Where("? BETWEEN games.estimated_playtime - ? AND games.estimated_playtime + ?", playtime, threshold, threshold)
	}

	// And run it
	var total int64
	result := query.
		Count(&total).
		Limit(limit).
		Offset(offset).
		Find(&games)

	if result.Error != nil {
		return []domain.Game{}, 0, result.Error
	}

	return games, int(total), nil
}

func (r *GameRepository) GameOfID(id uint) (*domain.Game, error) {
	game := &domain.Game{}
	result := r.DB.Preload(clause.Associations).Preload("Files", func(db *gorm.DB) *gorm.DB {
		return db.Order("files.order_by ASC")
	}).First(game, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return game, nil
}

func (r *GameRepository) RulesOfGame(id uint) ([]game.RulesSection, error) {
	rules := []game.RulesSection{}

	result := r.DB.Find(&rules).Where("game_id = ?", id).Order("rules_section.order_by ASC")

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return rules, nil
}

// Save will upsert a game record
func (r *GameRepository) Save(game *domain.Game) error {
	return r.DB.Transaction(func(db *gorm.DB) error {

		var result *gorm.DB
		if game.ID != 0 {
			err := db.Model(game).Association("Designers").Replace(game.Designers)
			if err != nil {
				return err
			}

			result = db.Omit("Designers").Save(game)
		} else {
			result = db.Omit("Designers.*").Create(game)
		}

		return result.Error
	})
}

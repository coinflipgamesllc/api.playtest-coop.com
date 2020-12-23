package app

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"go.uber.org/zap"
)

type (
	// UserService handles general interactions with users
	UserService struct {
		UserRepository domain.UserRepository
		Logger         *zap.Logger
	}

	// Request DTOs

	// ListUsersRequest query params
	ListUsersRequest struct {
		Name   string `form:"name" example:"New User"`
		Limit  int    `form:"limit" example:"100"`
		Offset int    `form:"offset" example:"50"`
		Sort   string `form:"sort" example:"name,desc"`
	}

	// Response DTOs

	// ListUsersResponse paginated users list
	ListUsersResponse struct {
		Users  []domain.User `json:"users"`
		Total  int           `json:"total" example:"1000"`
		Limit  int           `json:"limit" example:"100"`
		Offset int           `json:"offset" example:"50"`
	}
)

// ListUsers returns all users matching the specified query. The results are paginated
func (s *UserService) ListUsers(req *ListUsersRequest) ([]domain.User, int, error) {
	// Limit our limit
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Fetch users
	users, total, err := s.UserRepository.ListUsers(
		req.Name,
		req.Limit,
		req.Offset,
		req.Sort,
	)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, 0, err
	}

	return users, total, nil
}

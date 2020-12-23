package controller

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/gin-gonic/gin"
)

// UserController handles /users routes
type UserController struct {
	UserService *app.UserService
}

// ListUsers list users matching the query with pagination
// @Summary List users matching the query with pagination
// @Accept json
// @Produce json
// @Param query query app.ListUsersRequest false "Filters for users"
// @Success 200 {object} app.ListUsersResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags users
// @Router /users [get]
func (t *UserController) ListUsers(c *gin.Context) {
	// Validate request
	var req app.ListUsersRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	// Fetch users
	users, total, err := t.UserService.ListUsers(&req)

	if err != nil {
		serverErrorResponse(c, "failed to fetch users")
		return
	}

	c.JSON(200, app.ListUsersResponse{Users: users, Total: total, Limit: req.Limit, Offset: req.Offset})
}

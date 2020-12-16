package controller

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/gin-gonic/gin"
)

// AuthController handles /auth routes
type AuthController struct {
	AuthService *app.AuthService
}

// GetUser retrieves the authenticated user
// @Summary Retrieve the authenticated user
// @Description The authentication token includes the user's ID as the subject. We extract that and use it to pull the user from the database.
// @Produce json
// @Success 200 {object} app.UserResponse
// @Failure 401 {object} UnauthorizedResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags auth
// @Router /auth/user [get]
func (t *AuthController) GetUser(c *gin.Context) {
	userID := userID(c)

	// Fetch the user
	user, err := t.AuthService.FetchUser(userID)
	if err != nil {
		serverErrorResponse(c, "failed to fetch user")
		return
	}

	if user == nil {
		unauthorizedResponse(c)
		return
	}

	c.JSON(200, app.UserResponse{User: user})
}

// UpdateUser updates authenticated user
// @Summary Update authenticated user
// @Accept json
// @Produce json
// @Param params body app.UpdateUserRequest false "User data to update"
// @Success 200 {object} app.UserResponse
// @Failure 401 {object} UnauthorizedResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags auth
// @Router /auth/user [put]
func (t *AuthController) UpdateUser(c *gin.Context) {
	userID := userID(c)

	// Validate request
	var req app.UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	user, err := t.AuthService.UpdateUser(&req, userID)
	if err != nil {
		serverErrorResponse(c, "failed to update user")
		return
	}

	if user == nil {
		notFoundResponse(c, "user not found")
		return
	}

	c.JSON(200, app.UserResponse{User: user})
}

// Signup creates and authenticates a new user
// @Summary Create and authenticates a new user
// @Accept json
// @Produce json
// @Param credentials body app.SignupRequest true "User name, email, and password"
// @Success 201 {object} app.UserTokenResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/signup [post]
func (t *AuthController) Signup(c *gin.Context) {
	// Validate request
	var req app.SignupRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	user, at, rt, err := t.AuthService.Signup(req.Name, req.Email, req.Password)
	if err != nil {
		requestErrorResponse(c, "failed to create account")
		return
	}

	c.JSON(201, app.UserTokenResponse{User: user, AccessToken: at, RefreshToken: rt})
}

// Login authenticates a user
// @Summary Authenticate a user
// @Accept json
// @Produce json
// @Param credentials body app.LoginRequest true "User email/password combo"
// @Success 200 {object} app.UserTokenResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/login [post]
func (t *AuthController) Login(c *gin.Context) {
	// Validate request
	var req app.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Attempt to log in
	user, at, rt, err := t.AuthService.Login(req.Email, req.Password)
	if err != nil {
		requestErrorResponse(c, "failed to log in")
		return
	}

	c.JSON(200, app.UserTokenResponse{User: user, AccessToken: at, RefreshToken: rt})
}

// RefreshToken regenerates the access token and refresh token, given a valid refresh token.
// @Summary Regenerate the access token and refresh token, given a valid refresh token.
// @Accept json
// @Produce json
// @Param refresh_token body app.RefreshTokenRequest true "Refresh token originally acquired from /auth/token, /auth/signup, or /auth/login"
// @Success 200 {object} app.TokenResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/token [post]
func (t *AuthController) RefreshToken(c *gin.Context) {
	// Validate request
	var req app.RefreshTokenRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Validate token
	at, rt, err := t.AuthService.RefreshToken(req.RefreshToken)
	if err != nil {
		requestErrorResponse(c, "failed to refresh tokens")
		return
	}

	c.JSON(200, app.TokenResponse{AccessToken: at, RefreshToken: rt})
}

// VerifyEmail verifies that a user's email address is valid. A link is sent to their email and clicking it takes them here.
// this route isn't technically part of the API and does not serve JSON like the other routes.
func (t *AuthController) VerifyEmail(c *gin.Context) {
	// Fetch user by ID
	id := c.Param("id")
	err := t.AuthService.VerifyEmail(id)
	if err != nil {
		c.HTML(500, "500.html", gin.H{"error": err.Error()})
		return
	}

	// Send em home
	c.Redirect(307, "https://playtest-coop.com")
}

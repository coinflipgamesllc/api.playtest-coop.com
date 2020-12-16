package controller

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/gin-gonic/gin"
)

// AuthController handles /auth routes
type AuthController struct {
	AuthService *app.AuthService
}

// GetUserResponse wraps User object
type GetUserResponse struct {
	User *domain.User `json:"user"`
}

// GetUser retrieves the authenticated user
// @Summary Retrieve the authenticated user
// @Description The authentication token includes the user's ID as the subject. We extract that and use it to pull the user from the database.
// @Produce json
// @Success 200 {object} GetUserResponse
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

	c.JSON(200, GetUserResponse{User: user})
}

// UpdateUserRequest definition for updating a user
type UpdateUserRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2" example:"User McUserton"`
	Email       string `json:"email" binding:"omitempty,email" example:"user@example.com"`
	NewPassword string `json:"new_password" binding:"omitempty,nefield=OldPassword,min=10" example:"AVerySecurePassword123!"`
	OldPassword string `json:"old_password" binding:"omitempty" example:"NotASecurePassword"`
	Pronouns    string `json:"pronouns" binding:"omitempty,contains=/" example:"they/them"`
}

// UpdateUser updates authenticated user
// @Summary Update authenticated user
// @Accept json
// @Produce json
// @Param params body UpdateUserRequest false "User data to update"
// @Success 200 {object} GetUserResponse
// @Failure 401 {object} UnauthorizedResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags auth
// @Router /auth/user [put]
func (t *AuthController) UpdateUser(c *gin.Context) {
	userID := userID(c)

	// Validate request
	var req UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	user, err := t.AuthService.UpdateUser(userID, req.Name, req.Email, req.NewPassword, req.OldPassword, req.Pronouns)
	if err != nil {
		serverErrorResponse(c, "failed to update user")
		return
	}

	if user == nil {
		notFoundResponse(c, "user not found")
		return
	}

	c.JSON(200, GetUserResponse{User: user})
}

// SignupRequest params for signing up for a new account
type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2" example:"User McUserton"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=10" example:"AVerySecurePassword123!"`
}

// UserTokenResponse includes user object with access and refresh tokens
type UserTokenResponse struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgwNTY5NzksIm5hbWUiOiJSb2IgTmV3dG9uIiwic3ViIjoxfQ.KKUtLne51DqBPqQxZZmCFsjsGAeYRukZNcXCx6IpLN8"`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"`
}

// Signup creates and authenticates a new user
// @Summary Create and authenticates a new user
// @Accept json
// @Produce json
// @Param credentials body SignupRequest true "User name, email, and password"
// @Success 201 {object} UserTokenResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/signup [post]
func (t *AuthController) Signup(c *gin.Context) {
	// Validate request
	var req SignupRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	user, at, rt, err := t.AuthService.Signup(req.Name, req.Email, req.Password)
	if err != nil {
		requestErrorResponse(c, "failed to create account")
		return
	}

	c.JSON(201, UserTokenResponse{User: user, AccessToken: at, RefreshToken: rt})
}

// LoginRequest params for logging in
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=10" example:"AVerySecurePassword123!"`
}

// Login authenticates a user
// @Summary Authenticate a user
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User email/password combo"
// @Success 200 {object} UserTokenResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/login [post]
func (t *AuthController) Login(c *gin.Context) {
	// Validate request
	var req LoginRequest
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

	c.JSON(200, UserTokenResponse{User: user, AccessToken: at, RefreshToken: rt})
}

// RefreshTokenRequest param for refreshing access tokens
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"`
}

// TokenResponse wrapper for access and refresh tokens
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgwNTY5NzksIm5hbWUiOiJSb2IgTmV3dG9uIiwic3ViIjoxfQ.KKUtLne51DqBPqQxZZmCFsjsGAeYRukZNcXCx6IpLN8"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"`
}

// RefreshToken regenerates the access token and refresh token, given a valid refresh token.
// @Summary Regenerate the access token and refresh token, given a valid refresh token.
// @Accept json
// @Produce json
// @Param refresh_token body RefreshTokenRequest true "Refresh token originally acquired from /auth/token, /auth/signup, or /auth/login"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/token [post]
func (t *AuthController) RefreshToken(c *gin.Context) {
	// Validate request
	var req RefreshTokenRequest
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

	c.JSON(200, TokenResponse{AccessToken: at, RefreshToken: rt})
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

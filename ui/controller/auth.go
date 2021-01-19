package controller

import (
	"errors"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/gin-contrib/sessions"
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
// @Failure 400 {object} ValidationErrorResponse
// @Failure 401 {object} UnauthorizedResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags auth
// @Router /auth/user [put]
func (t *AuthController) UpdateUser(c *gin.Context) {
	userID := userID(c)

	// Validate request
	var req app.UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	user, err := t.AuthService.UpdateUser(&req, userID)
	if err != nil {
		if errors.Is(err, domain.UserNotFound{}) {
			notFoundResponse(c, err.Error())
			return
		}

		serverErrorResponse(c, "failed to update user")
		return
	}

	c.JSON(200, app.UserResponse{User: user})
}

// RequestResetPassword sends a password reset email to the specified email
// @Summary Send a password reset email to the specified email
// @Accept json
// @Produce json
// @Param email body app.ResetPasswordRequest true "User email to request a password reset for"
// @Success 200 {object} AckResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/reset-password [post]
func (t *AuthController) RequestResetPassword(c *gin.Context) {
	// Validate request
	var req app.ResetPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	err := t.AuthService.RequestResetPassword(req.Email)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	ackResponse(c)
}

// Signup creates and authenticates a new user
// @Summary Create and authenticates a new user
// @Accept json
// @Produce json
// @Param credentials body app.SignupRequest true "User name, email, and password"
// @Success 201 {object} app.UserResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/signup [post]
func (t *AuthController) Signup(c *gin.Context) {
	// Validate request
	var req app.SignupRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	user, err := t.AuthService.Signup(req.Name, req.Email, req.Password, c.ClientIP())
	if err != nil {
		requestErrorResponse(c, "failed to create account")
		return
	}

	// Cookie time!
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.JSON(201, app.UserResponse{User: user})
}

// Login authenticates a user
// @Summary Authenticate a user
// @Accept json
// @Produce json
// @Param credentials body app.LoginRequest true "User email/password combo"
// @Success 200 {object} app.UserResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/login [post]
func (t *AuthController) Login(c *gin.Context) {
	// Validate request
	var req app.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	// Attempt to log in
	user, err := t.AuthService.Login(req.Email, req.Password, c.ClientIP())
	if err != nil {
		if errors.Is(err, domain.UserNotFound{}) {
			notFoundResponse(c, err.Error())
			return
		}

		if errors.Is(err, domain.CredentialsIncorrect{}) {
			notFoundResponse(c, err.Error())
			return
		}

		requestErrorResponse(c, err.Error())
		return
	}

	// Cookie time!
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.JSON(200, app.UserResponse{User: user})
}

// Logout ends an authenticated session
// @Summary End an authenticated session
// @Produce json
// @Success 200 {object} AckResponse
// @Failure 400 {object} RequestErrorResponse
// @Tags auth
// @Router /auth/logout [get]
func (t *AuthController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	session.Save()

	ackResponse(c)
}

// Non-API routes

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

// ResetPassword confirms a reset password request and redirects the user to a
// page to set their actual password.
func (t *AuthController) ResetPassword(c *gin.Context) {
	// Pull otp by ID
	otp := c.Param("otp")

	err := t.AuthService.ResetPassword(otp)
	if err != nil {
		c.HTML(500, "500.html", gin.H{"error": "Invalid password reset link (maybe it was already used?)"})
		return
	}

	// Make em set a real password
	c.Redirect(307, "https://playtest-coop.com/set-password?p="+otp)
}

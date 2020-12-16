package app

import (
	"errors"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type (
	// AuthService handles both authentication and authorization
	AuthService struct {
		AuthToken      string
		Logger         *zap.SugaredLogger
		UserRepository domain.UserRepository
	}

	// Request DTOs

	// UpdateUserRequest definition for updating a user
	UpdateUserRequest struct {
		Name        string `json:"name" binding:"omitempty,min=2" example:"User McUserton"`
		Email       string `json:"email" binding:"omitempty,email" example:"user@example.com"`
		NewPassword string `json:"new_password" binding:"omitempty,nefield=OldPassword,min=10" example:"AVerySecurePassword123!"`
		OldPassword string `json:"old_password" binding:"omitempty" example:"NotASecurePassword"`
		Pronouns    string `json:"pronouns" binding:"omitempty,contains=/" example:"they/them"`
	}

	// SignupRequest params for signing up for a new account
	SignupRequest struct {
		Name     string `json:"name" binding:"required,min=2" example:"User McUserton"`
		Email    string `json:"email" binding:"required,email" example:"user@example.com"`
		Password string `json:"password" binding:"required,min=10" example:"AVerySecurePassword123!"`
	}

	// LoginRequest params for logging in
	LoginRequest struct {
		Email    string `json:"email" binding:"required,email" example:"user@example.com"`
		Password string `json:"password" binding:"required,min=10" example:"AVerySecurePassword123!"`
	}

	// RefreshTokenRequest param for refreshing access tokens
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"`
	}

	// Response DTOs

	// UserResponse wraps User object
	UserResponse struct {
		User *domain.User `json:"user"`
	}

	// UserTokenResponse includes user object with access and refresh tokens
	UserTokenResponse struct {
		User         *domain.User `json:"user"`
		AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgwNTY5NzksIm5hbWUiOiJSb2IgTmV3dG9uIiwic3ViIjoxfQ.KKUtLne51DqBPqQxZZmCFsjsGAeYRukZNcXCx6IpLN8"`
		RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"`
	}

	// TokenResponse wrapper for access and refresh tokens
	TokenResponse struct {
		AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgwNTY5NzksIm5hbWUiOiJSb2IgTmV3dG9uIiwic3ViIjoxfQ.KKUtLne51DqBPqQxZZmCFsjsGAeYRukZNcXCx6IpLN8"`
		RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"`
	}
)

func (s *AuthService) generateTokensForUser(user *domain.User) (string, string, error) {
	accessToken := jwt.New(jwt.GetSigningMethod("HS256"))
	accessToken.Claims = jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	}

	at, err := accessToken.SignedString([]byte(s.AuthToken))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.GetSigningMethod("HS256"))
	refreshToken.Claims = jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	rt, err := refreshToken.SignedString([]byte(s.AuthToken))
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

// UpdateUser will update the user with the specified values
func (s *AuthService) UpdateUser(userID uint, name, email, newPassword, oldPassword, pronouns string) (*domain.User, error) {
	user, err := s.FetchUser(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	// Update the user
	if name != "" {
		user.Rename(name)
	}

	if email != "" {
		user.ChangeEmail(email)
	}

	if newPassword != "" && oldPassword != "" {
		err := user.ChangePassword(newPassword, oldPassword)
		if err != nil {
			s.Logger.Error(err.Error(), "user", userID)
			return nil, err
		}
	}

	if pronouns != "" {
		user.SetPronouns(pronouns)
	}

	// Save changes
	err = s.UserRepository.Save(user)
	if err != nil {
		s.Logger.Error(err.Error(), "user", userID)
		return nil, err
	}

	return user, nil
}

// Signup will create a new account
func (s *AuthService) Signup(name, email, password string) (*domain.User, string, string, error) {
	user, err := domain.NewUser(name, email, password)
	if err != nil {
		s.Logger.Error(err)
		return nil, "", "", err
	}

	// Save
	err = s.UserRepository.Save(user)
	if err != nil {
		s.Logger.Error(err)
		return nil, "", "", err
	}

	// Generate tokens
	at, rt, err := s.generateTokensForUser(user)
	if err != nil {
		s.Logger.Error(err)
		return nil, "", "", err
	}

	return user, at, rt, nil
}

// Login attempts to log a user into their account
func (s *AuthService) Login(email, password string) (*domain.User, string, string, error) {
	// Retrieve user
	user, err := s.UserRepository.UserOfEmail(email)
	if err != nil {
		s.Logger.Error(err)
		return nil, "", "", err
	}

	if user == nil {
		s.Logger.Error(err)
		return nil, "", "", err
	}

	// Verify password
	ok, err := user.ValidPassword(password)
	if err != nil {
		s.Logger.Error(err)
		return nil, "", "", err
	}

	if !ok {
		return nil, "", "", errors.New("email and password combination not found")
	}

	// Generate tokens for future requests
	at, rt, err := s.generateTokensForUser(user)
	if err != nil {
		s.Logger.Error(err)
		return nil, "", "", err
	}

	return user, at, rt, nil
}

// RefreshToken will regenerate access and refresh tokens given a valid refresh token
func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	// Validate token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.AuthToken), nil
	})

	if err != nil {
		s.Logger.Error(err)
		return "", "", err
	}

	// Extract and validate that the user account is still active
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["sub"]

		user, err := s.UserRepository.UserOfID(uint(id.(float64)))
		if err != nil {
			s.Logger.Error(err)
			return "", "", err
		}

		if user == nil {
			return "", "", err
		}

		// Generate a new token pair
		at, rt, err := s.generateTokensForUser(user)
		if err != nil {
			s.Logger.Error(err)
			return "", "", err
		}

		return at, rt, nil
	}

	return "", "", errors.New("unauthorized")
}

// VerifyEmail will check for a verify id in the database and mark the corresponding user's email as valid
func (s *AuthService) VerifyEmail(id string) error {
	// Fetch user by ID
	user, err := s.UserRepository.UserOfVerificationID(id)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	if user == nil {
		return nil
	}

	// Mark verified and save
	user.VerifyEmail()
	if err := s.UserRepository.Save(user); err != nil {
		s.Logger.Error(err)
		return err
	}

	return nil
}

// FetchUser returns the user with the provided ID
func (s *AuthService) FetchUser(id uint) (*domain.User, error) {
	return s.UserRepository.UserOfID(id)
}

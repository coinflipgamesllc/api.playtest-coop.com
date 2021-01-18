package app

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"go.uber.org/zap"
)

type (
	// AuthService handles both authentication and authorization
	AuthService struct {
		AuthToken      string
		Logger         *zap.Logger
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
		Color       string `json:"color" binding:"omitempty,hexcolor" example:"#2a9d8f"`
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

	// ResetPasswordRequest params for requesting a password reset
	ResetPasswordRequest struct {
		Email string `json:"email" binding:"required,email" example:"user@example.com"`
	}

	// Response DTOs

	// UserResponse wraps User object
	UserResponse struct {
		User *domain.User `json:"user"`
	}
)

// UpdateUser will update the user with the specified values
func (s *AuthService) UpdateUser(req *UpdateUserRequest, userID uint) (*domain.User, error) {
	user, err := s.FetchUser(userID)
	if err != nil {
		return nil, domain.GenericServerError{}
	}

	if user == nil {
		return nil, domain.UserNotFound{ProvidedID: userID}
	}

	// Update the user
	if req.Name != "" {
		user.Rename(req.Name)
	}

	if req.Email != "" {
		user.ChangeEmail(req.Email)
	}

	if req.NewPassword != "" && req.OldPassword != "" {
		err := user.ChangePassword(req.NewPassword, req.OldPassword)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, domain.GenericServerError{}
		}
	}

	if req.Pronouns != "" {
		user.SetPronouns(req.Pronouns)
	}

	if req.Color != "" {
		user.SetColor(req.Color)
	}

	// Save changes
	err = s.UserRepository.Save(user)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, domain.GenericServerError{}
	}

	// Decorate the email for this response
	user.Email = user.Account.Email

	return user, nil
}

// RequestResetPassword will send a password reset email to the specified user
func (s *AuthService) RequestResetPassword(email string) error {
	// Retrieve user
	user, err := s.UserRepository.UserOfEmail(email)
	if err != nil {
		s.Logger.Error(err.Error())
		return domain.GenericServerError{}
	}

	if user == nil {
		return domain.UserNotFound{ProvidedEmail: email}
	}

	// Request reset password & save
	user.RequestResetPassword()
	err = s.UserRepository.Save(user)
	if err != nil {
		s.Logger.Error(err.Error())
		return domain.GenericServerError{}
	}

	return nil
}

// ResetPassword actually resets the user's password
func (s *AuthService) ResetPassword(otp string) error {
	// Retrieve user
	user, err := s.UserRepository.UserOfOneTimePassword(otp)
	if err != nil {
		s.Logger.Error(err.Error())
		return domain.GenericServerError{}
	}

	if user == nil {
		return domain.UserNotFound{}
	}

	// Actually reset password & save
	err = user.ResetPassword(otp)
	if err != nil {
		s.Logger.Error(err.Error())
		return domain.GenericServerError{}
	}

	err = s.UserRepository.Save(user)
	if err != nil {
		s.Logger.Error(err.Error())
		return domain.GenericServerError{}
	}

	return nil
}

// Signup will create a new account
func (s *AuthService) Signup(name, email, password string) (*domain.User, error) {
	user, err := domain.NewUser(name, email, password)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, domain.GenericServerError{}
	}

	// Save
	err = s.UserRepository.Save(user)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, domain.GenericServerError{}
	}

	// Decorate the email for this response
	user.Email = user.Account.Email

	return user, nil
}

// Login attempts to log a user into their account
func (s *AuthService) Login(email, password string) (*domain.User, error) {
	// Retrieve user
	user, err := s.UserRepository.UserOfEmail(email)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, domain.GenericServerError{}
	}

	if user == nil {
		return nil, domain.UserNotFound{ProvidedEmail: email}
	}

	// Verify password
	ok, err := user.ValidPassword(password)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, domain.GenericServerError{}
	}

	if !ok {
		return nil, domain.CredentialsIncorrect{}
	}

	// Decorate the email for this response
	user.Email = user.Account.Email

	return user, nil
}

// VerifyEmail will check for a verify id in the database and mark the corresponding user's email as valid
func (s *AuthService) VerifyEmail(id string) error {
	// Fetch user by ID
	user, err := s.UserRepository.UserOfVerificationID(id)
	if err != nil {
		s.Logger.Error(err.Error())
		return domain.GenericServerError{}
	}

	if user == nil {
		return domain.UserNotFound{}
	}

	// Mark verified and save
	user.VerifyEmail()
	if err := s.UserRepository.Save(user); err != nil {
		s.Logger.Error(err.Error())
		return domain.GenericServerError{}
	}

	return nil
}

// FetchUser returns the user with the provided ID
func (s *AuthService) FetchUser(id uint) (*domain.User, error) {
	user, err := s.UserRepository.UserOfID(id)
	if err != nil {
		return nil, err
	}

	// Decorate the email for this response
	user.Email = user.Account.Email

	return user, nil
}

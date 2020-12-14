package domain

import (
	"errors"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/user"
	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure/pubsub"
	"gorm.io/gorm"
)

// User represents a designer, tester, publisher, etc user of the system
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name     string       `json:"name"`
	Account  user.Account `json:"-" gorm:"embedded"`
	Pronouns string       `json:"pronouns"`
}

// UserRepository defines how to interact with the user in database
type UserRepository interface {
	UserOfID(uint) (*User, error)
	UserOfEmail(string) (*User, error)
	UserOfVerificationID(string) (*User, error)
	Save(*User) error
}

func userCreated(u *User) Event {
	return Event{
		Name: "User/Created",
		Data: map[string]interface{}{
			"id":             u.ID,
			"name":           u.Name,
			"email":          u.Account.Email,
			"verificationID": u.Account.VerificationID,
		},
	}
}

func userEmailUnverified(u *User) Event {
	return Event{
		Name: "User/EmailUnverified",
		Data: map[string]interface{}{
			"id":             u.ID,
			"name":           u.Name,
			"email":          u.Account.Email,
			"verificationID": u.Account.VerificationID,
		},
	}
}

// NewUser creates a new user with the specified name, email, and password.
func NewUser(name, email, password string) (*User, error) {
	account, err := user.NewAccount(email, password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:    name,
		Account: *account,
	}, nil
}

// VerifyEmail marks the user's email as verified
func (u *User) VerifyEmail() {
	u.Account.VerifyEmail()
}

// ValidPassword returns true if the provided password matches the account password
func (u *User) ValidPassword(password string) (bool, error) {
	return u.Account.ValidPassword(password)
}

// Rename updates the user's name
func (u *User) Rename(newName string) {
	if u.Name != newName && newName != "" {
		u.Name = newName
	}
}

// ChangeEmail updates the user's email
func (u *User) ChangeEmail(newEmail string) {
	if u.Account.Email != newEmail && newEmail != "" {
		u.Account.Email = newEmail
		u.Account.Verified = false
	}
}

// ChangePassword updates the user's password
func (u *User) ChangePassword(newPassword, oldPassword string) error {
	valid, err := u.Account.ValidPassword(oldPassword)
	if !valid {
		return errors.New("password mismatch")
	}
	if err != nil {
		return err
	}

	account, err := user.NewAccount(u.Account.Email, newPassword)
	if err != nil {
		return err
	}

	u.Account = *account
	u.Account.Verified = true // We don't need to revalidate email

	return nil
}

// SetPronouns updates the user's pronouns
func (u *User) SetPronouns(newPronouns string) {
	u.Pronouns = newPronouns
}

// AfterCreate hook for sending welcome emails
func (u *User) AfterCreate(tx *gorm.DB) error {
	event := userCreated(u)
	pubsub.Instance.Publish(event.Name, event.Data)

	return nil
}

// AfterUpdate hook for sending verification emails
func (u *User) AfterUpdate(tx *gorm.DB) error {
	if !u.Account.Verified {
		event := userEmailUnverified(u)
		pubsub.Instance.Publish(event.Name, event.Data)
	}

	return nil
}

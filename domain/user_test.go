package domain

import (
	"testing"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/user"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("Name", "email@example.com", "Password")

	if err != nil {
		t.Errorf("Error encountered while creating user: %s", err)
	}

	if u.Name != "Name" {
		t.Error("Name not set on new user")
	}

	if u.Account.Email != "email@example.com" {
		t.Error("Email is not set on new user")
	}

	if u.Account.Password == "Password" || u.Account.Password == "" {
		t.Error("Password not set or not hashed on new user")
	}

	if u.Pronouns != "" {
		t.Error("Users do not get default pronouns on new user")
	}
}

func TestValidPassword(t *testing.T) {
	u, _ := NewUser("Name", "email@example.com", "Password")
	ok, err := u.ValidPassword("Password")

	if err != nil {
		t.Errorf("Error encountered while checking password: %s", err)
	}

	if !ok {
		t.Error("Valid passwords not matching")
	}
}

func TestRenameUser(t *testing.T) {
	var tests = []struct {
		user         *User
		newName      string
		expectedName string
	}{
		{&User{Name: "Original Name"}, "New Name", "New Name"},
		{&User{Name: "Original Name"}, "", "Original Name"},
		{&User{Name: "Original Name"}, "Original Name", "Original Name"},
	}

	for _, tt := range tests {
		tt.user.Rename(tt.newName)
		actual := tt.user.Name
		if tt.expectedName != actual {
			t.Errorf("Rename incorrect")
		}
	}
}

func TestChangeEmail(t *testing.T) {
	u, _ := NewUser("Name", "email@example.com", "Password")
	u.VerifyEmail()

	u.ChangeEmail("new@email.com")
	if u.Account.Email != "new@email.com" {
		t.Error("Failed to change email")
	}

	if u.Account.Verified == true {
		t.Error("Changing email requires revalidation")
	}
}

func TestChangePassword(t *testing.T) {
	u, _ := NewUser("Name", "email@example.com", "Password")
	u.VerifyEmail()

	originalHash := u.Account.Password

	err := u.ChangePassword("New Password", "Password")
	if err != nil {
		t.Errorf("Encountered error on password change: %s", err)
	}

	if originalHash == u.Account.Password {
		t.Error("Changing password didn't actually change the password")
	}

	if u.Account.Verified == false {
		t.Error("Changing password shouldn't reset email verification")
	}

	err = u.ChangePassword("Another new password", "Definitely not correct")
	if err == nil {
		t.Error("Changing password requires _correct_ old password")
	}
	if _, ok := err.(user.PasswordMismatch); !ok {
		t.Error("PasswordMismatch error is required in specific")
	}
}

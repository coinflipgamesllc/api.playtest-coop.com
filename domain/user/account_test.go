package user

import (
	"math/rand"
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password := "starting_password"
	h, err := hash(password)
	if err != nil {
		t.Errorf("Failed to hash password: %s", err.Error())
	}

	equal, err := compare("starting_password", h)
	if err != nil {
		t.Errorf("Failed to compare passwords: %s", err.Error())
	}

	if !equal {
		t.Error("Passwords should produce the same hash and be equal.")
	}

	notEqual, err := compare("different_password", h)
	if err != nil {
		t.Errorf("Failed to compare passwords: %s", err.Error())
	}

	if notEqual {
		t.Error("Passwords should not produce the same hash and not be equal.")
	}
}

func TestVerifyEmail(t *testing.T) {
	account := Account{Verified: false, VerificationID: "Not Verified"}
	account.VerifyEmail()

	if account.Verified != true || account.VerificationID != "" {
		t.Error("Verification should clear the verification ID and set the account to verified")
	}
}

func TestAddOneTimePassword(t *testing.T) {
	rand.Seed(42)
	a := Account{}
	a.AddOneTimePassword()

	expected := "HRukpTTueZPtNeuvunhuksqVGzAdxlgghEjkMVeZJpmKqakmTRgKfBSWYjUNGkdm"
	actual := a.OneTimePassword
	if expected != actual {
		t.Errorf("Password reset should generate random OTP matching '%s', got '%s'", expected, actual)
	}
}

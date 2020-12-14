package user

import "testing"

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

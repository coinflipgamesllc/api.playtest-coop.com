package domain

import "fmt"

// GenericServerError error
type GenericServerError struct{}

func (e GenericServerError) Error() string {
	return "something went wrong"
}

// UserNotFound error
type UserNotFound struct {
	ProvidedID    uint
	ProvidedEmail string
}

func (e UserNotFound) Error() string {
	if e.ProvidedID != 0 {
		return fmt.Sprintf("user '%d' not found", e.ProvidedID)
	}
	if e.ProvidedEmail != "" {
		return fmt.Sprintf("user with email '%s' not found", e.ProvidedEmail)
	}

	return "user not found"
}

// CredentialsIncorrect error
type CredentialsIncorrect struct{}

func (e CredentialsIncorrect) Error() string {
	return "email and password combination not found"
}

// OneTimePasswordIncorrect error
type OneTimePasswordIncorrect struct{}

func (e OneTimePasswordIncorrect) Error() string {
	return "one-time-use password expired or invalid"
}

// Unauthorized error
type Unauthorized struct{}

func (e Unauthorized) Error() string {
	return "you are not allowed to do that"
}

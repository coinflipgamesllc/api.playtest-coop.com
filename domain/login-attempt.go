package domain

import "time"

// LoginAttempt tracks logins by email/ip address and whether they were successful
type LoginAttempt struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	Email      string
	Successful bool
	IP         string
}

// LoginAttemptRepository defines how to interact with the user in database
type LoginAttemptRepository interface {
	Save(*LoginAttempt) error
}

// RecordLoginAttempt creates the log entry for the attempt. It needs to be persisted
func RecordLoginAttempt(email, ip string, successful bool) *LoginAttempt {
	return &LoginAttempt{
		Email:      email,
		IP:         ip,
		Successful: successful,
	}
}

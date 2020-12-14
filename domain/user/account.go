package user

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

// Account tracks a user's authentication credentials
type Account struct {
	Email          string `json:"email" gorm:"unique"`
	Password       string `json:"-"`
	Verified       bool   `json:"-"`
	VerificationID string `json:"-"`
}

// NewAccount creates a new account with the provided email/password.
// The returned account will have a hashed password instead of the original.
func NewAccount(email, password string) (*Account, error) {
	hashedPassword, err := hash(password)
	if err != nil {
		return nil, err
	}

	return &Account{
		Email:          email,
		Password:       hashedPassword,
		Verified:       false,
		VerificationID: uuid.New().String(),
	}, nil
}

// ValidPassword returns true if the provided password matches the account password
func (a *Account) ValidPassword(password string) (bool, error) {
	return compare(password, a.Password)
}

// VerifyEmail marks the email as verified and removes the corresponding ID
func (a *Account) VerifyEmail() {
	a.Verified = true
	a.VerificationID = ""
}

// Password hashing functions

type passwordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	size    uint32
}

var c = &passwordConfig{
	time:    1,
	memory:  64 * 1024,
	threads: 4,
	size:    32,
}

func hash(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.size)

	encHash := base64.RawStdEncoding.EncodeToString(hash)
	encSalt := base64.RawStdEncoding.EncodeToString(salt)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, encSalt, encHash)

	return full, nil
}

func compare(providedPassword, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")

	lc := &passwordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &lc.memory, &lc.time, &lc.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	lc.size = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(providedPassword), salt, lc.time, lc.memory, lc.threads, lc.size)

	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}

package user

// PasswordMismatch error for when a provided password is incorrect
type PasswordMismatch struct{}

func (e PasswordMismatch) Error() string {
	return "passwords do not match"
}

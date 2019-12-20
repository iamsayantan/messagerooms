package messagerooms

import "errors"

// User type represents an user.
type User struct {
	ID       string
	Password string
	Nickname string
}

var (
	// ErrNicknameAlreadyTaken is returned when a nickname is already takne by someone.
	ErrNicknameAlreadyTaken = errors.New("the nickname is already taken")

	// ErrUserNotFound is returned when no user is found with the given field.
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository provides methods for interacting with User storage.
type UserRepository interface {
	Create(nickname, password string) (*User, error)
	FindByID(id string) (*User, error)
	FindByNickname(nickname string) (*User, error)
}

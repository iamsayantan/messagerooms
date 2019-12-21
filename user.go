package messagerooms

// User type represents an user.
type User struct {
	ID       string `json:"id"`
	Password string `json:"-"`
	Nickname string `json:"nickname"`
}

// UserRepository provides methods for interacting with User storage.
type UserRepository interface {
	Create(nickname, password string) (*User, error)
	FindByID(id string) (*User, error)
	FindByNickname(nickname string) (*User, error)
}

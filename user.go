package messagerooms

// User type represents an user.
type User struct {
	ID       string `json:"id"`
	Password string `json:"-"`
	Nickname string `json:"nickname"`
}

// GetPersonalTopics returns the personal subscription topics for the user.
// When a new persistent connection is made with the pubsub system we subscribe the connection
// with these topics by default.
func (u *User) GetPersonalTopics() []string {
	// personal topics are in the format of topicName:userID
	return []string{
		TopicNewMessage + ":" + u.ID,
		TopicNewRoom + ":" + u.ID,
	}
}

// UserRepository provides methods for interacting with User storage.
type UserRepository interface {
	Create(nickname, password string) (*User, error)
	FindByID(id string) (*User, error)
	FindByNickname(nickname string) (*User, error)
}

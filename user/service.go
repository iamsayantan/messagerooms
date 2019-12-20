package user

import messagerooms "github.com/iamsayantan/MessageRooms"

import "errors"

var (
	// ErrUserAlreadyExists is returned when an user registers with an existing nickname.
	ErrUserAlreadyExists = errors.New("an user already exists with the nickname")

	// ErrInvalidNickname returned when user tries to login with non existing nickname.
	ErrInvalidNickname = errors.New("invalid nickname")

	// ErrInvalidPassword is returned when user tries to login with wrong password.
	ErrInvalidPassword = errors.New("invalid password")
)

// Service is the interface that provides user related methods.
type Service interface {
	// NewUser creates a new user.
	NewUser(nickname, password string) (*messagerooms.User, error)

	// Login checks for valid nickname and password and returns the user.
	Login(nickname, password string) (*messagerooms.User, error)
}

type userService struct {
	user messagerooms.UserRepository
}

func (s *userService) NewUser(nickname, password string) (*messagerooms.User, error) {
	_, err := s.user.FindByNickname(nickname)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	var user *messagerooms.User
	user, err = s.user.Create(nickname, password)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(nickname, password string) (*messagerooms.User, error) {
	user, err := s.user.FindByNickname(nickname)
	if err != nil {
		return nil, ErrInvalidNickname
	}

	// passwords are not encrypted, so just string matching
	if user.Password != password {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// NewService creates an user service with required dependencies.
func NewService(user messagerooms.UserRepository) Service {
	return &userService{
		user: user,
	}
}

package mysql

import (
	"errors"

	messagerooms "github.com/iamsayantan/MessageRooms"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrUserNotFound is returned when not user data is found with the given option.
	ErrUserNotFound = errors.New("user not found")
)

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) Create(nickname, password string) (*messagerooms.User, error) {
	id, _ := uuid.NewV4()
	user := messagerooms.User{
		ID:       id.String(),
		Nickname: nickname,
		Password: password,
	}

	u.db.Create(&user)

	return &user, nil
}

func (u *userRepository) FindByID(id string) (*messagerooms.User, error) {
	user := messagerooms.User{}

	if notFound := u.db.Where("id = ?", id).First(&user).RecordNotFound(); notFound {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (u *userRepository) FindByNickname(nickname string) (*messagerooms.User, error) {
	user := messagerooms.User{}

	if notFound := u.db.Where("nickname = ?", nickname).First(&user).RecordNotFound(); notFound {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

// NewUserRepository returns implementation of UserRepository interface.
func NewUserRepository(db *gorm.DB) messagerooms.UserRepository {
	return &userRepository{db: db}
}

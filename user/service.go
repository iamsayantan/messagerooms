package user

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	messagerooms "github.com/iamsayantan/MessageRooms"
)

// JwtSigningSecret used for signing and verifying jwt tokens.
const JwtSigningSecret = "secret"

var (
	// ErrUserAlreadyExists is returned when an user registers with an existing nickname.
	ErrUserAlreadyExists = errors.New("an user already exists with the nickname")

	// ErrInvalidNickname returned when user tries to login with non existing nickname.
	ErrInvalidNickname = errors.New("invalid nickname")

	// ErrInvalidPassword is returned when user tries to login with wrong password.
	ErrInvalidPassword = errors.New("invalid password")

	// ErrInvalidAccessToken is returned when we are unable to validate access token.
	ErrInvalidAccessToken = errors.New("invalid access token")
)

// Service is the interface that provides user related methods.
type Service interface {
	// NewUser creates a new user.
	NewUser(nickname, password string) (*messagerooms.User, error)

	// Login checks for valid nickname and password and returns the user.
	Login(nickname, password string) (*messagerooms.User, error)

	// GenerateAuthToken generates an authentication token for the user. This is used after successful login.
	GenerateAuthToken(user messagerooms.User) (string, error)

	// VerifyAuthToken for valid authentication token.
	VerifyAuthToken(token string) (*messagerooms.User, error)
}

// JWTClaims represents the JWT token payload
type JWTClaims struct {
	Nickname string `json:"nickname"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
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

func (s *userService) GenerateAuthToken(user messagerooms.User) (string, error) {
	jwtKey := []byte(JwtSigningSecret)
	expirationTime := time.Now().Add(time.Hour * 24 * 365) // valid for one year

	claims := &JWTClaims{
		Nickname: user.Nickname,
		UserID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func (s *userService) VerifyAuthToken(token string) (*messagerooms.User, error) {
	jwtKey := []byte(JwtSigningSecret)

	claims := &JWTClaims{}
	tokn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !tokn.Valid {
		return nil, ErrInvalidAccessToken
	}

	u, err := s.user.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrInvalidAccessToken
	}

	return u, nil
}

// NewService creates an user service with required dependencies.
func NewService(user messagerooms.UserRepository) Service {
	return &userService{
		user: user,
	}
}

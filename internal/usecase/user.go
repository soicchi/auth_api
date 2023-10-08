package usecase

import (
	"fmt"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/utils"
)

type UserServiceImpl struct {
	UserRepo  UserRepository
	TokenRepo RefreshTokenRepository
}

type UserRepository interface {
	CreateUser(user *models.User) (uint, error)
	FetchUserByEmail(email string) (*models.User, error)
	FetchUsers() ([]models.User, error)
}

type ResponseUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func NewUserServiceImpl(userRepo UserRepository, tokenRepo RefreshTokenRepository) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
	}
}

func newResponseUser(user *models.User) ResponseUser {
	return ResponseUser{
		ID:    user.ID,
		Email: user.Email,
	}
}

func (s *UserServiceImpl) CreateUser(email string, password string) (map[string]string, error) {
	tokens := make(map[string]string)

	// hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return tokens, fmt.Errorf("error hashing password %v", err)
	}

	// generate refresh token
	token, err := utils.GenerateToken()
	if err != nil {
		return tokens, fmt.Errorf("error generating token %v", err)
	}

	refreshToken := models.NewRefreshToken(token)
	user := models.NewUser(email, hashedPassword, refreshToken)

	userID, err := s.UserRepo.CreateUser(user)
	if err != nil {
		return tokens, fmt.Errorf("error creating user %v", err)
	}

	// generate access token
	accessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		return tokens, fmt.Errorf("error generating access token %v", err)
	}

	tokens["accessToken"] = accessToken
	tokens["refreshToken"] = refreshToken.Token

	return tokens, nil
}

func (s *UserServiceImpl) CheckSignIn(email, password string) error {
	user, err := s.UserRepo.FetchUserByEmail(email)
	if err != nil {
		return fmt.Errorf("error getting user by email %v", err)
	}

	if !utils.ValidatePassword(user.Password, password) {
		return fmt.Errorf("error validating password")
	}

	return nil
}

func (s *UserServiceImpl) FetchAllUsers() ([]models.User, error) {
	users, err := s.UserRepo.FetchUsers()
	if err != nil {
		return nil, fmt.Errorf("error getting users %v", err)
	}

	if users == nil {
		return nil, nil
	}

	return users, nil
}

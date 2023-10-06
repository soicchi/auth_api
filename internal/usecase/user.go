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

type AllUsersResponse struct {
	Users []ResponseUser `json:"users"`
}

type CreateUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

func newAllUsersResponse(users []models.User) AllUsersResponse {
	responseUsers := make([]ResponseUser, len(users))
	for _, user := range users {
		responseUser := newResponseUser(&user)
		responseUsers = append(responseUsers, responseUser)
	}

	return AllUsersResponse{
		Users: responseUsers,
	}
}

func newCreateUserResponse(accessToken, refreshToken string) CreateUserResponse {
	return CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// It is not appropriate to return a structure of type Response in this function.
// But it returns a structure of type Response for the time being due to the responsibility of usecase.
func (s *UserServiceImpl) CreateUser(email string, password string) (CreateUserResponse, error) {
	var response CreateUserResponse

	// hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return response, fmt.Errorf("error hashing password %v", err)
	}

	// generate refresh token
	token, err := utils.GenerateToken()
	if err != nil {
		return response, fmt.Errorf("error generating token %v", err)
	}

	refreshToken := models.NewRefreshToken(token)
	user := models.NewUser(email, hashedPassword, refreshToken)

	userID, err := s.UserRepo.CreateUser(user)
	if err != nil {
		return response, fmt.Errorf("error creating user %v", err)
	}

	// generate access token
	accessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		return response, fmt.Errorf("error generating access token %v", err)
	}

	response = newCreateUserResponse(accessToken, token)
	return response, nil
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

func (s *UserServiceImpl) FetchAllUsers() (AllUsersResponse, error) {
	allUsersResponse := AllUsersResponse{}
	users, err := s.UserRepo.FetchUsers()
	if err != nil {
		return allUsersResponse, fmt.Errorf("error getting users %v", err)
	}

	allUsersResponse = newAllUsersResponse(users)

	return allUsersResponse, nil
}

package services

import (
	"errors"

	"inventory-api/internal/config"
	"inventory-api/internal/models"
	"inventory-api/internal/repositories"
	"inventory-api/internal/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
		config:   cfg,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}
	
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	}
	
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	
	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	
	if user == nil {
		return "", errors.New("invalid credentials")
	}
	
	if !user.CheckPassword(req.Password) {
		return "", errors.New("invalid credentials")
	}
	
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, user.Role, s.config.JWTSecret, s.config.JWTExpireHours)
	if err != nil {
		return "", err
	}
	
	return token, nil
}


func (s *AuthService) GetUserProfile(userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	user.Password = ""
	return user, nil
}
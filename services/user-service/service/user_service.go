package service

import (
	"errors"
	"services/user-service/model"
	"services/user-service/repository"
	"services/user-service/utils"

	"github.com/google/uuid"
)

type UserService interface {
	Register(req model.RegisterRequest) (*model.User, error)
	Login(req model.LoginRequest) (*model.User, error)
	GetProfile(userID uuid.UUID) (*model.User, error)
	UpdateProfile(userID uuid.UUID, req model.UpdateProfileRequest) (*model.User, error)
	DeleteProfile(userID uuid.UUID) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(req model.RegisterRequest) (*model.User, error) {
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashed,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Gender:       req.Gender,
		Role:         model.UserRole(req.Role),
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) Login(req model.LoginRequest) (*model.User, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) GetProfile(userID uuid.UUID) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) UpdateProfile(userID uuid.UUID, req model.UpdateProfileRequest) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Image != "" {
		user.Image = req.Image
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Description != "" {
		user.Description = req.Description
	}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) DeleteProfile(userID uuid.UUID) error {
	return s.repo.Delete(userID)
}
